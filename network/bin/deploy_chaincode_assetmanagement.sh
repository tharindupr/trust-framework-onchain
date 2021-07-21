
#!/bin/bash
./package_chaincode.sh "./../../gocc/src/github.com/tharindupr/asset_management" "assetcontract" $1
./install_chaincode.sh "./../../gocc/src/github.com/tharindupr/asset_management" "assetcontract" $1
./approve_chaincode.sh "./../../gocc/src/github.com/tharindupr/asset_management" "assetcontract" $1
./commit_chaincode.sh "./../../gocc/src/github.com/tharindupr/asset_management" "assetcontract" $1


echo "===================== Invoking Init ====================="
./invoke_init.sh assetcontract $1