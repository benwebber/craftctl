#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

import json
import os
import sys

import requests


def release(token, project, version):
    headers = {
        'Authorization': 'token {}'.format(token),
        'Accept': 'application/vnd.github.v3+json',
    }
    url = 'https://api.github.com/repos/benwebber/{}/releases'.format(project)
    tag = 'v{}'.format(version)
    payload = {
        'tag_name': tag,
        'name': '{} {}'.format(project, version),
    }
    return requests.post(url, headers=headers, data=json.dumps(payload))


def main():
    try:
        project = sys.argv[1]
        version = sys.argv[2]
    except IndexError:
        print('provide project name and version number', file=sys.stderr)
        return 1

    try:
        token = os.environ['GITHUB_API_TOKEN']
    except KeyError:
        print('set GITHUB_API_TOKEN', file=sys.stderr)
        return 1

    resp = release(token, project, version)
    print(json.dumps(resp.json(), indent=2))


if __name__ == '__main__':
    sys.exit(main())
