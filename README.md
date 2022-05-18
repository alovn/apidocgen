# apidocgen

apidocgen is a tool for Go to generate apis markdown docs.

## Install

```bash
go install github.com/alovn/apidocgen@latest
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
    --template:  Template name or custom template directory, built-in includes markdown and apidocs, default markdown.
    --single:    If true, generate a single markdown file, default false
```

built-in templates include `markdown`, `apidocs`, default is `markdown`.

run the command in the go module directory.

```bash
cd your-api-service-dir
apidocgen \
    --dir=svc-user,common \
    --output=./docs

apidocgen \
    --dir=svc-user,common \
    --template=apidocs \
    --output=./docs


```

## Template

The built-in includes `markdown` and `apidocs`.

The built-in template `apidocs` is the template for generate website [apidocs](git@github.com:alovn/apidocs.git).

You can also use the custom template:

```bash
apidocgen \
    --dir=svc-user,common \
    --template=/Users/xxx/workspace/apidocs/custom-template-direcoty \
    --output=./docs
```

## Examples

Some examples and generated markdown docs are here: [apidocgen/examples](https://github.com/alovn/apidocgen/tree/main/examples).
