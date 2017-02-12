import hashlib
import requests
import random
import json
import string
import datetime
import os

names_file = "names.txt"
names = open(names_file).read().splitlines()

adresses_file = "adresses.txt"
adresses = open(adresses_file).read().splitlines()


def genRandomHash():
    return hashlib.sha256(str(random.random()).encode('utf-8')).hexdigest()


def genRandomName(n):
    name = names[random.randint(0, len(names)-1)].split()
    return name[n]


def genRandomPhone():
    return "6" + ''.join(random.SystemRandom().choice(string.digits) for _ in range(10))


def genRandomNumberString(n):
    return ''.join(random.SystemRandom().choice(string.digits) for _ in range(9))


def genRandomAdress():
    return adresses[random.randint(0,len(adresses)-1)]+" "+str(random.randint(0,300))


def gendummyDict():
    docs = [
        {
            "path": "check.pdf",
            "hash": genRandomHash(),
            "timestamp": str(datetime.datetime.utcnow()),
            "source": "NestEgg",
            "location": "NestEgg"
        }
    ],
    params = {
        'firstname': genRandomName(0),
        "lastname": genRandomName(1),
        "address": genRandomAdress(),
        "phone": genRandomPhone(),
        "ssn": genRandomNumberString(9),
        "acctnum": genRandomNumberString(10),
        "docs": docs
    }
    pensionDict = {"request": "makepension", "params": params}
    return pensionDict


def genDummyData():
    if not os.path.exists("dummy"):
        os.system("sudo mkdir dummy")
    for i in range(50):
        f = open('dummy/dummy'+str(i)+'.json', 'w')
        json.dump(gendummyDict(), f)


if __name__ == "__main__":
    # genDummyData()
    r = requests.post('http://localhost:1337/POST', json=gendummyDict())
    print(r)
