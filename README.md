## Build

`$ go build -v main.go`

OR

`$ make`

OR

The Dockerfile will copy the executable (bin/data-linux) to the container so before deploying be
sure to `make xbuild`.

`$ docker build -t rithium/stor-data .`

## Cross Compile

`$ env GOOS=linux GOARCH=arm64 go -v build main.go`

OR

`$ make xbuild`
## Run

`$ env S_CASSANDRA_URL=... ./main`

OR

`$ docker run -it -e ZK_HOSTS=45.79.159.146:2181 -e S_CASSANDRA_KEYSPACE=stor -e S_URL=0.0.0.0 -p 80:80 rithium/stor-data`

Will allow you to make requests from the host which is useful for testing the REST endpoints with postman.

## Test

Integration tests will mock the database connection in order to run through
the HTTP handlers.

`$ go test -race -tags=unit *.go`

`$ go test -race -tags=integration *.go`


## Env Vars

Available configuration variables and their default values:

`S_URL=localhost`

`S_PORT=80`


`S_CASSANDRA_URL=127.0.0.1`

`S_CASSANDRA_USER=`

`S_CASSANDRA_PASS=`

`S_CASSANDRA_KEYSPACE=stor`

`S_CASSANDRA_PROTO=4`

## Docker
### Dockerfile

The Dockerfile is built on the rithium/smartstackgo image which contains Nerve and Synapse for service discovery and
the default golang:1.7 image.

### Race Flag
The race detection flag will not work on Alpine Linux due to it not using glibc.

### Deploy

`$ docker pull rithium/data:latest`

Replace `ZK_HOSTS` with the local discovery container/instance url:port.

`$ docker run -it -e ZK_HOSTS=172.17.0.2:2181 rithium/data`

