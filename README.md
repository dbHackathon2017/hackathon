# Run Instructions

There is no binary to run. There is no main, only test functions to execute some code.
I setup a remote factomd on 1 second blocks, so our entries go in immediatly. This entry credit
key has lots of money in it to use.

- Es3cpDrGJRZpJBqZ3PwdohDpmMcXqmr8PuN2yyzBdB2rZ2McEtu1
- EC29nUzTTopMuwEHgPGZ8eBvTGEgzPHErbJU8HVPXxTvKjP37hK6

factom-write/util.go : GetECAddress Returns this address

To use remote factomd

```
factom.SetFactomdServer(constants.REMOTE_HOST)
```

See factom-write pension_test.go for example

# Factom Types

## Pension

| PensionID|ChainID|
|---|---|
|ExtID (0)|"Pension Chain"|
|ExtID (1)|Pension Company|
|ExtID (2)| Authority PubKey|
|ExtID (3)| nonce |
|Content|Document data

## Transaction

### Value change

| Transaction|Factom Entry|
|---|---|
|ExtID (0)|"Transaction Value Change"|
|ExtID (1)|UserType|
|ExtID (3)|Value Change|
|ExtID (4)|PensionID|
|ExtID (5)|PersonSubmit|
|ExtID (6)| Timestamp|
|ExtID (6)|PubKey|
|ExtID (7)|Signature ExtID(0-6)|
|Content|Document data
