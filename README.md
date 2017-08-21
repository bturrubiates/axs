# axs -- Simplifying access

*axs* is a small tool to simplify accessing machines.

## Why?

I frequently find myself using ephemeral development setups that are configured
with conveniences such as serial connections to the host, serial connections to
pieces of hardware, ssh accessible BMC, ssh accessible power strips, etc.
*axs* allows defining groups of machines and providing ways to access the
various subsystems available in a hierarchical fashion.

## Installation

```sh
go install github.com/bturrubiates/axs
```

## Configuration

*axs* expects to find the list of machines in a configuration file. By default,
the configuration file lives at `$HOME/.axsrc.json`. The configuration file
could be JSON, YAML, TOML, or anything else that
[Viper](https://github.com/spf13/viper) supports. I personally prefer to use
JSON, but it shouldn't matter.

Access methods should be specified using URL format. Currently only SSH and
telnet are supported.

## Usage

```sh
ben at yggdrasil :: ./axs -h
Usage: ./axs [OPTIONS] target
  -config string
        Config file. (default "$HOME/.axsrc.json")
  -list
        List targets.
  -resolve
        Resolve command.
```

Passing the `-resolve` flag with a target will generate the access command,
but will print it instead of executing it.

## Completions

Completions are only available for `zsh`. Download the `completion/axs.zsh`
file and source it in your `~/.zshrc`.

### Example

Given the following configuration:

```json
{
    "cc1": {
        "bay-a": {
            "host": "ssh://ben@bay-a:22",
            "serial": "telnet://bay-a:23"
        },
        "bay-b": {
            "host": "ssh://root@192.168.0.1",
            "bmc": "ssh://admin@192.168.0.2"
        }
    }
}
```

```sh
ben at yggdrasil :: ./axs -config config.json -resolve cc1.bay-a.host
ssh -p 22 ben@bay-a
```

The target given must resolve to a string, otherwise `axs` will error:

```sh
ben at yggdrasil :: ./axs -config config.json -resolve cc1.bay-a
2017/07/11 22:43:11 target not found
```
