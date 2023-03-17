#!/usr/bin/python
# -*- coding: UTF-8 -*-
import re, json, sys

"""
Convert http json result to Go struct
usage: ./cov_json_to_struct.py '{"head_block_number":72958298,"head_block_id":"0459415adbcbb2162989d536daf3f84580932afa"}'
"""

def name_convert_to_camel(name: str) -> str:
    name = name.capitalize()
    return re.sub(r'(_[a-z])', lambda x: x.group(1)[1].upper(), name)

def cov_data_type(t) -> str:
    tmp = type(t)
    if tmp == type(1):
        return 'uint'
    elif tmp == type('a'):
        return 'string'
    elif tmp == type(True):
        return 'bool'
    else:
        return 'unknown'

try:
    json_data = sys.argv[1]
except e:
    print('unexpect error')
    print(e)

result = []
data = json.loads(json_data)
for k in data:
    struct_item = name_convert_to_camel(k)
    data_type = cov_data_type(data[k])
    result.append("\t"+struct_item+" "+data_type)

print("")
print("")
print("")
print("")
print("The result is:")
print("\n".join(result))

