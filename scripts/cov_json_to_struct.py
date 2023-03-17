#!/usr/bin/python
import re, json

"""Convert http json result to Go struct"""

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

"""Replace the json_data"""
json_data = '{"can_vote": true,"head_block_number":72957803,"head_block_id":"04593f6b752d4505140bab7d133131c0130cf17f"}'

result = []
data = json.loads(json_data)
for k in data:
    struct_item = name_convert_to_camel(k)
    data_type = cov_data_type(data[k])
    result.append("\t"+struct_item+" "+data_type)

print("\n".join(result))

