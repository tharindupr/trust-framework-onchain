#!/bin/bash

./package_chaincode.sh "./../../gocc/src/github.com/tharindupr/kpis" "kpicontract" $1
./install_chaincode.sh "./../../gocc/src/github.com/tharindupr/kpis" "kpicontract" $1
./approve_chaincode.sh "./../../gocc/src/github.com/tharindupr/kpis" "kpicontract" $1
./commit_chaincode.sh "./../../gocc/src/github.com/tharindupr/kpis" "kpicontract" $1


echo "===================== Invoking Init ====================="
./invoke_init.sh kpicontract $1
