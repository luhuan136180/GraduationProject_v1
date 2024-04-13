package contract

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

const (
	ContractNetworkIP         = "contract_network_ip"
	ContractAccountPrivateKey = "contract_account_private_key"
	ContractAddress           = "contract_address"
	ContractGasFeeCap         = "contract_gas_fee_cap"
	ContractGasLimit          = "contract_gas_limit"
)

type Options struct {
	NetworkIP         string // ip *必填
	AccountPrivateKey string // 所用账号私钥 *必填
	Address           string // 合约地址 *必填
	GasFeeCap         int64  // Gas费用上限
	GasLimit          uint64 // Gas费用上限(单次交易) ps：默认设置为0

	v *viper.Viper
}

func NewLoggerOptions() *Options {
	o := &Options{
		NetworkIP:         "",
		AccountPrivateKey: "",
		Address:           "",
		GasFeeCap:         0,
		GasLimit:          0,

		v: viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))),
	}

	o.v.AutomaticEnv()
	return o
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.NetworkIP, ContractNetworkIP, o.NetworkIP, "log level. env LOG_FILENAME")
	fs.StringVar(&o.AccountPrivateKey, ContractAccountPrivateKey, o.AccountPrivateKey, "log level. env LOG_FILENAME")
	fs.StringVar(&o.Address, ContractAddress, o.Address, "log level. env LOG_FILENAME")
	fs.Int64Var(&o.GasFeeCap, ContractGasFeeCap, o.GasFeeCap, "log level. env LOG_FILENAME")
	fs.Uint64Var(&o.GasLimit, ContractGasLimit, o.GasLimit, "log level. env LOG_FILENAME")

	_ = o.v.BindPFlags(fs)
}
