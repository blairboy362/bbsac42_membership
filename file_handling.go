package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

func loadMembershipDetailsFromCsv(path string) (members map[string]*Member, err error) {
	memberFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open %s", path)
	}
	defer memberFile.Close()

	loadedMembers := []*Member{}
	err = gocsv.UnmarshalFile(memberFile, &loadedMembers)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse membership from %s", path)
	}

	members = make(map[string]*Member)
	for _, member := range loadedMembers {
		members[member.MemberID] = member
	}

	return members, nil
}

func loadNewMembersFromCsv(path string) (members []*Member, err error) {
	memberFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open %s", path)
	}
	defer memberFile.Close()

	err = gocsv.UnmarshalFile(memberFile, &members)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse membership from %s", path)
	}

	return members, nil
}

type CsvBankTxn struct {
	Junk1       string `csv:"Date"`
	Junk2       string `csv:"Type"`
	Description string `csv:"Description"`
	Junk3       string `csv:"Paid out"`
	Amount      string `csv:"Paid in"`
	Junk4       string `csv:"Balance"`
}

func loadTxnsFromCsv(path string) (txns []*bankTxn, err error) {
	txnFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open %s", path)
	}
	defer txnFile.Close()

	loadedTxns := []*CsvBankTxn{}
	err = gocsv.UnmarshalFile(txnFile, &loadedTxns)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse transactions from %s", path)
	}

	for _, loadedTxn := range loadedTxns {
		txn, err := newBankTxn(loadedTxn.Description, loadedTxn.Amount)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse transaction (%v): %v", *loadedTxn, err)
		}

		txns = append(txns, txn)
	}

	return txns, nil
}

type MemberReference struct {
	Reference string `csv:"Reference"`
	MemberIDs string `csv:"MemberIds"`
}

func loadMemberReferencesFromCsv(path string) (references map[string][]string, err error) {
	referenceFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open %s", path)
	}
	defer referenceFile.Close()

	loadedReferences := []*MemberReference{}
	err = gocsv.UnmarshalFile(referenceFile, &loadedReferences)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse membership from %s", path)
	}

	references = map[string][]string{}
	for _, loadedReference := range loadedReferences {
		reference := strings.TrimSpace(strings.ToUpper(loadedReference.Reference))
		membershipIds := strings.Split(loadedReference.MemberIDs, "|")
		_, ok := references[reference]
		if !ok {
			references[reference] = membershipIds
		} else {
			fmt.Fprintf(os.Stderr, "Loaded duplicate reference %v\n", reference)
		}
	}

	return references, nil
}

func writeTxnsToCsv(path string, txns []*bankTxn) error {
	txnFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer txnFile.Close()

	csvWriter := csv.NewWriter(bufio.NewWriter(txnFile))
	err = csvWriter.Write([]string{"Description", "Amount"})
	if err != nil {
		return err
	}

	for _, txn := range txns {
		err := csvWriter.Write([]string{txn.description, txn.amount.String()})
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}

func writeMembersToCsv(path string, members []*Member) error {
	memberFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer memberFile.Close()

	err = gocsv.MarshalFile(members, memberFile)
	if err != nil {
		return err
	}

	return nil
}

func writeMemberIdsToCsv(path string, memberIds []string) error {
	memberIDFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer memberIDFile.Close()

	csvWriter := csv.NewWriter(bufio.NewWriter(memberIDFile))
	err = csvWriter.Write([]string{"MemberId"})
	if err != nil {
		return err
	}

	for _, memberID := range memberIds {
		err := csvWriter.Write([]string{memberID})
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}
