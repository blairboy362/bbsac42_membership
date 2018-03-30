package main

type Member struct {
	MemberID     string `csv:"MemberId"`
	Title        string `csv:"Title"`
	Forenames    string `csv:"Forenames"`
	Surname      string `csv:"Surname"`
	EmailAddress string `csv:"EmailAddress"`
}

func (m1 *Member) equal(m2 *Member) bool {
	return *m1 == *m2
}

func matchMembers(membersDetails map[string]*Member, memberIds []string) (matchedMembers []*Member, unmatchedMembers []string) {
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
