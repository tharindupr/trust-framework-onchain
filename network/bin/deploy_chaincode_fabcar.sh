
./package_chaincode.sh "./../../gocc/src/github.com/tharindupr/fabcar" "fabcar" $1
./install_chaincode.sh "./../../gocc/src/github.com/tharindupr/fabcar" "fabcar" $1
./approve_chaincode.sh "./../../gocc/src/github.com/tharindupr/fabcar" "fabcar" $1
./commit_chaincode.sh "./../../gocc/src/github.com/tharindupr/fabcar" "fabcar" $1


echo "===================== Invoking Init ====================="
./invoke_init.sh fabcar $1