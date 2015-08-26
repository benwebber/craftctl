#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Generate command list for craftctl.
"""

from __future__ import print_function

import argparse
import subprocess
import sys
import re

template = """package main

// THIS IS GENERATED CODE. DO NOT EDIT.

import "github.com/codegangsta/cli"

func init() {{
	cli.AppHelpTemplate = `Usage: {{{{.Name}}}} [options] <command> [args]...

   {{{{.Usage}}}}
{{{{if .Flags}}}}
Options:
   {{{{range .Flags}}}}{{{{.}}}}
   {{{{end}}}}{{{{end}}}}
Commands:
   {commands}
`
}}
""".format


def get_help_page(page=None):
    command = ['./dist/craftctl', 'help']
    if page:
        command.append(str(page))
    return subprocess.check_output(command).strip()


def get_commands():
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


def parse_args(argv):
    """
    Parse command-line arguments.

    Returns: argparse.Namespace object.
    """
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        '-i', '--input', type=argparse.FileType('r'),
        help='File containing commands. If not specified, generate commands using built-in help command.',
    )
    parser.add_argument(
        '-o', '--output', type=argparse.FileType('w'), default=sys.stdout,
        help='Where to write output. If not specified, default to standard output.',
    )
    parser.add_argument(
        '--format', choices=('text', 'cli.go'), default='text',
        help='Output format: choose between plaintext and template for use with cli.go.'
    )
    return parser.parse_args(argv)


def main(argv=None):
    if not argv:
        argv = sys.argv[1:]
    args = parse_args(argv)

    if args.input:
        commands = args.input.read().splitlines()
    else:
        commands = get_commands()

    if args.format == 'text':
        args.output.write('\n'.join(commands))
    else:
        data = {
            'commands': '\n   '.join(commands)
        }
        args.output.write(template(**data))


if __name__ == '__main__':
    main()
