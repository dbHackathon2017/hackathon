# Run Instructions

There is no binary to run, no main() function, only test functions to execute some sample code for demonstration purposes.
A remote factomd instance has been deployed that blocks produces blocks every second, so submitted entries are added to the chain nearly immediately.
The following entry credit key has lots of money in it to use for demonstration purposes:

- Es3cpDrGJRZpJBqZ3PwdohDpmMcXqmr8PuN2yyzBdB2rZ2McEtu1
- EC29nUzTTopMuwEHgPGZ8eBvTGEgzPHErbJU8HVPXxTvKjP37hK6

factom-write/util.go : GetECAddress() will return this address.

To make use of the remote factomd instance, use the following function:

```
factom.SetFactomdServer(constants.REMOTE_HOST)
```

See factom-write/pension_test.go for an example of the usage of the function.

# Requirements:

Golang go (tested with go version 1.7.4):
```
apt-get install golang-go
```
git (tested with go version 2.1.4):
```
sudo apt-get install git
```
Glide (tested with version 0.13.0-dev)
```
go get -u github.com/Masterminds/glide
cp ~/go/bin/glide /usr/local/bin/
```
# Installing Factom

Persistantly add $GOPATH to your system's environment variables:
```
mkdir $HOME/go && printf "export PATH=$PATH:/usr/local/go/bin\nexport GOPATH=$HOME/go\nexport PATH=$PATH:$GOPATH/bin" >> ~/.profile
```
Reload global system environment:
```
source ~/.profile
```
Download relevant factom source code
```
git clone https://github.com/FactomProject/factomd $GOPATH/src/github.com/FactomProject/factomd
git clone https://github.com/FactomProject/factom-cli $GOPATH/src/github.com/FactomProject/factom-cli
git clone https://github.com/FactomProject/factom-walletd $GOPATH/src/github.com/FactomProject/factom-walletd
git clone https://github.com/FactomProject/enterprise-wallet $GOPATH/src/github.com/FactomProject/enterprise-wallet
```
Get the dependencies and build each factom program one-by-one
```
glide cc
cd $GOPATH/src/github.com/FactomProject/factomd
glide install
go install -v -ldflags "-X github.com/FactomProject/factomd/engine.Build=`git rev-parse HEAD`"
cd $GOPATH/src/github.com/FactomProject/factom-cli
glide install
go install -v
cd $GOPATH/src/github.com/FactomProject/factom-walletd
glide install
go install -v
cd $GOPATH/src/github.com/FactomProject/enterprise-wallet
glide install
go install -v
cd $GOPATH/src/github.com/FactomProject/factomd
```
Create and make use of custom factom configurations
```
mkdir -p ~/.factom/m2/
cp $GOPATH/src/github.com/FactomProject/factomd/factomd.conf ~/.factom/m2/
```

# Factom Types

## Pension

| PensionID|ChainID|
|---|---|
|ExtID (0)|"Pension Chain"|
|ExtID (1)|Pension Company|
|ExtID (2)|Authority PubKey|
|ExtID (3)|Hash|
|ExtID (4)| nonce |
|Content|Document data

## Transaction

### Value change

| Transaction|Factom Entry|
|---|---|
|ExtID (0)|"Transaction Value Change"|
|ExtID (1)|UserType|
|ExtID (2)|Value Change|
|ExtID (3)|PensionID|
|ExtID (4)|ToPensionID|
|ExtID (5)|PersonSubmit|
|ExtID (6)|Timestamp|
|ExtID (7)|PubKey|
|ExtID (8)|Signature ExtID(0-6)|
|Content|Document data
