bbsac42_membership: main.go bank_txn.go reference_lookup.go file_handling.go
	go build -o bbsac42_membership

test:
	go test -v

clean:
	rm -f bbsac42_membership *csv

run-local: bbsac42_membership
	./bbsac42_membership /home/daniel/BSAC/Treasurer/BSACFeb14th2018.csv /home/daniel/BSAC/Treasurer/reference_member_mappings.csv /home/daniel/BSAC/Treasurer/membership_details.csv /home/daniel/BSAC/Treasurer/new_members.csv

.PHONY: clean run-local test
