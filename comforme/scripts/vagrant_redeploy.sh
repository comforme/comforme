#!/bin/bash
# Vagrant redeployment script
# This script synchronizes the copy of the Comforme source,
# cleans the cache, then reinstalls and restarts the applications

COMFORME_PATH=~/go/src/github.com/comforme/comforme
COMFORME_BIN=~/go/bin/comforme
export PORT=8080

rsync -avc --exclude *~ --exclude *.sw[op] --exclude .git /vagrant/ ${COMFORME_PATH}/
cd $COMFORME_PATH
pkill comforme
go clean && go install
$COMFORME_BIN &
