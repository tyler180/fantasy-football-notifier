#!/bin/bash

# This script is used to build the Go project and zip the binary.

# Set the directory and binary name
DIRECTORY="ffnotifier"
BINARY_NAME="bootstrap"
ZIP_FILE_NAME="ffnotifier.zip"

# Navigate to the directory
cd $DIRECTORY || { echo "Directory $DIRECTORY not found"; exit 1; }

if [ -e $BINARY_NAME.old ]; then
    rm $BINARY_NAME.old
fi

mv $BINARY_NAME $BINARY_NAME.old 2>/dev/null

# Build the Go binary
GOOS=linux GOARCH=arm64 go build -o $BINARY_NAME main.go

# Check if the build was successful
if [ $? -ne 0 ]; then
  echo "Go build failed"
  exit 1
fi
OLD_BINARY_HASH=$(md5sum $BINARY_NAME.old | awk '{print $1}')
NEW_BINARY_HASH=$(md5sum $BINARY_NAME | awk '{print $1}')
# Zip the binary
if [ $OLD_BINARY_HASH != $NEW_BINARY_HASH ]; then
    echo "Binary changed, zipping"
    if [ -e ${ZIP_FILE_NAME} ]; then
        rm ${ZIP_FILE_NAME}
    fi
    zip ${ZIP_FILE_NAME} $BINARY_NAME
else
    echo "Binary not changed, skipping zipping"
    exit 0
fi

# Check if the zip was successful
if [ $? -ne 0 ]; then
  echo "Failed to zip the binary"
  exit 1
fi

echo "Build and zip successful"
