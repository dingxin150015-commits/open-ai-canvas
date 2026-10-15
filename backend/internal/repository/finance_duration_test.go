package repository

import (
	"errors"
	"yingce/backend/internal/model"
	"testing"
)

func TestDurationBillingSettlementAndRecovery(t *testing.T) {
	for _, mode := range []string{"per_second", "per_second_usage"} {
		for _, recovered := range []bool{false, true} {
			t.Run(mode+map[bool]string{false: "settle", true: "recover"}[recovered], func(t *testing.T) {
				repo, db := newRefundedBillingRecoveryRepository(t)
				status := model.BillingStatusReserved
				reserved := int64(10000000)
				if recovered {
					status = model.BillingStatusRefunded
					reserved = 0
				}
				if err := db.Create(&model.CreditAccount{UserID: "user", AvailableMicrocredits: 90000000, ReservedMicrocredits: reserved}).Error; err != nil {
					t.Fatal(err)
				}
				order := model.BillingOrder{ID: "order", UserID: "user", IdempotencyKey: "key", Status: status, Capability: "video", BillingMode: mode, UnitPriceMicrocredits: 1000000, MultiplierBasisPoints: 10000, AmountMicrocredits: 10000000, ReservedAmountMicrocredits: 10000000}
				if err := db.Create(&order).Error; err != nil {
					t.Fatal(err)
				}
				if err := db.Create(&model.ApiCallLog{ID: "usage", BillingOrderID: "order", Status: model.ApiCallStatusSucceeded, RequestKind: "poll", UsageAvailable: true, UsageDurationMs: 8250}).Error; err != nil {
					t.Fatal(err)
				}
				settle := func() error {
					if recovered {
						return repo.RestoreRefundedBillingOrder("order", "provider")
					}
					return repo.SettleBillingOrder("order", "provider")
				}
				if err := settle(); err != nil {
					t.Fatal(err)
				}
				if err := settle(); err != nil {
					t.Fatal(err)
				}
				if err := db.First(&order, "id = ?", "order").Error; err != nil {
					t.Fatal(err)
				}
				expected := int64(10000000)
				if mode == "per_second_usage" {
					expected = 8250000
				}
				if order.ActualAmountMicrocredits != expected {
					t.Fatalf("actual=%d want=%d", order.ActualAmountMicrocredits, expected)
				}
				var count int64
				db.Model(&model.CreditLedgerEntry{}).Where("billing_order_id = ?", "order").Count(&count)
				// Settlement may add both consume/refund ledgers; rerunning must not duplicate them.
				if count < 1 || count > 2 {
					t.Fatalf("ledger count=%d", count)
				}
			})
		}
	}
}

func TestDurationBillingRequiresObservedUsage(t *testing.T) {
	order := model.BillingOrder{UnitPriceMicrocredits: 1000000, MultiplierBasisPoints: 10000}
	if _, err := durationUsageAmount(order, nil); !errors.Is(err, ErrBillingUsageUnavailable) {
		t.Fatal(err)
	}
	for _, ms := range []int64{1, 8250, 15000} {
		amount, err := durationUsageAmount(order, &BillingUsage{DurationMs: ms})
		if err != nil || amount != ms*1000 {
			t.Fatalf("ms=%d amount=%d err=%v", ms, amount, err)
		}
	}
}

func TestDurationBillingMissingUsageDoesNotSpendReservation(t *testing.T) {
	repo, db := newRefundedBillingRecoveryRepository(t)
	if err := db.Create(&model.CreditAccount{UserID: "user", AvailableMicrocredits: 90000000, ReservedMicrocredits: 10000000}).Error; err != nil {
		t.Fatal(err)
	}
	order := model.BillingOrder{ID: "order", UserID: "user", IdempotencyKey: "missing", Status: model.BillingStatusReserved, Capability: "video", BillingMode: "per_second_usage", UnitPriceMicrocredits: 1000000, MultiplierBasisPoints: 10000, AmountMicrocredits: 10000000, ReservedAmountMicrocredits: 10000000}
	if err := db.Create(&order).Error; err != nil {
		t.Fatal(err)
	}
	if err := repo.SettleBillingOrder("order", "provider"); !errors.Is(err, ErrBillingUsageUnavailable) {
		t.Fatalf("expected missing usage: %v", err)
	}
	var account model.CreditAccount
	if err := db.First(&account, "user_id = ?", "user").Error; err != nil {
		t.Fatal(err)
	}
	if account.AvailableMicrocredits != 90000000 || account.ReservedMicrocredits != 10000000 {
		t.Fatalf("reservation changed: %#v", account)
	}
	if err := db.First(&order, "id = ?", "order").Error; err != nil {
		t.Fatal(err)
	}
	if order.Status != model.BillingStatusReserved {
		t.Fatal("missing usage settled")
	}
}

func TestDurationBillingExcessUsageCannotOverdrawAccount(t *testing.T) {
	for _, restored := range []bool{false, true} {
		repo, db := newRefundedBillingRecoveryRepository(t)
		status := model.BillingStatusReserved
		if restored {
			status = model.BillingStatusRefunded
		}
		if err := db.Create(&model.CreditAccount{UserID: "user", AvailableMicrocredits: 0, ReservedMicrocredits: 1000000}).Error; err != nil {
			t.Fatal(err)
		}
		order := model.BillingOrder{ID: "order", UserID: "user", IdempotencyKey: "cap", Status: status, Capability: "video", BillingMode: "per_second_usage", UnitPriceMicrocredits: 1000000, MultiplierBasisPoints: 10000, AmountMicrocredits: 1000000, ReservedAmountMicrocredits: 1000000}
		if err := db.Create(&order).Error; err != nil {
			t.Fatal(err)
		}
		if err := db.Create(&model.ApiCallLog{ID: "usage", BillingOrderID: "order", Status: model.ApiCallStatusSucceeded, UsageAvailable: true, UsageDurationMs: 8250}).Error; err != nil {
			t.Fatal(err)
		}
		var err error
		if restored {
			err = repo.RestoreRefundedBillingOrder("order", "provider")
		} else {
			err = repo.SettleBillingOrder("order", "provider")
		}
		if !errors.Is(err, ErrDurationReservationExceeded) {
			t.Fatalf("expected spending cap: %v", err)
		}
		var account model.CreditAccount
		if err := db.First(&account, "user_id = ?", "user").Error; err != nil {
			t.Fatal(err)
		}
		if account.AvailableMicrocredits != 0 || account.ReservedMicrocredits != 1000000 {
			t.Fatalf("balance changed: %#v", account)
		}
		var count int64
		db.Model(&model.CreditLedgerEntry{}).Count(&count)
		if count != 0 {
			t.Fatalf("unexpected ledger count: %d", count)
		}
		if !restored {
			if err := repo.SettleReviewedDurationBillingOrder("order", "provider"); !errors.Is(err, ErrInsufficientCredits) {
				t.Fatalf("review overdraw allowed: %v", err)
			}
			if err := db.Model(&model.CreditAccount{}).Where("user_id = ?", "user").Update("available_microcredits", 10000000).Error; err != nil {
				t.Fatal(err)
			}
			if err := repo.SettleReviewedDurationBillingOrder("order", "provider"); err != nil {
				t.Fatal(err)
			}
			if err := repo.SettleReviewedDurationBillingOrder("order", "provider"); err != nil {
				t.Fatal(err)
			}
			if err := db.First(&account, "user_id = ?", "user").Error; err != nil {
				t.Fatal(err)
			}
			if account.AvailableMicrocredits != 2750000 || account.ReservedMicrocredits != 0 {
				t.Fatalf("reviewed balance wrong: %#v", account)
			}
		}
	}
}
