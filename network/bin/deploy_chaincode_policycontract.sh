
#!/bin/bash
./package_chaincode.sh "./../../gocc/src/github.com/tharindupr/policy_contract" "policycontract" $1
./install_chaincode.sh "./../../gocc/src/github.com/tharindupr/policy_contract" "policycontract" $1
./approve_chaincode.sh "./../../gocc/src/github.com/tharindupr/policy_contract" "policycontract" $1
./commit_chaincode.sh "./../../gocc/src/github.com/tharindupr/policy_contract" "policycontract" $1



echo "===================== Invoking Init ====================="

./invoke_init.sh policycontract $1
