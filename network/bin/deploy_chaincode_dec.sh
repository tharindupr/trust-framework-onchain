#!/bin/bash

./package_chaincode.sh "./../../gocc/src/github.com/tharindupr/dec" "deccontract" $1
./install_chaincode.sh "./../../gocc/src/github.com/tharindupr/dec" "deccontract" $1
./approve_chaincode.sh "./../../gocc/src/github.com/tharindupr/dec" "deccontract" $1
./commit_chaincode.sh "./../../gocc/src/github.com/tharindupr/dec" "deccontract" $1


echo "===================== Invoking Init ====================="
./invoke_init.sh deccontract $1
