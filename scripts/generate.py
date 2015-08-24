#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""

"""

from __future__ import print_function

import subprocess
import re

formatter = '   {}'.format


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
        # Remove header to isolate commands.
        page_commands = re.sub(header, '', help_page)
        # Split commands on /, except where the / presents alternate syntax,
        # e.g., `tp`.
        for line in re.sub(r'(\S)/(\S)', r'\1\n\2', page_commands).strip().split('\n'):
            # Remove any / from the middle of the command.
            yield line.replace('/', '')
        page += 1


def main():
    for c in commands():
        print(formatter(c))

if __name__ == '__main__':
    main()
