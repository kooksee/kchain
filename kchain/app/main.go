package app

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	kcfg "kchain/types/cfg"
	"kchain/types/cnst"
	"kchain/types/code"
)

//-----------------------------------------

var _ types.Application = (*PersistentApplication)(nil)

type PersistentApplication struct {
	types.BaseApplication
	ValUpdates       []*types.Validator
	GenesisValidator string
	blockHeader      *types.Header
	blockhash        []byte
}

func Run() *PersistentApplication {
	return NewPersistentApplication(
		"kchain",
		cfg().Config.RootDir,
	)
}

func NewPersistentApplication(name, dbDir string) *PersistentApplication {
	logger = kcfg.GetLogWithKeyVals("module", "app")

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
	key := []byte(cnst.ValidatorPrefix + hex.EncodeToString(pk.Bytes()))

	if !state.Has(key) {
		m := "Please contact the administrator to join the node"
		logger.Error(m, "key", key)
		return errors.New(m)
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
			Log:  err.Error(),
		}
	}

	d := strings.Split(tx.Path, ".")
	path := d[0]
	db := d[1]

	switch path {
	case "db", "const_db":
		data, _ := json.MarshalToString(map[string]interface{}{
			"sender":       tx.PubKey,
			"block_height": app.blockHeader.Height,
			"block_hash":   hex.EncodeToString(app.blockhash),
			"time":         app.blockHeader.Time,
			"data":         tx.Value,
		})
		state.Set([]byte(f("%s:%s", db, tx.Key)), []byte(data))

	case "validator":
		d, _ := hex.DecodeString(tx.Key)
		d1, _ := strconv.Atoi(tx.Value)
		if err := app.updateValidator(&types.Validator{PubKey: d, Power: int64(d1)}); err != nil {
			return types.ResponseDeliverTx{
				Code: code.ErrValidatorAdd.Code,
				Log:  err.Error(),
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
			Log:  err.Error(),
		}
	}

	// verify sign
	if err := tx.Verify(); err != nil {
		return types.ResponseCheckTx{
			Code: code.ErrTransactionVerify.Code,
			Log:  err.Error(),
		}
	}

	d := strings.Split(tx.Path, ".")

	if len(d) != 2 {
		return types.ResponseCheckTx{
			Code: code.ErrTransactionVerify.Code,
			Log:  fmt.Sprintf("the path %s is error", tx.Path),
		}
	}

	path := d[0]
	db := d[1]

	switch path {
	case "db":
	case "const_db":
		if state.Has([]byte(f("%s:%s", db, tx.Key) )) {
			return types.ResponseCheckTx{
				Code: code.ErrTransactionVerify.Code,
				Log:  fmt.Sprintf("the key %s already exists", tx.Key),
			}
		}
	case "validator":
		//logger.Error(tx.PubKey)
		//logger.Error(app.GenesisValidator)
		if tx.PubKey != app.GenesisValidator {
			return types.ResponseCheckTx{
				Code: code.ErrTransactionVerify.Code,
				Log:  "Please contact the administrator to add validator",
			}
		}
		if _, err := hex.DecodeString(tx.Key); err != nil {
			return types.ResponseCheckTx{
				Code: code.ErrHexDecode.Code,
				Log:  err.Error(),
			}
		}
		if d, err := strconv.Atoi(tx.Value); err != nil {
			return types.ResponseCheckTx{
				Code: code.ErrJsonDecode.Code,
				Log:  err.Error(),
			}
		} else {
			if d > 9 {
				return types.ResponseCheckTx{
					Code: code.ErrVerify.Code,
					Log:  "the node power must be less than 10",
				}
			}
		}

	default:
		return types.ResponseCheckTx{
			Code: code.ErrUnknownMathod.Code,
			Log:  "unknown path",
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
		logger.Info("Commit block", "height", height, "root", hex.EncodeToString(appHash))
		return types.ResponseCommit{Code: code.Ok.Code, Data: appHash}
	}
}

func (app *PersistentApplication) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {
	var (
		key = reqQuery.Data
	)

	d := strings.Split(reqQuery.Path, ".")
	path := d[0]
	db := ""
	if len(d) == 2 {
		db = d[1]
	}

	switch path {
	case "db", "const_db":
		index, value := state.Get([]byte(f("%s:%s", db, string(key))))
		res.Code = types.CodeTypeOK
		res.Index = int64(index)
		res.Key = key
		res.Value = value
		if value != nil {
			res.Log = "exists"
		} else {
			res.Log = "does not exist"
		}

	case "keys":

		s := strings.Split(string(key), ":")

		if len(s) != 2 {
			res.Code = code.ErrTransactionDecode.Code
			res.Log = f("error range %s", key)
			return
		}

		s_f := s[0]
		s_t := s[1]
		i_f, _ := strconv.Atoi(s_f)
		i_t, _ := strconv.Atoi(s_t)

		d := []string{}

		for i := i_f; i <= i_t; i++ {

			if k, _ := state.GetByIndex(i); k != nil {
				if !bytes.HasPrefix(k, []byte("val:")) && !bytes.HasPrefix(k, []byte("__app:")) {
					d = append(d, string(k))
				}
			} else {
				continue
			}
		}

		res.Value, _ = json.Marshal(d)

	default:
		res.Code = code.ErrUnknownMathod.Code
		res.Log = "unknown path"
	}
	return
}

// Save the validators in the merkle tree
func (app *PersistentApplication) InitChain(req types.RequestInitChain) types.ResponseInitChain {

	logger.Info("InitChain")
	for _, v := range req.Validators {
		// 最高权限拥有者
		if v.Power == 10 {

			state.Set([]byte("__app:genesis_validator"), v.PubKey)

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
	app.blockhash = req.Hash

	_, d := state.Get([]byte("__app:genesis_validator"))
	app.GenesisValidator = hex.EncodeToString(d)

	return types.ResponseBeginBlock{}
}

func (app *PersistentApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

//---------------------------------------------

// update validators
func (app *PersistentApplication) Validators() (validators []*types.Validator) {
	state.Iterate(func(key, value []byte) bool {
		if strings.HasPrefix(string(key), cnst.ValidatorPrefix) {
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

	if v.Power < 0 {
		state.Remove(key)
		v.Power = 0
		logger.Info("delete node ok", "key", key)
	}

	if v.Power == 0 {
		if !state.Has(key) {
			return errors.New(fmt.Sprintf("Cannot remove non-existent validator %X", key))
		}
	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(v, value); err != nil {
			return errors.New(fmt.Sprintf("Error encoding validator: %v", err))
		}
		state.Set(key, value.Bytes())

		logger.Info("save node ok", "key", key)
	}

	app.ValUpdates = append(app.ValUpdates, v)

	return nil
}
