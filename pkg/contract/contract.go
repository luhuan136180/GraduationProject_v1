package contract

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/big"
	"time"
	"v1/pkg/apiserver/contract"
	"v1/pkg/dao"
	"v1/pkg/model"
)

// var ContractClient Client

type Client struct {
	Client  *ethclient.Client
	Token   *contract.Token
	Options *Options
}

func InitContract(option *Options, db *gorm.DB) {
	ctx := context.Background()

	if option.NetworkIP == "" {
		zap.L().Debug("contract option is err.....")
		fmt.Println("start failed...")
		return
	}
	// 链接合约
	client, err := newClient(option)
	if err != nil {
		return
	}

	// 存储起始节点
	_, err = client.SaveValue("start", []string{"start"}, 0)
	if err != nil {
		zap.L().Error("client.SaveValue", zap.Error(err))
		return
	}

	// 获取起始节点
	resultArr, _, err := client.GetValue("start")
	if err != nil {
		zap.L().Error("client.GetValue err:", zap.Error(err))
		return
	}

	if resultArr[0] != "start" {
		return
	}
	// 写到这里了
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		fmt.Println("start .......contract")
		for {
			select {
			case <-ticker.C:
				err := client.Save(ctx, db)
				if err != nil {
					zap.L().Debug("save data is failed:", zap.Error(err))
				}
				fmt.Println("waiting....")
			}
		}
	}()
	// 最后成功后再存
	contractConfig := model.ContractConfig{
		NetworkIP:         option.NetworkIP,
		AccountPrivateKey: option.AccountPrivateKey,
		Address:           option.Address,
		GasFeeCap:         big.NewInt(option.GasFeeCap),
		GasLimit:          option.GasLimit,
	}

	dao.InitConfiguration(ctx, db, model.Configuration{
		Tag:     string(model.ContractTAG),
		Content: contractConfig,
	})
}

func newClient(option *Options) (*Client, error) {
	client, err := ethclient.Dial(option.NetworkIP)
	if err != nil {
		zap.L().Debug("contract start end : ethclient.Dial error :", zap.Error(err))
		return nil, err
	}

	token, err := contract.NewToken(common.HexToAddress(option.Address), client)
	if err != nil {
		zap.L().Debug("contract start end : NewToken error : ", zap.Error(err))
		return nil, err
	}

	ContractClient := Client{
		Client:  client,
		Token:   token,
		Options: option,
	}
	return &ContractClient, nil
}

func (client *Client) getOpt() (*bind.TransactOpts, error) {
	// 获取当前最新区块
	chainID, err := client.Client.ChainID(context.Background())
	if err != nil {
		fmt.Println("获取ChainID失败:", err)
		return nil, err
	}
	fmt.Println("chain id :", chainID)

	privateKeyECDSA, err := crypto.HexToECDSA(client.Options.AccountPrivateKey)
	if err != nil {
		fmt.Println("crypto.HexToECDSA error ,", err)
		return nil, err
	}

	gasTipCap, _ := client.Client.SuggestGasTipCap(context.Background())

	// 构建参数对象
	opts, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainID)
	if err != nil {
		fmt.Println("bind.NewKeyedTransactorWithChainID error ,", err)
		return nil, err
	}

	// 设置参数
	opts.GasFeeCap = big.NewInt(108694000460)
	opts.GasLimit = uint64(0)
	opts.GasTipCap = gasTipCap

	return opts, nil
}

func (client *Client) GetValue(key string) ([]string, *big.Int, error) {
	return client.Token.Get(nil, key)
}

func (client *Client) SaveValue(key string, arr []string, dataType int64) (*types.Transaction, error) {
	opts, err := client.getOpt()
	if err != nil {
		zap.L().Error("client.Token.SaveValue save:", zap.Error(err))
		return nil, err
	}

	Transaction, err := client.Token.Save(opts, key, arr, big.NewInt(dataType))
	if err != nil {
		zap.L().Error("client.Token.Save save:", zap.Error(err))
		return nil, err
	}
	return Transaction, err
}
