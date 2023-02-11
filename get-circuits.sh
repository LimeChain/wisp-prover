#!/bin/bash

# Download the tar file from the Amazon S3 bucket
echo "Starting download of compiled-circuits.tar from Amazon S3 bucket crc-circuits..."
curl https://s3.amazonaws.com/crc-circuits/circuits.tar -o circuits.tar

# Unzip the tar file into the ./circuits directory
echo "Unzipping compiled-circuits.tar into ./circuits directory..."
mkdir -p ./circuits
tar -xvf circuits.tar -C ./circuits

# Remove the downloaded tar file
echo "Deleting compiled-circuits.tar..."
rm circuits.tar