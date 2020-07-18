from swpag_client import Team
import sys
from time import sleep
import subprocess



TICK=30
MAX_SIZE=100
exploit = sys.argv[2]
service= sys.argv[1]


t = Team("http://teaminterface.ictf.love/", "g7iCTu9Gt6pj1DCG4XwP")
services = t.get_service_list()
print(services)
if service not in services:
    raise Exception("Check service name")

while True:
    targets= ["diocane"]
    targets = t.get_targets(service)
    for target in targets:
        flags=subprocess.check_output([exploit,target]).decode("utf-8").split("\n")
        if len(flags)> MAX_SIZE:
            list_of_flags=[flags[:MAX_SIZE],flags[MAX_SIZE:]]
        else:
            list_of_flags=[flags]
        print(list_of_flags)
        for max_flags in list_of_flags:
            t.submit_flags(max_flags)



