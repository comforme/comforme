#!/bin/bash
# Install comforme

DEBUG_MODE="$1"
GO_VERSION="1.4.2"
PKG_URL="https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz"
INSTALL_DIR="/usr/local"
DB_NAME="comforme"
TABLES="users categories pages posts communities community_memberships sessions"
SEQUENCES="users_id_seq pages_id_seq"
MANDRILL_APIKEY="$2"

# Install Go
cd /tmp && wget $PKG_URL && sudo tar -C $INSTALL_DIR -zxf go1.4.linux-amd64.tar.gz
mkdir ~/go
USERNAME="`whoami`"
cat >> ~/.bashrc <<-HERE
export DEBUG_MODE=${DEBUG_MODE}
export GOROOT=/usr/local/go
export GOPATH=\$HOME/go
export PATH=\$PATH:\$GOROOT/bin
export PORT=8080
export DATABASE_URL="host=/run/postgresql user=${USERNAME} dbname=comforme sslmode=disable"
HERE
if [ ! -z "$MANDRILL_APIKEY" ]; then
echo "export MANDRILL_APIKEY=${MANDRILL_APIKEY}" >> ~/.bashrc
fi
source ~/.bashrc

# Install Heroku Toolkit
wget -qO- https://toolbelt.heroku.com/install-ubuntu.sh | sh

# Install Go modules
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin
go get golang.org/x/crypto/bcrypt
go get golang.org/x/crypto/scrypt
for repo in lib/pq julienschmidt/httprouter keighl/mandrill comforme/comforme; do
    go get github.com/$repo
done

# Install PostgreSQL
sudo apt-get install -y postgresql
sudo -u postgres psql -c "CREATE DATABASE ${DB_NAME};"
sudo -u postgres psql -d $DB_NAME < /vagrant/schema.sql
sudo -u postgres psql -d $DB_NAME -c "CREATE USER ${USERNAME}"
for seq in $SEQUENCES; do
    sudo -u postgres psql -d $DB_NAME -c "GRANT USAGE, SELECT ON SEQUENCE ${seq} to ${USERNAME};"
done

for table in $TABLES; do
    sudo -u postgres psql -d $DB_NAME -c "GRANT ALL PRIVILEGES ON TABLE ${table} TO ${USERNAME};"
done

# Generate secret key used for hashing salt
SECRET=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

# Run comforme executable with environmental variables
DEBUG_MODE=$DEBUG_MODE SECRET=$SECRET MANDRILL_APIKEY=$MANDRILL_APIKEY PORT=8080 nohup ~/go/bin/comforme &

