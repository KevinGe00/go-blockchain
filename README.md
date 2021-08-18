# Go-Blockchain

This is a project I build to demonstrate the foundational ideas of a blockchain works. It imagines blocks on a blockchain using a single server and does not use a peer-to-peer network as this project is not meant to be a demonstration of decentralization. There are also some basic transaction/cyrptography examples related to blockchain/cryptocurrency shown through the code.

## Installation/Run
```
git clone https://github.com/KevinGe00/go-blockchain.git
```
In root directory, to run the server:
```
go build
./go-blockchain
```
"Listening on port 10000" should appear and indicates that the server is now up and running.

## Blockchain as an API
Please find the [Postman collection](go-blockchain.postman_collection.json), [import it into Postman](https://learning.postman.com/docs/getting-started/importing-and-exporting-data/) to see examples how to call the API in order to view the blockchain and mine blocks to be added to the blockchain.

## Testing
In root directory, run this to run all the *_test.go files:
```
go test -v
```
You can also look through the code coverage of the tests in detail by running:
```
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```
