# bbsac42_membership
Manage a list of current active members using the bank account as the source of truth.

## Usage
```
./bbsac_membership <path_to_txns_file> <path_to_ref_mapping_file> <path_to_members_details_file> <path_to_new_members_file>
```

It spits out a CSV file with details of all current members. In addition, it writes several other helpful files with details of transactions that didn't match, or missing member ID info.

## txns_file
A CSV file from the bank.
```
Date,Type,Description,Paid out,Paid in,Balance
16-Jan-18,CR,Some Ref , ,18.5,12345.67
...
```

## ref_mapping_file
A simple CSV file that maps bank references to BSAC member IDs.
```
Reference,MemberIds
Some Ref,A123456
...
```

## members_details_file
A simple CSV file that contains some member details (based on the monthly email received from BSAC HQ).
```
MemberId,Title,Forenames,Surname,EmailAddress
A123456,Mr,Joe,Bloggs,joe.bloggs@example.com
...
```

## new_members_file
A simple CSV file that contains details of members that have not yet got their BSAC membership number, or that otherwise should be included in bulletins but are no longer paying.
```
Title,Forenames,Surname,EmailAddress
Miss,Jane,Doe,jane.doe@example.com
...
```
