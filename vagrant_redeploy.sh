#!/bin/bash
# Vagrant redeployment script
# This script synchronizes the copy of the Comforme source,
# cleans the cache, then reinstalls and restarts the applications

COMFORME_PATH=~/go/src/github.com/comforme/comforme
COMFORME_BIN=~/go/bin/comforme
PORT=8080

rsync -avc --exclude *~ --exclude *.sw[op] /vagrant/ ${COMFORME_PATH}/
cd $COMFORME_PATH
go clean && go install
pkill comforme
PORT=$PORT $COMFORME_BIN &
