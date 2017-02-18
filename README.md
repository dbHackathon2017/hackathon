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

See factom-write pension_test.go for an example of the usage of the function.

# Requirements:

Golang go (tested on version 1.6)

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
