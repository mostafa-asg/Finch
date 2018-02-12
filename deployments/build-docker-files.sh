#!/bin/bash

# finch-REST
go build -o finch-rest/finch ../
cp ../configs/finch.yml finch-rest/
cd finch-rest & docker build --tag finch-rest:latest ./finch-rest

# finch load balancer
cd lb & docker build --tag finch-lb:latest ./lb

# cassandra
cd cassandraPrepare & docker build --tag finch-cassandra-prepare:latest ./cassandraPrepare