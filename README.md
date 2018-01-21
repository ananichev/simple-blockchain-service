# Simple blockchain service

### Quick start
```
# download service
go get github.com/ananichev/simple-blockchain-service
cd $GOPATH/src/github.com/ananichev/simple-blockchain-service

# install dependencies
go get ./...

# build service
go build -o server .

# run server on 3000 port
./server
```
Populate store:
```
./populate_db
```

### API

`POST /add_data` - adds row to internal buffer. When rows count reaches 5 in buffer - service will create block.

`GET /last_blocks/{N}` - returns N last blocks
