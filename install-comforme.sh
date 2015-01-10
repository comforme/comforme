#!/bin/bash
# Install comforme

GO_VERSION="1.4"
PKG_URL="https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz"
INSTALL_DIR="/usr/local"

# Install Go
cd /tmp && wget $PKG_URL && sudo tar -C $INSTALL_DIR -zxf go1.4.linux-amd64.tar.gz
mkdir ~/go
cat >> ~/.bashrc <<-HERE
export GOROOT=/usr/local/go
export GOPATH=\$HOME/go
export PATH=\$PATH:\$GOROOT/bin
export PORT=8080
HERE
source ~/.bashrc

# Install Heroku Toolkit
wget -qO- https://toolbelt.heroku.com/install-ubuntu.sh | sh

# Install Go modules
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin
for repo in lib/pq go-zoo/bone comforme/comforme; do
    go get github.com/$repo
done

PORT=8080 nohup ~/go/bin/comforme &

