#!/bin/sh

curl -XPUT -v -d '{"dateOfBirth": "2021-10-01"}' 'http://localhost:8080/hello/apple'
curl -XPUT -v -d '{"dateOfBirth": "2021-11-01"}' 'http://localhost:8080/hello/pear'

curl -XGET -v 'http://localhost:8080/hello/apple'
