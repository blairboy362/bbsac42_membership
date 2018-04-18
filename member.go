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

func identifyLeaversAndJoiners(previousMembers, currentMembers []*Member) (leavers, joiners []*Member) {
	leavers = []*Member{}
	joiners = []*Member{}
	for _, previousMember := range previousMembers {
		found := false
		for _, currentMember := range currentMembers {
			if len(previousMember.MemberID) > 0 && len(currentMember.MemberID) > 0 {
				if previousMember.MemberID == currentMember.MemberID {
					found = true
					break
				}
			} else if len(previousMember.EmailAddress) > 0 && len(currentMember.EmailAddress) > 0 {
				if previousMember.EmailAddress == currentMember.EmailAddress {
					found = true
					break
				}
			}
		}

		if !found {
			leavers = append(leavers, previousMember)
		}
	}

	for _, currentMember := range currentMembers {
		found := false
		for _, previousMember := range previousMembers {
			if len(previousMember.MemberID) > 0 && len(currentMember.MemberID) > 0 {
				if previousMember.MemberID == currentMember.MemberID {
					found = true
					break
				}
			} else if len(previousMember.EmailAddress) > 0 && len(currentMember.EmailAddress) > 0 {
				if previousMember.EmailAddress == currentMember.EmailAddress {
					found = true
					break
				}
			}
		}

		if !found {
			joiners = append(joiners, currentMember)
		}
	}

	return leavers, joiners
}
