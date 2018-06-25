bbsac42_membership: main.go bank_txn.go reference_lookup.go file_handling.go
	go build -o bbsac42_membership

test:
	go test -v

clean:
	rm -f bbsac42_membership *csv

run-local: bbsac42_membership
	./bbsac42_membership "/home/daniel/Dropbox/BBSAC42 Committee/Branch Management/Treasurership - Membership/bbsac42_membership"

.PHONY: clean run-local test
