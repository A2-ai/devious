#!/bin/sh

INSTALL_DIR=.local/bin

# If no URL is provided, use the default
INSTALL_URL=$1
if [ -z "$INSTALL_URL" ]; then
    echo "No URL provided, using default"
    INSTALL_URL=https://github.com/A2-ai/devious/releases/latest/download/dvs_Linux_x86_64.tar.gz
fi

# Create the user bin directory and add it to the PATH if it doesn't exist
mkdir ~/$INSTALL_DIR && (echo 'export PATH="$HOME/'$INSTALL_DIR':$PATH"' | tee -a ~/.bashrc ~/.profile) && source ~/.bashrc

# Download and extract the binary
cd ~/$INSTALL_DIR
wget $INSTALL_URL -O dvs.tar.gz
tar -xzf dvs.tar.gz dvs
chmod +x dvs
rm dvs.tar.gz