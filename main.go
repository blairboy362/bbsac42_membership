package main

import (
	"fmt"
	"os"
	"time"

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

func newMembership(fc *fileConfig) (*membership, error) {
	var err error
	m := membership{}
	m.references, err = loadMemberReferencesFromCsv(fc.getSourcePath("reference_member_mappings.csv"))
	if err != nil {
		return nil, err
	}

	m.members, err = loadMembershipDetailsFromCsv(fc.getSourcePath("membership_details.csv"))
	if err != nil {
		return nil, err
	}

	m.newMembers, err = loadNewMembersFromCsv(fc.getSourcePath("new_members.csv"))
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
	if len(os.Args) < 2 {
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
	currentUtcTime := time.Now().UTC()
	fileConfig := fileConfig{
		os.Args[1],
		currentUtcTime.Format("200601"),
		currentUtcTime.AddDate(0, -1, 0).Format("200601"),
	}

	ignoreTxnsPath := fileConfig.getCurrentDestinationPath(DefaultIgnoreTxnsPath)
	incorrectMembershipTxnsPath := fileConfig.getCurrentDestinationPath(DefaultIncorrectMembershipTxnsPath)
	unmatchedTxnsPath := fileConfig.getCurrentDestinationPath(DefaultUnmatchedTxnsPath)
	paidMembersPath := fileConfig.getCurrentDestinationPath(DefaultPaidMembersPath)
	unmatchedMemberIDsPath := fileConfig.getCurrentDestinationPath(DefaultUnmatchedMemberIDsPath)
	allMembersPath := fileConfig.getCurrentDestinationPath(DefaultAllMembersPath)

	membership, err := newMembership(&fileConfig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %v references.\n", len(membership.references))
	fmt.Printf("Loaded %v members details.\n", len(membership.members))
	fmt.Printf("Loaded %v new members details.\n", len(membership.newMembers))
	activeMembers, err := membership.loadAndFilterTxns(fileConfig.getCurrentSourcePath("bank_acct_txns.csv"), correctMembershipAmounts, membershipAmounts)
	if err != nil {
		panic(err)
	}

	destinationDir := fileConfig.getCurrentDestinationPath("")
	fmt.Printf("Deleting directory %v.\n", destinationDir)
	err = os.RemoveAll(destinationDir)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Creating directory %v.\n", destinationDir)
	err = os.MkdirAll(destinationDir, 0755)
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
