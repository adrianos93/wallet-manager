# wallet-manager

wallet-manager is a small microservice that can be used to manage a users wallets, and perform various operations on them.

Current functionality

- create users
- create wallets for users
- check wallet balance
- deposit money into a wallet
- withdraw money from a wallet
- transfer money between wallets

The microservice was implemented as a REST API and has no other dependencies, so it can be stood up straight from the command line.
The microservice runs on port `8080`, so ensure the port is available.

## Running the service

Requirements:
- Go 1.17 installed

The service binary is included in the source code, so the service can be run by merely invoking.

`./manager` whilst you are in the root directory of the repository.

To rebuild the binary you can run (assuming you are in the root directory of the repo):

`go build -o ./manager /cmd/wallet-manager/main.go`

## Routes and usage

The service has the following available routes:

- GET `/v1/health/wallet-manager` (healthcheck endpoints)
- POST `/v1/user` (creates a user)
- POST `/v1/user/{userId}/wallet` (creates a wallet for the given user)
- GET `/v1/user/{userId}/wallet/{walletId}/balance` (returns the balance on the given wallet for the given user)
- POST `/v1/user/{userId}/wallet/{walletId}/deposit` (processes a deposit on the given wallet for the given user)
- POST `/v1/user/{userId}/wallet/{walletId}/withdraw` (processes a withdrawal on the given wallet for the given user)
- POST `/v1/user/{userId}/wallet/{walletId}/payment` (initiates a payment from the given wallet for the given user)


## Payloads and Responses

The following JSON payloads (these are examples) are required to call the following endpoints:

`POST /v1/user/{userId}/wallet/{walletId}/deposit` and `POST /v1/user/{userId}/wallet/{walletId}/withdrawal`

```json
{"Amount": 100.25}
```

and will respond with:

```json
{"Balance":100}
```

`GET /v1/user/{userId}/wallet/{walletId}/balance` will respond with:

```json
{"Balance":100}
```

`POST /v1/user/{userId}/wallet/{walletId}/payment` will accept:

```json
{
    "Creditor": "wallet1",
    "Amount": 50,
}
```
and respond with:

```json
{
    "TransactionId": "123456",
    "Balance": 50,
}
```

`POST /v1/user` responds with:

```json
{
  "Id": "3fdba7bf30c091836b82b57ab49a0cca"
}
```

`POST /v1/user/{userId}/wallet` responds with:

```json
{
    "Id":"8d3f349c582245d797419754e77d1d82",
    "Balance":0,
}
```

## Design

The design for wallet-manager implements effective Go, where packages should be small and descriptive and functionality is limited to its intended function. 
Ideally, if the implementation were to use a proper backend, such as a SQL database, each package would interact with each other using interfaces, which would create a layer of abstraction and allow for isolated unit testing using mocks.

The app is formed of 3 packages:

- wallet

The wallet package is responsible for performing operations on any given wallet. It also contains a map holding a record of every wallet that exists within wallet-manager. A wallet is given a unique identifier and a balance.

- user

The user package is responsible for creating users and performing user-based actions. A user is formed of a unique identifier and a map containing all of the wallets belonging to them. User based actions entail performing transactions on a wallet the user owns, and it achieves that by invoking the wallet package.

- server

The server packages contains the handlers for each of the routes that are made available by this microservice. It will handle incoming requests and correctly unmarshall them into the appropriate structs to be processed by the user and wallet packages, as well as perform some validation to ensure the requests are valid.
Server was placed in it's own package so it can be more easily tested.

- manager

This package contains global variables, such as the service name, and helper functions.

## Satisfying requirements

The requirements for this app were:

- User can deposit money into her wallet (satisfied by `user` and `wallet` pacakges)
- User can withdraw money from her wallet (satisfied by `user` and `wallet` pacakges)
- User can send money to another user (satisfied by `user` and `wallet` pacakges)
- User can check her wallet balance (satisfied by `user` and `wallet` pacakges)
- The Wallet App should keep track of all the users and wallets it contains (satisfied by the global `Users` and `Wallets` maps)


Non-functional requirements satisfied:

- Input is sanitized using regexes provided by `gorilla/mux`
- Users cannot access other users wallets
- Test coverage is 100%
- Errors are handled and returned with the appropriate response codes
- Descriptive variable and method names were used. Code is self-documenting.
- As a result of using maps, time complexity of this app is O(1), allowing it to maintain the same blazingly fast response time at scale.
- Space complexity is O(n), so it will still scale quite well.


## How to review code.

I would recommend reviewing the code bottom to top, meaning starting from the `wallet` package, then moving on to `user` and finally onto `server`, running unit tests using `go test`.

## Improvements:

- Add GoDocs to the packages
- Add context for tracking requests.
- Add logging.
- Add a client
- Improve server tests by checking response bodies.
- Allow a user to view all their wallets
- Add json validation
- Add SQL backend

## Time spent on solution

- Design: 30 mins
- Writing tests: 2 hours
- Writing code: 2 hours
- writing readme: 30 minutes

Time was divided over 2 days. 