package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

type fileConfig struct {
	baseDir            string
	currentFolderName  string
	previousFolderName string
}

func (fc *fileConfig) getSourcePath(fileName string) string {
	return filepath.Join(fc.baseDir, "in", fileName)
}

func (fc *fileConfig) getCurrentSourcePath(fileName string) string {
	return filepath.Join(fc.baseDir, "in", fc.currentFolderName, fileName)
}

func (fc *fileConfig) getCurrentDestinationPath(fileName string) string {
	return filepath.Join(fc.baseDir, "out", fc.currentFolderName, fileName)
}

func (fc *fileConfig) getPreviousDestinationPath(fileName string) string {
	return filepath.Join(fc.baseDir, "out", fc.previousFolderName, fileName)
}

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

func loadAllMembersFromCsv(path string) (members []*Member, err error) {
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
	TxnType     string `csv:"Type"`
	Description string `csv:"Description"`
	Junk2       string `csv:"Paid out"`
	Amount      string `csv:"Paid in"`
	Junk3       string `csv:"Balance"`
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
		if loadedTxn.TxnType == "CR" {
			txn, err := newBankTxn(loadedTxn.Description, loadedTxn.Amount)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse transaction (%v): %v", *loadedTxn, err)
			}

			txns = append(txns, txn)
		}
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

func loadEmailsFromCsv(path string) (emails []string, err error) {
	emailsFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer emailsFile.Close()

	emails = []string{}
	csvReader := csv.NewReader(bufio.NewReader(emailsFile))
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		emails = append(emails, line[0])
	}

	return emails, nil
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
	targetFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	csvWriter := csv.NewWriter(bufio.NewWriter(targetFile))
	err = csvWriter.Write([]string{"MemberId"})
	if err != nil {
		return err
	}

	for _, line := range memberIds {
		err := csvWriter.Write([]string{line})
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}

func writeStringsToCsv(path string, content []string) error {
	targetFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	csvWriter := csv.NewWriter(bufio.NewWriter(targetFile))
	for _, line := range content {
		err := csvWriter.Write([]string{line})
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}
