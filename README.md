# panacea-dep-consumer

A HTTP server for DEP(Data Exchange Protocol) consumers which

## Features

- Request the storing a data to consumer after the oracle verify the data.

## Installation

```bash
make build
make test
make install

consumerd -listen-addr="" -grpc-addr=""
# ex) consumerd -listen-addr="127.0.0.1:8080" -grpc-addr="http://127.0.0.1:9090"
# The `-grpc-addr`values should be with URL scheme such as `http`, `https`.
```

## Request Store a Data
```bash
curl -v -X POST "http://${YOUR_HTTP_SERVER}/v1/deals/${dealID}/data/${dataHash}" -H "Authorization: Bearer ${ORACLE_JWT}"

## The ORACLE_JWT is an JWT which is signed by oracle private key.
```
If the storing a data success, the response will show following message:
```bash
success to store data
```

And the data will be stored in $HOME/${dealID} directory with the file name ${dataHash}.