package main

import (
	"testing"
)

func TestMatchMembers(t *testing.T) {
	membersDetails := map[string]*Member{
		"A123456": {"A123456", "Mr", "Joe", "Blogg", "joebloggs@example.com"},
	}
	memberIds := []string{"A123456", "A789012"}
	expectedMatchedMembers := []*Member{
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
		if !expectedMatchedMembers[i].equal(actualMatchedMembers[i]) {
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

func TestIdentifyLeaversAndJoiners(t *testing.T) {
	previousMembers := []*Member{
		{"A123456", "Mr", "Joe", "Blogg", "joebloggs@example.com"},
		{"A789012", "Ms", "Jane", "Doe", "janedoe@example.com"},
		{"", "Mr", "Non", "Member", "nonmember@example.com"},
	}
	currentMembers := []*Member{
		{"A123456", "Mr", "Joe", "Blogg", "joebloggs@example.com"},
		{"", "Mr", "New", "Nonmember", "newnonmember@example.com"},
		{"3456789", "Mr", "New", "Member", "newmember@example.com"},
		{"A456789", "Mr", "Non", "Member", "nonmember@example.com"},
	}
	expectedLeavers := []*Member{
		{"A789012", "Ms", "Jane", "Doe", "janedoe@example.com"},
	}
	expectedJoiners := []*Member{
		{"", "Mr", "New", "Nonmember", "newnonmember@example.com"},
		{"3456789", "Mr", "New", "Member", "newmember@example.com"},
	}

	actualLeavers, actualJoiners := identifyLeaversAndJoiners(previousMembers, currentMembers)

	if actualLeavers == nil || actualJoiners == nil {
		t.Fatalf("One or more returns are nil!")
	}

	if len(actualLeavers) != len(expectedLeavers) {
		t.Fatalf("Leaver counts are not the same (%v, %v)", len(actualLeavers), len(expectedLeavers))
	}

	if len(actualJoiners) != len(expectedJoiners) {
		t.Fatalf("Joiner counts are not the same (%v, %v)", len(actualJoiners), len(expectedJoiners))
	}

	for i := range expectedLeavers {
		if !expectedLeavers[i].equal(actualLeavers[i]) {
			t.Fatalf(
				"%v != %v",
				expectedLeavers[i],
				actualLeavers[i])
		}
	}

	for i := range expectedJoiners {
		if !expectedJoiners[i].equal(actualJoiners[i]) {
			t.Fatalf(
				"%v != %v",
				expectedJoiners[i],
				actualJoiners[i])
		}
	}
}
