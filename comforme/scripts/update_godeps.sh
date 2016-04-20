#!/bin/bash
# Updates Godeps directory

git commit
git rm -rf Godeps
./format.sh
godep save
git add Godeps
git commit -m "Updated Godeps"
