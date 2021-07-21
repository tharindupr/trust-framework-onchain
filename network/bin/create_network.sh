#!/bin/bash
export FABRIC_LOGGING_SPEC=INFO


# Removing previous assets
echo    '================ Removing previous Assets================'
rm -r ./../config/*.tx
rm -r ./../config/*.block
rm -r ./../config/crypto-config/*


echo    '================Generating the Crypto Assets================'
cryptogen generate --config=../config/crypto-config.yaml --output=../config/crypto-config



# Create the Genesis Block
echo    '=============== Generating the Genesis Block ================'
GENESIS_BLOCK=../config/genesis.block
SYS_CHANNEL="sys-channel"
CONFIGTX_PATH=../config/

configtxgen -profile OrdererGenesis -configPath $CONFIGTX_PATH -channelID $SYS_CHANNEL -outputBlock $GENESIS_BLOCK 




echo    '================ Generate channel configuration block ================'
CHANNEL_ID="mychannel"
CHANNEL_CREATE_TX=../config/mychannel.tx
configtxgen -profile BasicChannel -configPath $CONFIGTX_PATH -outputCreateChannelTx $CHANNEL_CREATE_TX -channelID $CHANNEL_ID


echo    '================ Generate the anchor Peer updates ======'

ANCHOR_UPDATE_TX=../config/Org1MSPanchors.tx
configtxgen -profile BasicChannel -configPath ../config/ -outputAnchorPeersUpdate $ANCHOR_UPDATE_TX -channelID $CHANNEL_ID -asOrg Org1MSP

ANCHOR_UPDATE_TX=../config/Org2MSPanchors.tx
configtxgen -profile BasicChannel -configPath ../config/ -outputAnchorPeersUpdate $ANCHOR_UPDATE_TX -channelID $CHANNEL_ID -asOrg Org2MSP

# ANCHOR_UPDATE_TX=$DIR/../config/Digi-03MSPanchors.tx
# configtxgen -profile DigiBlocksChannel -outputAnchorPeersUpdate $ANCHOR_UPDATE_TX -channelID $CHANNEL_ID -asOrg Digi-03MSP

# ANCHOR_UPDATE_TX=$DIR/../config/Digi-04MSPanchors.tx
# configtxgen -profile DigiBlocksChannel -outputAnchorPeersUpdate $ANCHOR_UPDATE_TX -channelID $CHANNEL_ID -asOrg Digi-04MSP

# ANCHOR_UPDATE_TX=$DIR/../config/Digi-05MSPanchors.tx
# configtxgen -profile DigiBlocksChannel -outputAnchorPeersUpdate $ANCHOR_UPDATE_TX -channelID $CHANNEL_ID -asOrg Digi-05MSP


# export FABRIC_LOGGING_SPEC=INFO
# export FABRIC_CFG_PATH=$DIR/../config
# export COMPOSE_PROJECT_NAME=digiblocks
# export IMAGE_TAG=latest
# source   $DIR/.env

# docker-compose -f $DIR/../devenv/composer/docker-compose.base.yaml up

# echo '###################### Stoping previous containers ###############'
# docker stop $(docker ps -a -q)
# docker rm $(docker ps -a -q)
