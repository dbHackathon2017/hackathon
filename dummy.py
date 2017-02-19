import hashlib
import requests
import random
import json
import string
import os

names_file = "names.txt"
names = open(names_file).read().splitlines()

adresses_file = "adresses.txt"
adresses = open(adresses_file).read().splitlines()


def dump_json(data_dict, output_handle):
    """
    Dumps a python dictionary data-structure to JSON structured file
    :param data_dict: python dictionary
    :param output_handle: output file handler to write the JSON-data to
    """
    output_handle.write('{}\n'.format(json.dumps(data_dict, indent=4)))


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
    params = {
        'firstname': genRandomName(0),
        "lastname": genRandomName(1),
        "address": genRandomAdress(),
        "phone": genRandomPhone(),
        "ssn": genRandomNumberString(9),
        "acctnum": genRandomNumberString(10),
    }
    pensionDict = {"request": "makepension", "params": params}
    return pensionDict


def genDummyData():
    if not os.path.exists("dummy"):
        os.system("sudo mkdir dummy")
    for i in range(50):
        f = open('dummy/dummy'+str(i)+'.json', 'w+')
        dump_json(gendummyDict(), f)

if __name__ == "__main__":
    # genDummyData()
    for i in range(35):
        r = requests.post('http://localhost:1337/POST', json=gendummyDict())
        print(r.raw)
