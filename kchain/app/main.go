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
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	kcfg "kchain/types/cfg"
	"kchain/types/cnst"
	"kchain/types/code"
	"encoding/binary"
)

// -----------------------------------------

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")
	dataHeight      = "dataHeight:"
)

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func loadState(db dbm.DB) State {
	stateBytes := db.Get(stateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.db = db
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

var _ types.Application = (*PersistentApplication)(nil)

type PersistentApplication struct {
	types.BaseApplication
	ValUpdates       []types.Validator
	GenesisValidator string
	blockHeader      types.Header
	blockhash        []byte
}

func Run() *PersistentApplication {
	return NewPersistentApplication(
		"kchain",
		cfg().Config.DBDir(),
	)
}

func NewPersistentApplication(name, dbDir string) *PersistentApplication {
	logger = kcfg.GetLogWithKeyVals("module", "app")

	db, err := dbm.NewGoLevelDB(name, dbDir)
	if err != nil {
		panic(err)
	}

	state = loadState(db)

	return &PersistentApplication{}
}

func (app *PersistentApplication) SetLogger(l log.Logger) {
	logger = l
}

// 新节点连接过滤
func (app *PersistentApplication) PubKeyFilter(pk crypto.PubKey) error {
	key := []byte(cnst.ValidatorPrefix + hex.EncodeToString(pk.Bytes()))

	if !state.db.Has(key) {
		m := "Please contact the administrator to join the node"
		logger.Error(m, "key", key)
		return errors.New(m)
	}
	return nil
}

// 实现abci的Info协议
func (app *PersistentApplication) Info(req types.RequestInfo) (res types.ResponseInfo) {

	res.Data = cfg().Config.Moniker
	res.LastBlockHeight = int64(state.Height)
	res.LastBlockAppHash = state.AppHash
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

	m, _ := hex.DecodeString(string(txBytes))

	// decode tx
	if err := tx.FromBytes(m); err != nil {
		return types.ResponseDeliverTx{
			Code: code.ErrTransactionDecode.Code,
			Log:  err.Error(),
		}
	}

	d := strings.Split(tx.Path, ".")
	path := d[0]
	db := d[1]

	switch path {
	case "db", "const_db", "admin_db":

		for k, v := range tx.DecodeValues() {

			data, _ := json.Marshal(map[string]interface{}{
				"sender":       tx.PubKey,
				"block_height": app.blockHeader.Height,
				"data_height":  state.Size,
				"block_hash":   hex.EncodeToString(app.blockhash),
				"time":         app.blockHeader.Time,
				"data":         v,
			})

			k := []byte(f("%s:%s", db, k))
			state.db.Set(k, data)
			state.Size += 1
			state.db.Set([]byte(f("%s%d", dataHeight, state.Size)), k)
		}

	case "validator":
		for k, v := range tx.DecodeValues() {
			d, _ := hex.DecodeString(k)
			d1, _ := strconv.Atoi(f("%d", v))
			if err := app.updateValidator(types.Validator{PubKey: d, Power: int64(d1)}); err != nil {
				return types.ResponseDeliverTx{
					Code: code.ErrValidatorAdd.Code,
					Log:  err.Error(),
				}
			}
		}
	}

	return types.ResponseDeliverTx{Code: code.Ok.Code}
}

// 实现abci的CheckTx协议
func (app *PersistentApplication) CheckTx(txBytes []byte) types.ResponseCheckTx {

	tx := NewTransaction()
	m, _ := hex.DecodeString(string(txBytes))

	// decode tx
	if err := tx.FromBytes(m); err != nil {
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
		for k := range tx.DecodeValues() {
			if state.db.Has([]byte(f("%s:%s", db, k) )) {
				return types.ResponseCheckTx{
					Code: code.ErrTransactionVerify.Code,
					Log:  fmt.Sprintf("the key %s already exists", k),
				}
			}
		}

	case "admin_db":
		if tx.PubKey != app.GenesisValidator {
			return types.ResponseCheckTx{
				Code: code.ErrTransactionVerify.Code,
				Log:  "Please contact the administrator to operate the tx",
			}
		}

	case "validator":
		for k, v := range tx.DecodeValues() {
			if tx.PubKey != app.GenesisValidator {
				return types.ResponseCheckTx{
					Code: code.ErrTransactionVerify.Code,
					Log:  "Please contact the administrator to add validator",
				}
			}
			if _, err := hex.DecodeString(k); err != nil {
				return types.ResponseCheckTx{
					Code: code.ErrHexDecode.Code,
					Log:  err.Error(),
				}
			}
			if d, err := strconv.Atoi(f("%d", v)); err != nil {
				return types.ResponseCheckTx{
					Code: code.ErrJsonDecode.Code,
					Log:  err.Error(),
				}
			} else {
				// power等于10是最高的权限
				if d > 9 {
					return types.ResponseCheckTx{
						Code: code.ErrVerify.Code,
						Log:  "the node power must be less than 10",
					}
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

	appHash := make([]byte, 8)
	binary.PutVarint(appHash, state.Size)
	state.AppHash = appHash
	state.Height += 1
	saveState(state)
	return types.ResponseCommit{Data: appHash}
}

func (app *PersistentApplication) Query(reqQuery types.RequestQuery) (res types.ResponseQuery) {
	var (
		key = string(reqQuery.Data)
	)

	d := strings.Split(reqQuery.Path, ".")
	path := d[0]
	db := ""
	if len(d) == 2 {
		db = d[1]
	}

	switch path {
	case "db", "const_db", "admin_db":
		res.Code = types.CodeTypeOK
		res.Value = state.db.Get([]byte(f("%s:%s", db, key)))
		if res.Value != nil {
			res.Code = 0
			res.Log = "exists"
		} else {
			res.Code = 1
			res.Log = "does not exist"
		}

	case "keys":

		s := strings.Split(key, ":")

		if len(s) != 2 {
			res.Code = code.ErrTransactionDecode.Code
			res.Log = f("error range %s", key)
			return
		}

		s_f := s[0]
		s_t := s[1]
		i_f, _ := strconv.Atoi(s_f)
		i_t, _ := strconv.Atoi(s_t)

		// 比较最大值,查询最大值为数据最大高度
		if i_t > int(state.Size) {
			i_t = int(state.Size)
		}

		// 最大查询范围值为1000
		if i_t-i_f > 1000 {
			i_t = i_f + 1000
		}

		d := map[string]int{}

		for i := i_f; i <= i_t; i++ {
			k := []byte(f("%s%d", dataHeight, i))
			v := state.db.Get(k)

			if v == nil {
				continue
			}

			if bytes.HasPrefix(v, []byte("__app:")) {
				continue
			}

			if bytes.HasPrefix(v, []byte(cnst.ValidatorPrefix)) {
				continue
			}

			d[string(v)] = i
		}

		res.Value, _ = json.Marshal(d)
		res.Code = 0

	case "accounts":

		d := map[string]int{}
		for i := 0; i <= int(state.Size); i++ {
			k := []byte(f("%s%d", dataHeight, i))
			v := state.db.Get(k)

			if v == nil {
				continue
			}

			if bytes.HasPrefix(v, []byte(cnst.AccountPrefix)) {
				d[string(bytes.Trim(v, cnst.AccountPrefix))] = i
			}
		}

		res.Value, _ = json.Marshal(d)
		res.Code = 0
	case "metadatas":
		s := strings.Split(key, ":")

		if len(s) != 2 {
			res.Code = code.ErrTransactionDecode.Code
			res.Log = f("error range %s", key)
			return
		}

		s_f := s[0]
		s_t := s[1]
		i_f, _ := strconv.Atoi(s_f)
		i_t, _ := strconv.Atoi(s_t)

		// 比较最大值,查询最大值为数据最大高度
		if i_t > int(state.Size) {
			i_t = int(state.Size)
		}

		// 最大查询范围值为1000
		if i_t-i_f > 1000 {
			i_t = i_f + 1000
		}

		d := map[string]int{}

		for i := i_f; i <= i_t; i++ {
			k := []byte(f("%s%d", dataHeight, i))
			v := state.db.Get(k)

			if v == nil {
				continue
			}

			if bytes.HasPrefix(v, []byte("metadata:")) {
				d[string(bytes.Trim(v, "metadata:"))] = i
			}
		}

		res.Value, _ = json.Marshal(d)
		res.Code = 0

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

			state.db.Set([]byte("__app:genesis_validator"), v.PubKey)

			app.GenesisValidator = hex.EncodeToString(v.PubKey)
		}

		if r := app.updateValidator(v); r != nil {
			logger.Error("Error updating validators", "r", r.Error())
		}
	}
	return types.ResponseInitChain{}
}

func (app *PersistentApplication) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	app.ValUpdates = make([]types.Validator, 0)
	app.blockHeader = req.Header
	app.blockhash = req.Hash

	d := state.db.Get([]byte("__app:genesis_validator"))
	app.GenesisValidator = hex.EncodeToString(d)

	return types.ResponseBeginBlock{}
}

func (app *PersistentApplication) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{ValidatorUpdates: app.ValUpdates}
}

// ---------------------------------------------

// 更新validator
func (app *PersistentApplication) updateValidator(v types.Validator) error {
	key := []byte(cnst.ValidatorPrefix + hex.EncodeToString(v.PubKey))

	// power等于-1的时候,开放节点的权限
	if v.Power == -1 {
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&v, value); err != nil {
			return errors.New(fmt.Sprintf("Error encoding validator: %v", err))
		}

		state.db.Set(key, value.Bytes())
		state.Size += 1
		state.db.Set([]byte(f("%s%d", dataHeight, state.Size)), key)

		logger.Info("save node ok", "key", key)

		v.Power = 0
		app.ValUpdates = append(app.ValUpdates, v)
		return nil
	}

	// power等于-2的时候,删除节点
	if v.Power == -2 {
		state.db.Delete(key)
		logger.Info("delete node ok", "key", key)

		v.Power = 0
		app.ValUpdates = append(app.ValUpdates, v)
		return nil
	}

	// power小于等于0的时候,删除验证节点
	if v.Power >= 0 {
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(&v, value); err != nil {
			return errors.New(fmt.Sprintf("Error encoding validator: %v", err))
		}

		state.db.Set(key, value.Bytes())
		state.Size += 1
		state.db.Set([]byte(f("%s%d", dataHeight, state.Size)), key)

		logger.Info("save node ok", "key", key)

		app.ValUpdates = append(app.ValUpdates, v)
	}
	return nil
}
