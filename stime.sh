#!/bin/sh

# Installer:
#
# installConfig/
#	config.txt
#	hive.tf
# installSTime.sh --> Will just work over installConfig folder 
#	1. Read properties of config.txt and create variables
#	2. Parse/Prepare hive.tf
#	3. Use terraform to deploy hive to VPS



# Check if an installation already exists

case "$1" in

    "remove" )
cd ./installConfig/
rm stclient
rm -rf electronGUI
rm -rf go*
rm -rf vpskeys/*
rm -rf reports/*
rm -rf ../src/github.com/
rm -rf ../src/golang.org/
rm -rf ../src/google.golang.org/
rm -rf ../src/go.opencensus.io/
rm -rf ../src/cloud.google.com/
rm -rf ../pkg 
        ;;
    "install" )


USERNAME=$2
PASSWORD=$3
HASH=$(htpasswd -bnBC 14 "" ${PASSWORD} | tr -d ':\n')
HIVIP=$4
HIVPORT=$5
HIVTLSHASH=$6

# Download GO and Compile Hive
wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz -P ./installConfig/
tar xvf ./installConfig/go1.10.3.linux-amd64.tar.gz -C ./installConfig/
export GOPATH="$(pwd)"
./installConfig/go/bin/go get "github.com/mattn/go-sqlite3"
./installConfig/go/bin/go get "github.com/gorilla/mux"
./installConfig/go/bin/go get "golang.org/x/crypto/blowfish"
./installConfig/go/bin/go get "golang.org/x/crypto/bcrypt"
./installConfig/go/bin/go get "golang.org/x/net/context"
./installConfig/go/bin/go get "golang.org/x/oauth2"
./installConfig/go/bin/go get "golang.org/x/oauth2/google"
./installConfig/go/bin/go get "google.golang.org/api/gmail/v1"


#Compile client with target variables and prepare electron front-end
cd ./installConfig
GOOS=linux GOARCH=amd64 ./go/bin/go build --ldflags "-X main.username=${USERNAME} -X main.password=${PASSWORD} -X main.roasterString=${HIVIP}:${HIVPORT} -X main.fingerPrint=${HIVTLSHASH}" -o stclient client
cp -r ../src/client/electronGUI/ .
cd electronGUI/
sudo apt-get install -y nodejs
npm install

exit 1
        ;;

*) 	cd ./installConfig/
	./stclient &
	cd electronGUI
	npm start
	pkill stclient
	exit 1
   ;;
esac
