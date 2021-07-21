#!/bin/bash

./package_chaincode.sh "./../../gocc/src/github.com/tharindupr/identity" "identitycontract" $1
./install_chaincode.sh "./../../gocc/src/github.com/tharindupr/identity" "identitycontract" $1
./approve_chaincode.sh "./../../gocc/src/github.com/tharindupr/identity" "identitycontract" $1
./commit_chaincode.sh "./../../gocc/src/github.com/tharindupr/identity" "identitycontract" $1


echo "===================== Invoking Init ====================="
./invoke_init.sh identitycontract $1
