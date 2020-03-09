#!/usr/bin/env /usr/bin/python2.6
from pymongo import Connection
import time,datetime
import os,sys
connection = Connection('127.0.0.1', 27017)
db = connection.admin
db.authenticate('admin','admin')

connection=pymongo.Connection('10.110.16.164',7107)
db = connection['ddd']
def func_time(func):
    def _wrapper(*args,**kwargs):
        start = time.time()
        func(*args,**kwargs)
        print func.__name__,'run:',time.time()-start
    return _wrapper
@func_time
def ainsert(num):
    posts = db.userinfo
    for x in range(num):
        post = {"_id" : str(x),
        "author": str(x)+"Mike",
        "text": "My first blog post!",
        "tags": ["ddd", "dd.cc", "ddd.cto"],
        "date": datetime.datetime.utcnow()}
        posts.insert(post)
if __name__ == "__main__":
    num = sys.argv[1]
    ainsert(int(num))
