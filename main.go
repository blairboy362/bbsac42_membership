package main

import (
	"fmt"
	"os"

	"github.com/shopspring/decimal"
)

type membership struct {
	references []referenceLookup
	members    map[string]member
	newMembers []member
}

type transactions struct {
	ignored   []bankTxn
	incorrect []bankTxn
	candidate []bankTxn
	unmatched []bankTxn
}

type members struct {
	paying       []member
	unmatchedIDs []string
}

type activeMembers struct {
	txns    transactions
	members members
}

func newMembership(referencePath, membersPath, newMembersPath string) (*membership, error) {
	m := membership{}
	references, err := loadMemberReferencesFromCsv(referencePath)
	if err != nil {
		return nil, err
	}

	m.references = references
	members, err := loadMembershipDetailsFromCsv(membersPath)
	if err != nil {
		return nil, err
	}

	m.members = members
	newMembers, err := loadNewMembersFromCsv(newMembersPath)
	if err != nil {
		panic(err)
	}

	m.newMembers = newMembers

	return &m, nil
}

func (m *membership) loadAndFilterTxns(txnsPath string, correctMembershipAmounts, otherMembershipAmounts []decimal.Decimal) (am *activeMembers, err error) {
	txns, err := loadTxnsFromCsv(txnsPath)
	if err != nil {
		return nil, err
	}

	var transactions = transactions{}
	membershipAmounts := append(correctMembershipAmounts, otherMembershipAmounts...)
	transactions.candidate, transactions.ignored = filterInterestingTxns(txns, membershipAmounts)
	_, transactions.incorrect = filterInterestingTxns(transactions.candidate, correctMembershipAmounts)
	var memberIds []string
	memberIds, transactions.unmatched = identifyMembers(transactions.candidate, m.references)
	var members = members{}
	members.paying, members.unmatchedIDs = matchMembers(m.members, memberIds)

	am = new(activeMembers)
	am.txns = transactions
	am.members = members

	return am, err
}

func main() {
	if len(os.Args) != 5 {
		panic("Missing argument!")
	}

	correctMembershipAmounts := []decimal.Decimal{
		decimal.New(185, -1),
		decimal.New(30, 0),
	}
	membershipAmounts := []decimal.Decimal{
		decimal.New(165, -1),
		decimal.New(15, 0),
		decimal.New(25, 0),
	}
	membershipAmounts = append(membershipAmounts, correctMembershipAmounts...)
	ignoreTxnsPath := "ignored_txns.csv"
	incorrectMembershipTxnsPath := "incorrect_membership_txns.csv"
	unmatchedTxnsPath := "unmatched_txns.csv"
	paidMembersPath := "paid_members.csv"
	unmatchedMemberIDsPath := "unmatched_memberids.csv"
	allMembersPath := "all_members.csv"
	txnsPath := os.Args[1]
	referencePath := os.Args[2]
	membersPath := os.Args[3]
	newMembersPath := os.Args[4]

	membership, err := newMembership(referencePath, membersPath, newMembersPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %v references.\n", len(membership.references))
	fmt.Printf("Loaded %v members details.\n", len(membership.members))
	fmt.Printf("Loaded %v new members details.\n", len(membership.newMembers))
	activeMembers, err := membership.loadAndFilterTxns(txnsPath, correctMembershipAmounts, membershipAmounts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %v transactions.\n", len(activeMembers.txns.candidate)+len(activeMembers.txns.ignored))
	if len(activeMembers.txns.ignored) > 0 {
		fmt.Printf("Writing %v ignored records to %v.\n", len(activeMembers.txns.ignored), ignoreTxnsPath)
		err = writeTxnsToCsv(ignoreTxnsPath, activeMembers.txns.ignored)
		if err != nil {
			panic(err)
		}
	}

	if len(activeMembers.txns.incorrect) > 0 {
		fmt.Printf("Writing %v incorrect membership records to %v.\n", len(activeMembers.txns.incorrect), incorrectMembershipTxnsPath)
		err = writeTxnsToCsv(incorrectMembershipTxnsPath, activeMembers.txns.incorrect)
		if err != nil {
			panic(err)
		}
	}

	if len(activeMembers.txns.unmatched) > 0 {
		fmt.Printf("Writing %v unmatched transactions to %v.\n", len(activeMembers.txns.unmatched), unmatchedTxnsPath)
		err = writeTxnsToCsv(unmatchedTxnsPath, activeMembers.txns.unmatched)
		if err != nil {
			panic(err)
		}
	}

	if len(activeMembers.members.unmatchedIDs) > 0 {
		fmt.Printf("Writing %v unmatched member IDs to %v.\n", len(activeMembers.members.unmatchedIDs), unmatchedMemberIDsPath)
		err = writeMemberIdsToCsv(unmatchedMemberIDsPath, activeMembers.members.unmatchedIDs)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Writing %v paid members details to %v\n", len(activeMembers.members.paying), paidMembersPath)
	err = writeMembersToCsv(paidMembersPath, activeMembers.members.paying)
	if err != nil {
		panic(err)
	}

	allMembers := append(activeMembers.members.paying, membership.newMembers...)
	fmt.Printf("Writing %v members details to %v\n", len(allMembers), allMembersPath)
	err = writeMembersToCsv(allMembersPath, allMembers)
	if err != nil {
		panic(err)
	}
}
