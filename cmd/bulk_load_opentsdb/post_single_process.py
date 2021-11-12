#!/usr/bin/env python3
# -*- coding:utf-8 -*-


import http.client
import urllib
import urllib.request
import urllib.error
import json
import time
import random
import traceback



def post_insert(filename):
    base_url = "192.168.1.84"
    headerdata = {'Content-type': 'application/json'}
    
    with open(filename, "r") as fd:
        count = 0
        for line in fd.readlines():
            if count%1000 == 0:                                                      
                print(count, " times: " , time.time())
                time.sleep(2)
                
            try:
                conn = http.client.HTTPConnection(base_url, 4242)
                conn.request(method="POST", url="/api/put?details", body=line, headers=headerdata)
                res = conn.getresponse()
                print(res.read())
                return
            except Exception as e:
                traceback.print_exc()                                                                                                                                                                                                                
                raise e
            count += 1


if __name__ == '__main__':
    datafile = "/data/opentsdb/data/cpu1000_01"
    
    try:
        post_insert(datafile)
    except Exception as e:
        traceback.print_exc()
        raise e
    except KeyboardInterrupt as e:
        exit(0)
        
