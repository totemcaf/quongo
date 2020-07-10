#!/usr/bin/env bash

# https://ant0ine.github.io/go-json-rest/
#go get github.com/ant0ine/go-json-rest/tree/v3.3.2/rest
dep ensure -add github.com/ant0ine/go-json-rest/tree/v3.3.2/rest

# https://github.com/StephanDollberg/go-json-rest-middleware-jwt

# http://labix.org/mgo
#go get gopkg.in/mgo.v2
dep ensure -add gopkg.in/mgo.v2

# MongoDB 3.4.6
# https://en.wikipedia.org/wiki/ApacheBench
#dep ensure -add github.com/mongodb/mongo-go-driver/mongo
