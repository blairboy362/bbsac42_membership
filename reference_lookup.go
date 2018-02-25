package main

import (
	"fmt"
	"os"
	"strings"
)

type referenceLookup struct {
	reference     string
	membershipIds []string
}

func newReferenceLookup(reference string, membershipIds []string) *referenceLookup {
	return &referenceLookup{
		strings.TrimSpace(strings.ToUpper(reference)),
		membershipIds,
	}
}

func (rl1 *referenceLookup) equal(rl2 *referenceLookup) bool {
	if rl1.reference != rl2.reference {
		return false
	}

	if len(rl1.membershipIds) != len(rl2.membershipIds) {
		return false
	}

	for i := range rl1.membershipIds {
		if rl1.membershipIds[i] != rl2.membershipIds[i] {
			return false
		}
	}

	return true
}

func identifyMembers(txns []bankTxn, references []referenceLookup) (memberIds []string, unmatchedTxns []bankTxn) {
	memberSet := make(map[string]bool)
	for _, txn := range txns {
		found := false
		for _, lookup := range references {
			if txn.description == lookup.reference {
				for _, memberID := range lookup.membershipIds {
					_, ok := memberSet[memberID]
					if ok {
						fmt.Fprintf(os.Stderr, "Duplicate member ID: %v", memberID)
					}
					memberSet[memberID] = true
				}

				found = true
				break
			}
		}

		if !found {
			unmatchedTxns = append(unmatchedTxns, txn)
		}
	}

	memberIds = make([]string, len(memberSet))
	i := 0
	for k := range memberSet {
		memberIds[i] = k
		i++
	}

	return memberIds, unmatchedTxns
}
