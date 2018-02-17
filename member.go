package main

type member struct {
	memberID, title, forenames, surname, emailAddress string
}

func (m1 *member) equal(m2 *member) bool {
	return m1.memberID == m2.memberID &&
		m1.title == m2.title &&
		m1.forenames == m2.forenames &&
		m1.surname == m2.surname &&
		m1.emailAddress == m2.emailAddress
}

func matchMembers(membersDetails map[string]member, memberIds []string) (matchedMembers []member, unmatchedMembers []string) {
	for _, memberID := range memberIds {
		member, ok := membersDetails[memberID]
		if ok {
			matchedMembers = append(matchedMembers, member)
		} else {
			unmatchedMembers = append(unmatchedMembers, memberID)
		}
	}

	return matchedMembers, unmatchedMembers
}
