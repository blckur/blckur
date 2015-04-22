import os
import re
import requests

base = os.path.dirname(os.path.realpath(__file__))
resp = requests.get('https://api.github.com/repos/blckur/blckur/commits')
assert(resp.status_code == 200)
data = resp.json()
commit = data[0]['sha']

for name in os.listdir(base):
    path = os.path.join(base, name)
    if not os.path.isdir(path):
        continue

    path = os.path.join(path, 'Dockerfile')
    if not os.path.isfile(path):
        continue

    file_data = open(path, 'r').read()

    file_data = re.sub(
        'RUN go get github.com/blckur/blckur # [a-f0-9]',
        'RUN go get github.com/blckur/blckur # %s' % commit,
        file_data)

    with open(path, 'w') as dockerfile:
        dockerfile.write(file_data)
