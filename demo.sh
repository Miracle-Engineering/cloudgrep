#!/bin/bash

TAG=$(curl https://api.github.com/repos/run-x/cloudgrep/releases/latest | jq -r .tag_name | sed 's/v//')
echo "Tag: ${TAG}"

DOWNLOAD_URL=https://github.com/run-x/cloudgrep/releases/download/v${TAG}/cloudgrep_${TAG}_linux_amd64.tar.gz
echo "Download URL: ${DOWNLOAD_URL}"

wget ${DOWNLOAD_URL}
tar -xzf cloudgrep_${TAG}_linux_amd64.tar.gz

chmod 777 cloudgrep
./cloudgrep demo --bind=0.0.0.0
