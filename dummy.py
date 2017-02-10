import hashlib
import random
import json


def genRandomHash():
    return hashlib.sha256(str(random.random())).hexdigest()

def genHashList():
    hashList = []
    for i in range(random.randint(0, 20)):
        hashList.append(str(genRandomHash()))
    return hashList

def gendummyDict():
    dummyList = []
    for i in range(1000):
        dummy = {'pensionID': str(genRandomHash()), 'company': 'NestEgg', 'tokens': random.randint(0, 100), 'PublicKey': str(genRandomHash()),'transactions': genHashList()}
        dummyList.append(dummy)
    dummyDict = {"Pensions":dummyList}
    return dummyDict
f = open('dummy.json', 'w')
json.dump(gendummyDict(), f)
#print json.dumps(gendummyDict(), sort_keys=True, indent=4, separators=(',', ': '))
