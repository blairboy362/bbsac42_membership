package main

import (
	"strings"

	"github.com/shopspring/decimal"
)

type bankTxn struct {
	description string
	amount      decimal.Decimal
}

func newBankTxn(description, amount string) (*bankTxn, error) {
	amt, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, err
	}

	return &bankTxn{strings.ToUpper(strings.TrimSpace(description)), amt}, nil
}

func (t1 *bankTxn) equal(t2 *bankTxn) bool {
	return t1.description == t2.description && t1.amount.Equal(t2.amount)
}

func filterInterestingTxns(source []bankTxn, interestingAmounts []decimal.Decimal) (interestedTxns []bankTxn, junkTxns []bankTxn) {
	for _, txn := range source {
		found := false
		for _, amount := range interestingAmounts {
			if amount.Equal(txn.amount) {
				interestedTxns = append(interestedTxns, txn)
				found = true
				break
			}
		}

		if !found {
			junkTxns = append(junkTxns, txn)
		}
	}

	return interestedTxns, junkTxns
}
