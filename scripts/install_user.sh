#!/bin/sh

INSTALL_DIR=.local/bin

mkdir ~/$INSTALL_DIR
cd ~/$INSTALL_DIR
echo 'export PATH="$HOME/'$INSTALL_DIR':$PATH"' | tee -a ~/.bashrc ~/.profile
wget $1 -O dvs.tar.gz
tar -xzf dvs.tar.gz dvs
chmod +x dvs
rm dvs.tar.gz