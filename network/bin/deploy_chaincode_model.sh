#!/bin/bash

./package_chaincode.sh "./../../gocc/src/github.com/tharindupr/models" "modelcontract" $1
./install_chaincode.sh "./../../gocc/src/github.com/tharindupr/models" "modelcontract" $1
./approve_chaincode.sh "./../../gocc/src/github.com/tharindupr/models" "modelcontract" $1
./commit_chaincode.sh "./../../gocc/src/github.com/tharindupr/models" "modelcontract" $1


echo "===================== Invoking Init ====================="
./invoke_init.sh modelcontract $1
