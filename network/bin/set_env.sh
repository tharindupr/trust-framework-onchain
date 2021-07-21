export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/../config/crypto-config/ordererOrganizations/digiblocks.com/orderers/orderer.digiblocks.com/msp/tlscacerts/tlsca.digiblocks.com-cert.pem
export PEER0_ORG1_CA=${PWD}/../config/crypto-config/peerOrganizations/org1.digiblocks.com/peers/peer0.org1.digiblocks.com/tls/ca.crt
export PEER0_ORG2_CA=${PWD}/../config/crypto-config/peerOrganizations/org2.digiblocks.com/peers/peer0.org2.digiblocks.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/../config/


export CHANNEL_NAME=mychannel

setGlobalsForOrderer(){
    export CORE_PEER_LOCALMSPID="OrdererMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/../config/crypto-config/ordererOrganizations/digiblocks.com/orderers/orderer.digiblocks.com/msp/tlscacerts/tlsca.digiblocks.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../config/crypto-config/ordererOrganizations/digiblocks.com/users/Admin@digiblocks.com/msp
    
}

setGlobalsForPeer0Org1(){
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../config/crypto-config/peerOrganizations/org1.digiblocks.com/users/Admin@org1.digiblocks.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
}

setGlobalsForPeer1Org1(){
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../config/crypto-config/peerOrganizations/org1.digiblocks.com/users/Admin@org1.digiblocks.com/msp
    export CORE_PEER_ADDRESS=localhost:8051
    
}

setGlobalsForPeer0Org2(){
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../config/crypto-config/peerOrganizations/org2.digiblocks.com/users/Admin@org2.digiblocks.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
    
}

setGlobalsForPeer1Org2(){
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/../config/crypto-config/peerOrganizations/org2.digiblocks.com/users/Admin@org2.digiblocks.com/msp
    export CORE_PEER_ADDRESS=localhost:10051
    
}


# CHANNEL_NAME="mychannel"
# CC_RUNTIME_LANGUAGE="golang"
# VERSION="1"
# CC_SRC_PATH="./../../gocc/src/github.com/tharindupr/fabcar"
# CC_NAME="fabcar"