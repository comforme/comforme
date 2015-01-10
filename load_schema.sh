#!/bin/bash
# Loads default schema and promotes database

# Promote database to ensure consistent name
DATABASE=$(heroku pg:info | grep -oE 'HEROKU_POSTGRESQL_[A-Z]+_URL' | head -n 1)
heroku pg:promote $DATABASE

cat schema.sql | heroku pg:psql --app comforme
