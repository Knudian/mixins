![Github Actions](https://github.com/Knudian/mixins/actions/workflow/go.yaml/badge.svg)

# mixins

This golang package allow users to compose configuration file, that may rely on dynamically or user defined configuration.

## Installation

> Requires: go >= 1.20

```shell
go get github.com/Knudian/mixins
```

## Configuration

For the `mixin`s to be configured, you have to provide it three environment variables:

- `SECRETS_ROOT_DIR`: a directory absolute path. (defaults to `/run/secrets`)
- `MIXIN_PREFIX_ENV_VAR`: the prefix used within your configuration files, for accessing custom environment variables. (defaults to `%env=`)
- `MIXIN_PREFIX_SECRET`: the prefix used when you want to access secrets from within your app. (defaults to `%secret=`)

## Usage

Let's say you have a product using a YAML structure like the following:

```yaml
foo:
  bar: "baz"
  doo: "bee"
```

Your users may instead wish to provide a configuration such as:

```yaml
foo:
  bar: "${SOME_RANDOM_ENV_VAR}"
```

Implementing shell and/or brace expansion may be overkill, only to access a single environment variable.

### With environment vars

Instead, you can have your user configure your product with:

```yaml
foo:
  bar: "%env=SOME_RANDOM_ENV_VAR"
```

### With secrets

Instead of using an environment variable, your user may wish to use _secrets_.

In order to do so, they will only have to

- set the `SECRETS_ROOT_DIR` environment variable to the absolute path housing their secrets
- then configure your product with:

```yaml
foo:
  bar: "%secret=secretName/secretKey"
```

### Licence

This package is following the using the GPL-v3 license.