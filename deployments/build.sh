#!/bin/bash

go build -o finch-rest/finch ../
cp ../configs/finch.yml finch-rest/
cd finch-rest & docker build --tag finch-rest:latest ./finch-rest