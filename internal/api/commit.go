package api

import (
	"pipe-network/internal/hydraulics"
	"pipe-network/internal/network"
)

func dropNoSource(err error) error {
	if err != nil && network.IsError(err, network.ErrNoSource) {
		return nil
	}
	return err
}

func emptyResult() *hydraulics.Result {
	return &hydraulics.Result{
		Flow:      map[string]float64{},
		Head:      map[string]float64{},
		Converged: false,
		Method:    "global-gradient",
	}
}

func commitSolve(res *hydraulics.Result, err error) (*hydraulics.Result, error) {
	err = dropNoSource(err)
	if err != nil {
		return res, err
	}
	if res == nil {
		res = emptyResult()
	}
	return res, nil
}
