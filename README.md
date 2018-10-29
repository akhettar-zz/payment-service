# Payment Service

## Setting up GO environment and package management tool

### Install docker
You can download docker from [here](https://docs.docker.com/docker-for-mac/install/#what-to-know-before-you-install)

### Install dep
The payment service uses `dep` to install its dependencies. Details on how to install `dep` can be found [here](https://github.com/golang/dep)

- Adding the project for the first time:    
  ```bash 
  dep ensure -vendor-only 
  ```

- Adding a new dependency 
  ```bash 
  dep ensure -add <dependecy path> 
  ```
### Install GO

1. Install GO: `brew install go`
2. Install DEP `brew install dep`
3. Set the GOPATH env variable and add it to the path in your bash profile (~/.profile): 
```bash
export GOPATH=$WAVE_PROJECT_HOME/go
export PATH=$PATH:$GOPATH/bin
```

## Running test
Note when you run the test for the first time it will take some time. The test run a `cockroach database` docker container in the background 

`scripts/./run-tests.sh`

## Swagger

We are using gin-swagger (https://github.com/swaggo/gin-swagger) - see the comments added in each endpoint handler (`api/handler.go`)

You can access swagger documentation: http://localhost:8080/swagger/index.html. The swagger doc has already been generated

To generate a newer version of the operations to be exposed. Run these steps:
1. Install the swag tool using this script: `./scripts/install_swag.sh`
2. Run `swag init`. This will generate a new swagger docs file see - `./docs/docs.go`


## Running the server

`docker-compose up --build`

## Interacting with the server

### Health endpoint
`curl http://localhost:8080/health`

```json
HTTP/1.1 200 OK
Content-Length: 29
Content-Type: application/json; charset=utf-8
Date: Mon, 20 Aug 2018 07:18:49 GMT

{
    "message": "Up and running!"
}
```

### Create Payment

`curl -d @samples/paymentRequest.json -H "Content-Type: application/json" -X POST http://localhost:8080/payment`

`{
     "id": "5bd7506a9900b30008edf576",
     "organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"
 }`

### Query All Payments

`curl -X GET http://localhost:8000/payment`


### Query Given A Payment

`curl -X GET http://localhost:8080/payment/5bd7506a9900b30008edf576`


### Delete Payment

`curl -X DELETE http://localhost:8080/payment/5bd7506a9900b30008edf576`


### Update Payment

`curl -d @samples/paymentRequest.json -H "Content-Type: application/json" -X PUT http://localhost:8080/payment/5bd7506a9900b30008edf576`

## Mock
To generate a mock for an interface run the followings:
1- Install `gomock` `go get github.com/golang/mock/gomock`
2- Install `mockgen` `go get github.com/golang/mock/mockgen`

`mockgen` binary should be installed in $GOPATH/bin

To gnerate the mock repository run the following command:
1 - Create `mock` directory.
2 - `$GOPATH/bin/mockgen -destination=mocks/mock_repository.go -package=mocks github.com/payment-service/repository Repository`
