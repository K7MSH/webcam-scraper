#!/usr/bin/env python2
import os
from time import strftime
import urllib2
import json

abspath = os.path.abspath(__file__)
dname = os.path.dirname(abspath)
os.chdir(dname)


with open("cameras.py.json") as f:
    urls = json.load(f)
# Syntax of json file:
# {
#   "cam1": "url1.jpg",
#   "cam2": "url2.jpg"
# }

directory = "./"

stamp = strftime("%Y%m%d-%H%M%S")
print("Downloading "+str(stamp)+".jpg ...")

#passwd = urllib2.HTTPPasswordMgrWithDefaultRealm()
#passwd.add_password(None, str(url), str(username), str(password))
#auth_handler = urllib2.HTTPBasicAuthHandler(passwd)
#opener = urllib2.build_opener(auth_handler)
##urllib2.install_opener(opener)

for name, url in urls.items():
    folder = os.path.join(directory, name)
    try:
        os.makedirs(folder)
    except:
        pass
    try:
        response = urllib2.urlopen(str(url), timeout=30)
        with open(os.path.join(folder, str(stamp)+'.jpg'),'wb') as output:
            output.write(response.read())
        print("[%s] Download Complete!" % name)
    except:
        print("[%s] Download failed!" % name)

