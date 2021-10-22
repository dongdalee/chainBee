#!/bin/bash
set -ev

docker-compose -f docker-compose-ca.yaml up -d ca.sales1.chainbee.com

sleep 1
cd $GOPATH/src/chainBee/application/sdk
node enrollAdmin.js
sleep 1
node registUsers.js