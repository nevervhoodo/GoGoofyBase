# GoGoofyBase
## No info about task2 now!
Simple database realization
Used Golang-1.9

### Dependenses

github.com/julienschmidt/httprouter

### Usage

Placed on port 8080

```
GET /v1/records => returns whole table 
GET / => just home page
GET /v1/records/:id => returns record with selected id
GET /v1/records/key => return record with selected key
POST /v1/records "{"Key": "new_key", "Value": "new_value"}" => adds record to table
	*possible curl example*
	*curl -XPOST -H "Content-Type: application/json" -d '{"Key": "new", "Value": "hhh"}' localhost:8080/v1/records*
```

### Database table structure
Stores pairs key:value

Field Name 	| Type 		| Meaning
----------------|---------------|-----------------------------
ID 		| int 		| id for enumeration of records
Key 		| string 	| key field for database
Value 		| string 	| value for selected key



### Authors

...

### TODO list
- [x] create README
- [ ] add record remove
- [ ] add table clear
- [x] add update record
- [ ] make distributed database
- [ ] secure input validation and filtration
- [ ] add OOP-like structure 
- [x] divide into several package files

