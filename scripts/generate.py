#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""

"""

from __future__ import print_function

import subprocess
import re

def get_help_page(page=None):
    command = ['./dist/craftctl', 'help']
    if page:
        command.append(str(page))
    return subprocess.check_output(command).strip()


def commands():
    page = 1
    # --- Showing help page 1 of 9 (/help <page>) ---
    header = r'---.*(\d+).*(\d+).*---'
    help_page = get_help_page(page)
    _, max_page = re.findall(header, help_page)[0]
    while page < int(max_page):
        # Trivial: avoid unnecessary request.
        if page > 1:
            help_page = get_help_page(page)
        # Remove header.
        help_page = re.sub(header, '', help_page)
        # Split alternative command forms.
        help_page = help_page.replace('OR', '')
        # Split commands.
        for command in [c for c in help_page.split('/') if c]:
            yield command
        page += 1


def main():
    for c in commands():
        print(c)

if __name__ == '__main__':
    main()
