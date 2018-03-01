package commands

import (
	"github.com/spf13/cobra"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/types"
	"fmt"
	"encoding/hex"

	"kchain/services"
)

var msg string = "hello"
var sign_msg string = "type[hex string]"

// AddNodeFlags exposes some common configuration options on the command-line
// These are exposed for convenience of commands embedding a tendermint node
func AddNodeFlags1(cmd *cobra.Command) *cobra.Command {

	sign := &cobra.Command{
		Use:   "sign",
		Short: "签名命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			pvfs := types.LoadOrGenPrivValidatorFS(config.PrivValidatorFile())
			fmt.Println(hex.EncodeToString(pvfs.PrivKey.Sign([]byte(msg)).Bytes()))
			return nil
		},
	}
	sign.Flags().StringVar(&msg, "msg", msg, "需要签名的数据")

	verify := &cobra.Command{
		Use:   "verify",
		Short: "校验命令",
		RunE: func(cmd *cobra.Command, args []string) error {

			if d, e := hex.DecodeString(sign_msg); e != nil {
				panic(e.Error())
			} else {
				if sig, err := crypto.SignatureFromBytes(d); err != nil {
					panic(err.Error())
				} else {
					pvfs := types.LoadOrGenPrivValidatorFS(config.PrivValidatorFile())
					fmt.Println(pvfs.PubKey.VerifyBytes([]byte(msg), sig))
					return nil
				}
			}
		},
	}
	verify.Flags().StringVar(&sign_msg, "s_msg", sign_msg, "签名后的信息")
	verify.Flags().StringVar(&msg, "msg", msg, "签名信息")

	node := &cobra.Command{
		Use:   "node",
		Short: "服务启动",
		RunE: func(cmd *cobra.Command, args []string) error {
			services.Run()
			return nil

		},
	}
	cmd.Flags().StringVar(&kcfg().App.Addr, "addr", kcfg().App.Addr, "kchain services port")

	cmd.AddCommand(sign, verify, node)
	return cmd
}



// NewRunNodeCmd returns the command that allows the CLI to start a
// node. It can be used with a custom PrivValidator and in-process ABCI application.
func ServicesCmd() *cobra.Command {
	return AddNodeFlags1(&cobra.Command{
		Use:   "services",
		Short: "kchain工具服务",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
}
