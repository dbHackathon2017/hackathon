import hashlib
import random
import json
import string
import datetime
import csv

names_file = "names.txt"
names = open(names_file).read().splitlines()

adresses_file = "adresses.txt"
adresses = open(adresses_file).read().splitlines()


def genRandomHash():
    return hashlib.sha256(str(random.random())).hexdigest()

def genHashList():
    hashList = []
    for i in range(random.randint(0, 20)):
        hashList.append(str(genRandomHash()))
    return hashList

def genRandomName(n):
    name = names[random.randint(0,len(names)-1)].split()
    return name[n]

def genRandomPhone():
    return"6"+''.join(random.SystemRandom().choice(string.digits) for _ in range(10))

def genRandomNumberString(n):
    return ''.join(random.SystemRandom().choice(string.digits) for _ in range(9))

def genRandomAdress():
    return adresses[random.randint(0,len(adresses)-1)]+" "+str(random.randint(0,300))


def gendummyDict():
    dummyList = []
    for i in range(50):
        docs = {"path": "check.pdf",
                "hash":genRandomHash(),
                "timestamp": str(datetime.datetime.utcnow()),
                "source":"NestEgg" ,
                "location": "NestEgg"}
        params = {'firstname': genRandomName(0),
            "lastname":genRandomName(1),
            "address":genRandomAdress(),
            "phone":genRandomPhone(),
            "BSN":genRandomNumberString(9),
            "AccountNumber":genRandomNumberString(10),
            "docs": docs}
        pensionDict = {"request": "pension", "params":params}
        dummyList.append(pensionDict)
    dummyDict = {"pensions": dummyList}
    return dummyDict
f = open('dummy.json', 'w')
json.dump(gendummyDict(), f)
print json.dumps(gendummyDict(), sort_keys=True, indent=4, separators=(',', ': '))

#print genRandomName(0),genRandomName(1)
#print gendummyDict()


"""{
   "request":"pension",
   "params":{
      "firstname":"FIRST_NAME",
      "lastname":"LAST_NAME",
      "address":"ADDRESS",
      "phone":"PHONE_NUM",
      "ssn":"SSN",
      "acctnum":"ACCT_NUM",
      "docs":[
         {
            "path":"FILE_PATH",
            "hash":"FILE_HASH",
            "timestamp":"TIME_STAMP", // formatted string
            "source":"WHO_MADE",
            "location":"WHERE"
         }
      ]
   }
}"""
