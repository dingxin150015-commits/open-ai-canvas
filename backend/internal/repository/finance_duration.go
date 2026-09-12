package repository

import (
	"errors"
	"infinite-canvas/backend/internal/model"
	"math/big"
)

func durationUsageAmount(order model.BillingOrder, usage *BillingUsage) (int64, error) {
	if usage == nil || usage.DurationMs <= 0 {
		return 0, ErrBillingUsageUnavailable
	}
	if usage.DurationMs > 3600000 || order.UnitPriceMicrocredits < 0 || order.MultiplierBasisPoints <= 0 {
		return 0, errors.New("invalid duration usage")
	}
	amount := new(big.Int).Mul(big.NewInt(usage.DurationMs), big.NewInt(order.UnitPriceMicrocredits))
	amount.Mul(amount, big.NewInt(order.MultiplierBasisPoints))
	amount.Add(amount, big.NewInt(9999999)).Quo(amount, big.NewInt(10000000))
	if !amount.IsInt64() {
		return 0, errors.New("duration amount overflow")
	}
	return amount.Int64(), nil
}
