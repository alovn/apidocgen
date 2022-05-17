# apidocgen

apidocgen is a tool for Go to generate apis markdown docs.

## Install

```bash
go install github.com/alovn/apidocgen/@latest
```

## Usage

```bash
$ apidocgen --help
apidocgen is a tool for Go to generate apis markdown docs.

Usage:
  apidocgen --dir= --excludes= --output= --template= --single

Flags:
    --dir:       Search apis dir, comma separated, default .
    --excludes:  Exclude directories and files when searching, comma separated
    --output:    Generate markdown files dir, default ./docs/
    --template:  Custom template files dir.
    --single:    If true, generate a single markdown file, default false
```

run the command in the go module directory.
