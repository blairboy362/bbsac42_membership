package main

import (
	"fmt"
	"os"

	"github.com/shopspring/decimal"
)

const (
	DefaultIgnoreTxnsPath              = "ignored_txns.csv"
	DefaultIncorrectMembershipTxnsPath = "incorrect_membership_txns.csv"
	DefaultUnmatchedTxnsPath           = "unmatched_txns.csv"
	DefaultPaidMembersPath             = "paid_members.csv"
	DefaultUnmatchedMemberIDsPath      = "unmatched_memberids.csv"
	DefaultAllMembersPath              = "all_members.csv"
)

type membership struct {
	references map[string][]string
	members    map[string]*Member
	newMembers []*Member
}

type transactions struct {
	ignored   []*bankTxn
	incorrect []*bankTxn
	candidate []*bankTxn
	unmatched []*bankTxn
}

type members struct {
	paying       []*Member
	unmatchedIDs []string
}

type activeMembers struct {
	txns    transactions
	members members
}

func newMembership(referencePath, membersPath, newMembersPath string) (*membership, error) {
	var err error
	m := membership{}
	m.references, err = loadMemberReferencesFromCsv(referencePath)
	if err != nil {
		return nil, err
	}

	m.members, err = loadMembershipDetailsFromCsv(membersPath)
	if err != nil {
		return nil, err
	}

	m.newMembers, err = loadNewMembersFromCsv(newMembersPath)
	if err != nil {
		panic(err)
	}

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
		fmt.Printf("Writing %v ignored records to %v.\n", len(activeMembers.txns.ignored), DefaultIgnoreTxnsPath)
		err = writeTxnsToCsv(DefaultIgnoreTxnsPath, activeMembers.txns.ignored)
		if err != nil {
			panic(err)
		}
	}

	if len(activeMembers.txns.incorrect) > 0 {
		fmt.Printf("Writing %v incorrect membership records to %v.\n", len(activeMembers.txns.incorrect), DefaultIncorrectMembershipTxnsPath)
		err = writeTxnsToCsv(DefaultIncorrectMembershipTxnsPath, activeMembers.txns.incorrect)
		if err != nil {
			panic(err)
		}
	}

	if len(activeMembers.txns.unmatched) > 0 {
		fmt.Printf("Writing %v unmatched transactions to %v.\n", len(activeMembers.txns.unmatched), DefaultUnmatchedTxnsPath)
		err = writeTxnsToCsv(DefaultUnmatchedTxnsPath, activeMembers.txns.unmatched)
		if err != nil {
			panic(err)
		}
	}

	if len(activeMembers.members.unmatchedIDs) > 0 {
		fmt.Printf("Writing %v unmatched member IDs to %v.\n", len(activeMembers.members.unmatchedIDs), DefaultUnmatchedMemberIDsPath)
		err = writeMemberIdsToCsv(DefaultUnmatchedMemberIDsPath, activeMembers.members.unmatchedIDs)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Writing %v paid members details to %v\n", len(activeMembers.members.paying), DefaultPaidMembersPath)
	err = writeMembersToCsv(DefaultPaidMembersPath, activeMembers.members.paying)
	if err != nil {
		panic(err)
	}

	allMembers := append(activeMembers.members.paying, membership.newMembers...)
	fmt.Printf("Writing %v members details to %v\n", len(allMembers), DefaultAllMembersPath)
	err = writeMembersToCsv(DefaultAllMembersPath, allMembers)
	if err != nil {
		panic(err)
	}
}
