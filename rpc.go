package blockchain

import (
	"errors"
	"strconv"
	"time"

	"github.com/onrik/ethrpc"
	"github.com/sirupsen/logrus"
)

type RPC struct {
	rpcImpl       *ethrpc.EthRPC
	maxRetryTimes int
}

func NewRPC(uri string, maxRetryCount int, options ...func(rpc *ethrpc.EthRPC)) *RPC {
	rpc := ethrpc.New(uri, options...)

	return &RPC{rpc, maxRetryCount}
}

func (rpc RPC) getCurrentBlockNum() (uint64, error) {
	num, err := rpc.rpcImpl.EthBlockNumber()
	return uint64(num), err
}

func (rpc RPC) GetCurrentBlockNum() (rst uint64, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.getCurrentBlockNum()
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}

func (rpc RPC) GetBlockByNum(num uint64) (rst Block, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.getBlockByNum(num, true)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}

func (rpc RPC) getBlockByNum(num uint64, withTx bool) (Block, error) {
	b, err := rpc.rpcImpl.EthGetBlockByNumber(int(num), withTx)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("nil block")
	}

	return &EthereumBlock{b}, err
}

func (rpc RPC) GetLiteBlockByNum(num uint64) (rst Block, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.getBlockByNum(num, false)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}

func (rpc RPC) getTransactionReceipt(txHash string) (TransactionReceipt, error) {
	receipt, err := rpc.rpcImpl.EthGetTransactionReceipt(txHash)
	if err != nil {
		return nil, err
	}
	if receipt == nil {
		return nil, errors.New("nil receipt")
	}

	return &EthereumTransactionReceipt{receipt}, err
}

func (rpc RPC) GetTransactionReceipt(txHash string) (rst TransactionReceipt, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.getTransactionReceipt(txHash)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}

func (rpc RPC) getLogs(
	fromBlockNum, toBlockNum uint64,
	addresses []string,
	topics []string,
) ([]IReceiptLog, error) {

	filterParam := ethrpc.FilterParams{
		FromBlock: "0x" + strconv.FormatUint(fromBlockNum, 16),
		ToBlock:   "0x" + strconv.FormatUint(toBlockNum, 16),
		Address:   addresses,
		Topics:    [][]string{topics},
	}

	logs, err := rpc.rpcImpl.EthGetLogs(filterParam)
	if err != nil {
		logrus.Errorf("eth_getlogs: failed to retrieve logs: %v", err)
		return nil, err
	}

	logrus.Tracef("eth_getlogs: log count at block(%d - %d): %d", fromBlockNum, toBlockNum, len(logs))

	var result []IReceiptLog
	for i := 0; i < len(logs); i++ {
		l := logs[i]
		result = append(result, ReceiptLog{Log: &l})

		logrus.Tracef("eth_getlogs: receipt log: %+v", l)
	}

	return result, err
}

func (rpc RPC) GetLogs(
	fromBlockNum, toBlockNum uint64,
	addresses []string,
	topics []string,
) (rst []IReceiptLog, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.getLogs(fromBlockNum, toBlockNum, addresses, topics)
		if err == nil {
			break
		} else {
			time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
		}
	}

	return
}
