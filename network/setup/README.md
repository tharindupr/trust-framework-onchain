# Fabric 2.0
# May 2020

# how to use vagrant (If you are on Windows)
Make sure to install Virtual Box and Vagrant
cd into this directory.
vagrant up


# 1 Open a terminal and give execute permissions
cd network/setup  <br/>
chmod -R 755 ./*  <br/>

# 2 Install Pre-Requisites & Validate
./install-prereqs.sh  <br/>
source ~/.profile  <br/>
source ~/.bashrc   <br/>
# Logout and Login Again (To reflect the changes to GOPATH)
sudo ./validate-prereqs.sh


# 3 Install the Fabric binaries & images
sudo -E ./install-fabric.sh  <br/>
sudo ./validate-fabric.sh  <br/>

<!-- # 4 Install Hyperledger Explorer tool
./install-explorer.sh
sudo ./validate-explorer.sh -->

# 5 Install the Go Tools
./install-gotools.sh





