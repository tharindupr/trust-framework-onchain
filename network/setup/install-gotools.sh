#!/bin/bash


echo "Installing govendor....please wait"

# Changes the object owner for all folder under $HOME
sudo find ~ -user root -exec sudo chown $USER: {} +

# Installs the tools for Go
# https://github.com/kardianos/govendor/wiki/Govendor-CheatSheet
go get -u github.com/kardianos/govendor

