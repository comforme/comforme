#!/bin/bash
# Postdeploy script for automated Heroku deployment.
cat schema.sql | psql
