package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// -----------------------------
type FailoverBackend struct {
    *ethclient.Client
    mgr *ClientManager
}

var _ bind.ContractBackend = (*FailoverBackend)(nil)
func (f *FailoverBackend) CallContract(
	ctx context.Context,
	call ethereum.CallMsg,
	blockNumber *big.Int,
) ([]byte, error) {

	client := f.mgr.Current()
	res, err := client.CallContract(ctx, call, blockNumber)
	if err == nil {
		return res, nil
	}

	if IsRetryableRPCError(err) {
		f.mgr.Rotate()
		return f.mgr.Current().CallContract(ctx, call, blockNumber)
	}

	return nil, err
}

func (f *FailoverBackend) CodeAt(
	ctx context.Context,
	account common.Address,
	blockNumber *big.Int,
) ([]byte, error) {

	client := f.mgr.Current()
	code, err := client.CodeAt(ctx, account, blockNumber)
	if err == nil {
		return code, nil
	}

	if IsRetryableRPCError(err) {
		f.mgr.Rotate()
		return f.mgr.Current().CodeAt(ctx, account, blockNumber)
	}

	return nil, err
}

func (f *FailoverBackend) PendingNonceAt(
	ctx context.Context,
	account common.Address,
) (uint64, error) {

	client := f.mgr.Current()
	nonce, err := client.PendingNonceAt(ctx, account)
	if err == nil {
		return nonce, nil
	}

	if IsRetryableRPCError(err) {
		f.mgr.Rotate()
		return f.mgr.Current().PendingNonceAt(ctx, account)
	}

	return 0, err
}

func (f *FailoverBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	client := f.mgr.Current()
	price, err := client.SuggestGasPrice(ctx)
	if err == nil {
		return price, nil
	}

	if IsRetryableRPCError(err) {
		f.mgr.Rotate()
		return f.mgr.Current().SuggestGasPrice(ctx)
	}

	return nil, err
}

func (f *FailoverBackend) EstimateGas(
	ctx context.Context,
	call ethereum.CallMsg,
) (uint64, error) {

	client := f.mgr.Current()
	gas, err := client.EstimateGas(ctx, call)
	if err == nil {
		return gas, nil
	}

	if IsRetryableRPCError(err) {
		f.mgr.Rotate()
		return f.mgr.Current().EstimateGas(ctx, call)
	}

	return 0, err
}

func (f *FailoverBackend) SendTransaction(
	ctx context.Context,
	tx *types.Transaction,
) error {

	client := f.mgr.Current()
	err := client.SendTransaction(ctx, tx)
	if err == nil {
		return nil
	}

	if IsRetryableRPCError(err) {
		f.mgr.Rotate()
		return f.mgr.Current().SendTransaction(ctx, tx)
	}

	return err
}

func (f *FailoverBackend) ChainID(ctx context.Context) (*big.Int, error) {
	client := f.mgr.Current()
	id, err := client.ChainID(ctx)
	if err == nil {
		return id, nil
	}

	if IsRetryableRPCError(err) {
		f.mgr.Rotate()
		return f.mgr.Current().ChainID(ctx)
	}

	return nil, err
}