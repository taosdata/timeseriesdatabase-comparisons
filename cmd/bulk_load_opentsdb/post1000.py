#!/usr/bin/env python3
# -*- coding:utf-8 -*-


import http.client
import urllib
import urllib.request
import urllib.error
import json
import time
import traceback

from multiprocessing import Pool


def post_insert(filename):
    base_url = "192.168.1.84"
    headerdata = {'Content-type': 'application/json'}
    
    with open(filename, "r") as fd:
        count = 0
        for line in fd:
            if count % 1000 == 0:                                                      
                print(filename, " : ", count, " times: " , time.time())
                time.sleep(10)
                
            try:
                conn = http.client.HTTPConnection(base_url, 4242)
                conn.request(method="POST", url="/api/put?details", body=line, headers=headerdata)
                res = conn.getresponse()
            except Exception as e:
                traceback.print_exc()
                raise e
            count += 1


if __name__ == '__main__':
    datafile_list = ["/data/opentsdb/data/cpu1000_01", "/data/opentsdb/data/cpu1000_02"
            "/data/opentsdb/data/cpu1000_03", "/data/opentsdb/data/cpu1000_04", "/data/opentsdb/data/cpu1000_05",
            "/data/opentsdb/data/cpu1000_06", "/data/opentsdb/data/cpu1000_07", "/data/opentsdb/data/cpu1000_08",
            "/data/opentsdb/data/cpu1000_09"]
    
    p = Pool(8)
    li = []
    try:
        for datfile in datafile_list:
            res = p.apply_async(post_insert, args=(datfile,))
            li.append(res)
        p.close()
        p.join()
    except Exception as e:
        traceback.print_exc()
        raise e
    except KeyboardInterrupt as e:
        exit(0)
        
