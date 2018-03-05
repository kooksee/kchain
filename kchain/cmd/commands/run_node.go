package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	nm "github.com/tendermint/tendermint/node"

	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/proxy"
	"kchain/types/cfg"

	"kchain/abci"
)

var kcfg = cfg.GetConfig()

// AddNodeFlags exposes some common configuration options on the command-line
// These are exposed for convenience of commands embedding a tendermint node
func AddNodeFlags(cmd *cobra.Command) *cobra.Command {
	_cfg := kcfg()

	// app falgs
	cmd.Flags().StringVar(&_cfg.App.Addr, "addr", _cfg.App.Addr, "kchain web port")

	// bind flags
	cmd.Flags().StringVar(&config.Moniker, "moniker", config.Moniker, "Node Name")

	// node flags
	cmd.Flags().BoolVar(&config.FastSync, "fast_sync", config.FastSync, "Fast blockchain syncing")

	// abci flags
	//cmd.Flags().StringVar(&config.ProxyApp, "proxy_app", config.ProxyApp, "Proxy app address, or 'nilapp' or 'dummy' for local testing.")
	//cmd.Flags().StringVar(&config.ABCI, "abci", config.ABCI, "Specify abci transport (socket | grpc)")

	// rpc flags
	//cmd.Flags().StringVar(&config.RPC.GRPCListenAddress, "rpc.grpc_laddr", config.RPC.GRPCListenAddress, "GRPC listen address (BroadcastTx only). Port required")
	//cmd.Flags().BoolVar(&config.RPC.Unsafe, "rpc.unsafe", config.RPC.Unsafe, "Enabled unsafe rpc methods")

	// p2p flags
	cmd.Flags().StringVar(&config.P2P.ListenAddress, "p2p.laddr", config.P2P.ListenAddress, "Node listen address. (0.0.0.0:0 means any interface, any port)")
	cmd.Flags().StringVar(&config.P2P.Seeds, "p2p.seeds", config.P2P.Seeds, "Comma delimited host:port seed nodes")
	//cmd.Flags().BoolVar(&config.P2P.SkipUPNP, "p2p.skip_upnp", config.P2P.SkipUPNP, "Skip UPNP configuration")
	//cmd.Flags().BoolVar(&config.P2P.PexReactor, "p2p.pex", config.P2P.PexReactor, "Enable/disable Peer-Exchange")

	// consensus flags
	//cmd.Flags().BoolVar(&config.Consensus.CreateEmptyBlocks, "consensus.create_empty_blocks", config.Consensus.CreateEmptyBlocks, "Set this to false to only produce blocks when there are txs or when the AppHash changes")

	return cmd
}

// NewRunNodeCmd returns the command that allows the CLI to start a
// node. It can be used with a custom PrivValidator and in-process ABCI application.
func NewRunNodeCmd() *cobra.Command {
	return AddNodeFlags(&cobra.Command{
		Use:   "node",
		Short: "Run the kchain node",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 初始化配置
			kcfg().Config = config

			abci_app := abci.Run()

			pvfs := types.LoadOrGenPrivValidatorFS(config.PrivValidatorFile())

			// 启动abci服务和tendermint节点
			n, err := nm.NewNode(
				config,
				pvfs,
				proxy.NewLocalClientCreator(abci_app),
				nm.DefaultGenesisDocProviderFunc(config),
				nm.DefaultDBProvider,
				logger,
			)
			if err != nil {
				return fmt.Errorf("Failed to create node: %v", err)
			}

			//if otherEd, ok := pvfs.PrivKey.Unwrap().(crypto.PrivKeyEd25519); ok {
			//	logger.Error(fmt.Sprintf("save key %s ok", hex.EncodeToString(otherEd.Bytes())))
			//	n.Switch().SetNodePrivKey(otherEd)
			//}

			// 新加入节点的过滤逻辑
			n.Switch().SetPubKeyFilter(abci_app.PubKeyFilter)

			if err := n.Start(); err != nil {
				return fmt.Errorf("Failed to start node: %v", err)
			} else {
				logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())
			}

			// 得到正在运行的tendermint
			kcfg().Node = n

			// 启动应用
			//app.Run()

			n.RunForever()

			return nil
		},
	})
}
