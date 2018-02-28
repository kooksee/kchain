package cfg

import (
	cfg "github.com/tendermint/tendermint/config"
	//dbm "github.com/tendermint/tmlibs/db"
	tlog "github.com/tendermint/tmlibs/log"
	tmflags "github.com/tendermint/tmlibs/cli/flags"
	//ttypes "github.com/tendermint/tendermint/types"
	crypto "github.com/tendermint/go-crypto"
	c "github.com/tendermint/tendermint/rpc/client"
	"os"
)

func DbSet(k, v string) {
	GetConfig()().store.Set([]byte(k), []byte(v))
}

func DbGet(k string) string {
	return string(GetConfig()().store.Get([]byte(k)))
}

func GetLogWithKeyVals(keyvals ...interface{}) tlog.Logger {
	return GetConfig()().log.With(keyvals...)
}

func GetPrivKey() crypto.PrivKey {
	return GetConfig()().pk.PrivKey
}

func GetPubKey() crypto.PubKey {
	return GetConfig()().pk.GetPubKey()
}

func GetAddress() []byte {
	return GetConfig()().pk.GetAddress().Bytes()
}

func Sign(msg []byte) ([]byte, error) {
	if _sign, err := GetConfig()().pk.Sign(msg); err != nil {
		return nil, err
	} else {
		return _sign.Bytes(), nil
	}
}

func Verify(msg, sign []byte) (bool, error) {
	if _sign, err := crypto.SignatureFromBytes(sign); err != nil {
		return false, err
	} else {
		return GetConfig()().pk.PubKey.VerifyBytes(msg, _sign), nil
	}
}

func Abci() *c.HTTP {
	return GetConfig()().client
}

func initServices() *services {
	s := &services{
		Config:cfg.DefaultConfig(),
		App:&appConfig{
			Name:"kchain",
			Addr:":9000",
		},
		Node:nil,
	}

	// 初始化db
	//if store, err := dbm.NewGoLevelDB(
	//	"app",
	//	s.Config.DBDir(),
	//); err != nil {
	//	panic(err.Error())
	//} else {
	//	s.store = store
	//}

	// 初始化log
	if klog, err := tmflags.ParseLogLevel(
		s.Config.LogLevel,
		tlog.NewTMLogger(tlog.NewSyncWriter(os.Stderr)),
		"error",
	); err != nil {
		panic(err.Error())
	} else {
		s.log = klog
	}

	// 私钥获得
	//s.pk = ttypes.LoadOrGenPrivValidatorFS(s.Config.PrivValidatorFile())

	s.client = c.NewHTTP(s.Config.RPC.ListenAddress, "/ws")

	return s
}

func GetConfig() func() *services {
	return func() *services {
		once.Do(func() {
			instance = initServices()
		})
		return instance
	}
}


