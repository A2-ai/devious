#!/bin/sh

# Wipe the testing environment
rm -rf testing/environment

# Make directories
mkdir testing/environment
mkdir testing/environment/local
mkdir testing/environment/remote

# Initialize mock git repository
mkdir testing/environment/local/.git

# Copy test files to local
cp -r testing/data/* testing/environment/local