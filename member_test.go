package main

import (
	"testing"
)

func TestMatchMembers(t *testing.T) {
	membersDetails := map[string]member{
		"A123456": {"A123456", "Mr", "Joe", "Blogg", "joebloggs@example.com"},
	}
	memberIds := []string{"A123456", "A789012"}
	expectedMatchedMembers := []member{
		{"A123456", "Mr", "Joe", "Blogg", "joebloggs@example.com"},
	}
	expectedUnmatchedMembers := []string{"A789012"}

	actualMatchedMembers, actualUnmatchedMembers := matchMembers(membersDetails, memberIds)

	if actualMatchedMembers == nil || actualUnmatchedMembers == nil {
		t.Fatalf("One or more returns are nil!")
	}

	if len(actualMatchedMembers) != len(expectedMatchedMembers) {
		t.Fatalf(
			"Matched member counts are not the same (%v, %v)",
			len(actualMatchedMembers),
			len(expectedMatchedMembers),
		)
	}

	if len(actualUnmatchedMembers) != len(expectedUnmatchedMembers) {
		t.Fatalf(
			"Unmatched member counts are not the same (%v, %v)",
			len(actualUnmatchedMembers),
			len(expectedUnmatchedMembers),
		)
	}

	for i := range expectedMatchedMembers {
		if !expectedMatchedMembers[i].equal(&actualMatchedMembers[i]) {
			t.Fatalf(
				"%v != %v",
				expectedMatchedMembers[i],
				actualMatchedMembers[i])
		}
	}

	for i := range expectedUnmatchedMembers {
		if expectedUnmatchedMembers[i] != actualUnmatchedMembers[i] {
			t.Fatalf(
				"%v != %v",
				expectedUnmatchedMembers[i],
				actualUnmatchedMembers[i])
		}
	}
}
