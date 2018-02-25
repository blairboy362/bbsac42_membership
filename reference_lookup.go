package main

import (
	"fmt"
	"os"
)

/*
This file used to include a struct with additional boilerplate. Keeping
this here as I don't currently have a better place for it.
*/
func identifyMembers(txns []bankTxn, references map[string][]string) (memberIds []string, unmatchedTxns []bankTxn) {
	memberSet := make(map[string]bool)
	for _, txn := range txns {
		memberIDs, ok := references[txn.description]
		if ok {
			for _, memberID := range memberIDs {
				_, ok := memberSet[memberID]
				if ok {
					fmt.Fprintf(os.Stderr, "Duplicate member ID: %v", memberID)
				}
				memberSet[memberID] = true
			}
		} else {
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
