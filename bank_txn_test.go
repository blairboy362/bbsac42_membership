package main

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestFilterInterestingTxns(t *testing.T) {
	allTxns := []*bankTxn{
		{"JoeBloggs", decimal.New(15, 1)},
		{"JaneDoe", decimal.New(1234, -2)},
	}
	interestedAmounts := []decimal.Decimal{
		decimal.New(15, 1),
		decimal.New(165, -1),
		decimal.New(25, 1),
		decimal.New(30, 1),
	}
	expectedInterestingTxns := []*bankTxn{
		{"JoeBloggs", decimal.New(15, 1)},
	}
	expectedJunkTxns := []*bankTxn{
		{"JaneDoe", decimal.New(1234, -2)},
	}

	actualInterestingTxns, actualJunkTxns := filterInterestingTxns(allTxns, interestedAmounts)

	if actualInterestingTxns == nil || actualJunkTxns == nil {
		t.Fatalf("One or more returns are nil!")
	}

	if len(actualInterestingTxns) != len(expectedInterestingTxns) {
		t.Fatalf(
			"Interesting txn counts are not the same (%v, %v)",
			len(actualInterestingTxns),
			len(expectedInterestingTxns),
		)
	}

	if len(actualJunkTxns) != len(expectedJunkTxns) {
		t.Fatalf(
			"Junk txn counts are not the same (%v, %v)",
			len(actualJunkTxns),
			len(expectedJunkTxns),
		)
	}

	for i := range expectedInterestingTxns {
		if !expectedInterestingTxns[i].equal(actualInterestingTxns[i]) {
			t.Fatalf(
				"%v != %v",
				expectedInterestingTxns[i],
				actualInterestingTxns[i])
		}
	}

	for i := range expectedJunkTxns {
		if !expectedJunkTxns[i].equal(actualJunkTxns[i]) {
			t.Fatalf(
				"%v != %v",
				expectedJunkTxns[i],
				actualJunkTxns[i])
		}
	}
}
