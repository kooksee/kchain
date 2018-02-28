package abci

import (
	"bytes"
	"fmt"
	"strings"
	"github.com/tendermint/tmlibs/log"
	"github.com/tendermint/iavl"
	"github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
	crypto "github.com/tendermint/go-crypto"

	"kchain/types/code"
	"kchain/types/cnst"
	kcfg "kchain/types/cfg"
	"encoding/hex"
	"github.com/pkg/errors"
	"strconv"
)

//-----------------------------------------

var _ types.Application = (*PersistentApplication)(nil)

type PersistentApplication struct {
	types.BaseApplication
	ValUpdates       []*types.Validator
	GenesisValidator string
	blockHeader      *types.Header
}

func Run() *PersistentApplication {
	return NewPersistentApplication(
		"kchain",
		cfg().Config.DBDir(),
	)
}

func NewPersistentApplication(name, dbDir string) *PersistentApplication {
	logger = kcfg.GetLogWithKeyVals("module", "abci")

	db, err := dbm.NewGoLevelDB(name, dbDir)
	if err != nil {
		panic(err.Error())
	}

	state = iavl.NewVersionedTree(0, db)
	state.Load()

	return &PersistentApplication{}
}

func (app *PersistentApplication) SetLogger(l log.Logger) {
	logger = l
}

// 新节点连接过滤
func (app *PersistentApplication) PubKeyFilter(pk crypto.PubKeyEd25519) error {
	key := []byte(cnst.ValidatorPrefix + hex.EncodeToString(pk.Wrap().Bytes()))
	if !state.Has(key) {
		return errors.New("Please contact the administrator to join the node")
	}
	return nil
}

// 实现abci的Info协议
func (app *PersistentApplication) Info(req types.RequestInfo) (res types.ResponseInfo) {
	res.Data = cfg().Config.Moniker
	res.LastBlockHeight = int64(state.LatestVersion())
	res.LastBlockAppHash = state.Hash()
	res.Version = req.Version

	return
}

// 实现abci的SetOption协议
func (app *PersistentApplication) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{Code: types.CodeTypeOK}
}

// 实现abci的DeliverTx协议
func (app *PersistentApplication) DeliverTx(txBytes []byte) types.ResponseDeliverTx {
	tx := NewTransaction()

	// decode tx
	if err := tx.FromBytes(txBytes); err != nil {
		return types.ResponseDeliverTx{
			Code: code.ErrTransactionDecode.Code,
			Log:err.Error(),
		}
	}

	switch tx.Path {
	case "db":
		data := fmt.Sprintf("%d,,%s", app.blockHeader.Height, tx.Value)
		state.Set([]byte("db:" + tx.Key), []byte(data))
	case "const_db":
		data := fmt.Sprintf("%d,,%s", app.blockHeader.Height, tx.Value)
		state.Set([]byte("const_db:" + tx.Key), []byte(data))
	case "validator":
		d, _ := hex.DecodeString(tx.Key)
		d1, _ := strconv.Atoi(tx.Value)
		if err := app.updateValidator(&types.Validator{PubKey:d, Power:int64(d1)}); err != nil {
			return types.ResponseDeliverTx{
				Code: code.ErrValidatorAdd.Code,
				Log:err.Error(),
			}
		}
	}

	return types.ResponseDeliverTx{Code: code.Ok.Code}
}

// 实现abci的CheckTx协议
func (app *PersistentApplication) CheckTx(txBytes []byte) types.ResponseCheckTx {
	tx := NewTransaction()

	// decode tx
	if err := tx.FromBytes(txBytes); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrTransactionDecode.Code,
			Log:err.Error(),
		}
	}

	// verify sign
	if err := tx.Verify(); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrTransactionVerify.Code,
			Log:err.Error(),
		}
	}

	switch tx.Path {
	case "db":
	case "const_db":
		if state.Has([]byte("const_db:" + tx.Key)) {
			return types.ResponseCheckTx{
				Code: code.ErrTransactionVerify.Code,
				Log:fmt.Sprintf("the key %s already exists", tx.Key),
			}
		}
	case "validator":
		logger.Error(tx.PubKey)
		logger.Error(app.GenesisValidator)
		if strings.Compare(tx.PubKey, app.GenesisValidator) != 0 {
			return types.ResponseCheckTx{
				Code: code.ErrTransactionVerify.Code,
				Log:"Please contact the administrator to add validator",
			}
		}
		if _, err := hex.DecodeString(tx.Key); err != nil {
			return types.ResponseCheckTx{
				Code: code.ErrHexDecode.Code,
				Log:err.Error(),
			}
		}
		if _, err := strconv.Atoi(tx.Value); err != nil {
			return types.ResponseCheckTx{
				Code: code.ErrJsonDecode.Code,
				Log:err.Error(),
			}
		}

	default:
		return types.ResponseCheckTx{
			Code:code.ErrUnknownMathod.Code,
			Log:"unknown path",
		}
	}

	return types.ResponseCheckTx{Code: code.Ok.Code}
}

// Commit will panic if InitChain was not called
func (app *PersistentApplication) Commit() (res types.ResponseCommit) {
	// Save a new version for next height

	height := state.LatestVersion() + 1
	if appHash, err := state.SaveVersion(height); err != nil {
		panic(err)
	} else {
		logger.Info("Commit block", "height", height, "root", appHash)
		return types.ResponseCommit{Code: code.Ok.Code, Data: appHash}
	}
}

func (app *PersistentApplication) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {
	var (
		path = reqQuery.Path
		key = reqQuery.Data
	)

	switch path {
	case "db":
		index, value := state.Get([]byte("db:" + string(key)))
		data := strings.Split(string(value), ",,")
		res.Code = types.CodeTypeOK
		res.Index = int64(index)
		res.Key = key
		res.Value = []byte(data[1])
		h, _ := strconv.Atoi(data[0])
		res.Height = int64(h)
		if value != nil {
			res.Log = "exists"
		} else {
			res.Log = "does not exist"
		}


	case "const_db":
		index, value := state.Get([]byte("const_db:" + string(key)))
		data := strings.Split(string(value), ",,")
		res.Code = types.CodeTypeOK
		res.Index = int64(index)
		res.Key = key
		res.Value = []byte(data[1])
		h, _ := strconv.Atoi(data[0])
		res.Height = int64(h)
		if value != nil {
			res.Log = "exists"
		} else {
			res.Log = "does not exist"
		}

	default:
		res.Code = code.ErrUnknownMathod.Code
		res.Log = "unknown path"
	}
	return
}

// Save the validators in the merkle tree
func (app *PersistentApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	logger.Error(req.String())
	logger.Error("InitChain")

	for _, v := range cfg().Node.GenesisDoc().Validators {
		// 最高权限拥有者
		if v.Power == 10 {
			app.GenesisValidator = hex.EncodeToString(v.PubKey.Bytes())
		}
	}

	for _, v := range req.Validators {
		// 最高权限拥有者
		if v.Power == 10 {
			app.GenesisValidator = hex.EncodeToString(v.PubKey)
		}

		if r := app.updateValidator(v); r != nil {
			logger.Error("Error updating validators", "r", r.Error())
		}
	}
	return types.ResponseInitChain{}
}

func (app *PersistentApplication) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	app.ValUpdates = make([]*types.Validator, 0)
	app.blockHeader = req.Header

	d, _ := json.MarshalToString(req)
	logger.Error(d)

	return types.ResponseBeginBlock{}
}

func (app *PersistentApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

//---------------------------------------------


func isValidatorTx(tx []byte) bool {
	return strings.HasPrefix(string(tx), cnst.ValidatorPrefix)
}
// update validators
func (app *PersistentApplication) Validators() (validators []*types.Validator) {
	state.Iterate(func(key, value []byte) bool {
		if isValidatorTx(key) {
			validator := new(types.Validator)
			err := types.ReadMessage(bytes.NewBuffer(value), validator)
			if err != nil {
				panic(err)
			}
			validators = append(validators, validator)
		}
		return false
	})
	return
}

// add, update, or remove a validator
func (app *PersistentApplication) updateValidator(v *types.Validator) error {
	key := []byte(cnst.ValidatorPrefix + hex.EncodeToString(v.PubKey))
	if v.Power == 0 {
		if !state.Has(key) {
			return errors.New(fmt.Sprintf("Cannot remove non-existent validator %X", key))
		}
		state.Remove(key)
	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(v, value); err != nil {
			return errors.New(fmt.Sprintf("Error encoding validator: %v", err))
		}
		state.Set(key, value.Bytes())
	}

	app.ValUpdates = append(app.ValUpdates, v)

	return nil
}
