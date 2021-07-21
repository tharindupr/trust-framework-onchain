#!/bin/bash

# DO NOT Execute this script with sudo
if [ $SUDO_USER ]; then
    echo "Do not execute with sudo !!!    ./install-prereqs.sh"
    echo "Aborting!!!"
    exit 0
fi

# Install JQ
sudo apt-get install -y jq

sudo ./tools/docker.sh    
sudo ./tools/compose.sh   
sudo -E ./tools/go.sh     
sudo ./tools/node.sh


