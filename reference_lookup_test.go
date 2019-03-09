package main

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestIdentifyMembers(t *testing.T) {
	txns := []*bankTxn{
		{"JOE BLOGGS", decimal.New(15, 1)},
		{"NO MATCH", decimal.New(30, 1)},
		{"JOE BLOGGS", decimal.New(15, 1)},
	}
	refs := map[string][]string{
		"JOE BLOGGS": {"A123456", "A789012"},
	}
	expectedMemberIds := []string{
		"A123456",
		"A789012",
	}
	expectedUnmatchedTxns := []*bankTxn{
		{"NO MATCH", decimal.New(30, 1)},
	}

	actualMemberIds, actualUnmatchedTxns := identifyMembers(txns, refs)

	if actualMemberIds == nil || actualUnmatchedTxns == nil {
		t.Fatalf("One or more return values are nil!")
	}

	if len(actualMemberIds) != len(expectedMemberIds) {
		t.Fatalf(
			"MemberId lengths don't match (%v, %v)",
			len(actualMemberIds),
			len(expectedMemberIds),
		)
	}

	if len(actualUnmatchedTxns) != len(expectedUnmatchedTxns) {
		t.Fatalf(
			"UnmatchedTxn lengths don't match (%v, %v)",
			len(actualUnmatchedTxns),
			len(expectedUnmatchedTxns),
		)
	}

	for i := range expectedMemberIds {
		if expectedMemberIds[i] != actualMemberIds[i] {
			t.Fatalf("%v != %v", expectedMemberIds[i], actualMemberIds[i])
		}
	}

	for i := range expectedUnmatchedTxns {
		if !expectedUnmatchedTxns[i].equal(actualUnmatchedTxns[i]) {
			t.Fatalf("%v != %v", expectedUnmatchedTxns[i], actualUnmatchedTxns[i])
		}
	}
}
