package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func loadMembershipDetailsFromCsv(path string) (members map[string]member, err error) {
	memberFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer memberFile.Close()

	members = make(map[string]member)
	var csvFields map[string]int
	csvFields = make(map[string]int)
	csvReader := csv.NewReader(bufio.NewReader(memberFile))
	firstLine := true
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if len(line) < 5 {
			return nil, fmt.Errorf("Line does not have enough fields: %v", line)
		}

		if firstLine {
			firstLine = false
			for i, fieldName := range line {
				csvFields[fieldName] = i
			}
			continue
		}

		members[line[csvFields["MemberId"]]] = member{
			line[csvFields["MemberId"]],
			line[csvFields["Title"]],
			line[csvFields["Forenames"]],
			line[csvFields["Surname"]],
			line[csvFields["EmailAddress"]],
		}
	}

	return members, nil
}

func loadNewMembersFromCsv(path string) (members []member, err error) {
	memberFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer memberFile.Close()

	var csvFields map[string]int
	csvFields = make(map[string]int)
	csvReader := csv.NewReader(bufio.NewReader(memberFile))
	firstLine := true
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if len(line) < 4 {
			return nil, fmt.Errorf("Line does not have enough fields: %v", line)
		}

		if firstLine {
			firstLine = false
			for i, fieldName := range line {
				csvFields[fieldName] = i
			}
			continue
		}

		members = append(members, member{
			"",
			line[csvFields["Title"]],
			line[csvFields["Forenames"]],
			line[csvFields["Surname"]],
			line[csvFields["EmailAddress"]],
		})
	}

	return members, nil
}

func loadTxnsFromCsv(path string) (txns []bankTxn, err error) {
	txnFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer txnFile.Close()

	csvReader := csv.NewReader(bufio.NewReader(txnFile))
	firstLine := true
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if len(line) < 5 {
			return nil, fmt.Errorf("Line %v does not have enough fields", line)
		}

		if firstLine {
			firstLine = false
			continue
		}

		if line[1] == "CR" {
			txn, err := newBankTxn(line[2], line[4])
			if err != nil {
				return nil, fmt.Errorf("Failed to parse transaction (%v): %v", line, err)
			}

			txns = append(txns, *txn)
		} else {
			fmt.Fprintf(os.Stderr, "Skipping record (%v)\n", line)
		}
	}

	return txns, nil
}

func loadMemberReferencesFromCsv(path string) (references []referenceLookup, err error) {
	txnFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer txnFile.Close()

	csvReader := csv.NewReader(bufio.NewReader(txnFile))
	firstLine := true
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if len(line) < 2 {
			return nil, fmt.Errorf("Line %v does not have enough fields", line)
		}

		if firstLine {
			firstLine = false
			continue
		}

		membershipIds := strings.Split(line[1], "|")
		references = append(references, *newReferenceLookup(line[0], membershipIds))
	}

	return references, nil
}

func writeTxnsToCsv(path string, txns []bankTxn) error {
	txnFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer txnFile.Close()

	csvWriter := csv.NewWriter(bufio.NewWriter(txnFile))
	for _, txn := range txns {
		err := csvWriter.Write([]string{txn.description, txn.amount.String()})
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}

func writeMembersToCsv(path string, members []member) error {
	memberFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer memberFile.Close()

	csvWriter := csv.NewWriter(bufio.NewWriter(memberFile))
	for _, member := range members {
		err := csvWriter.Write(
			[]string{
				member.memberID,
				member.title,
				member.forenames,
				member.surname,
				member.emailAddress,
			})
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}

func writeMemberIdsToCsv(path string, memberIds []string) error {
	memberIDFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer memberIDFile.Close()

	csvWriter := csv.NewWriter(bufio.NewWriter(memberIDFile))
	for _, memberID := range memberIds {
		err := csvWriter.Write([]string{memberID})
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}
