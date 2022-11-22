# fun-run

[![Go](https://github.com/a-poor/fun-run/actions/workflows/go.yml/badge.svg)](https://github.com/a-poor/fun-run/actions/workflows/go.yml)
[![goreleaser](https://github.com/a-poor/fun-run/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/a-poor/fun-run/actions/workflows/goreleaser.yml)

_created by Austin Poor_

A simple CLI for executing multiple processes simultaneously.


<p align="center">
    <img src="./demo.gif" width="640" />
</p>


## Installation

`fun-run` can be installed with go...

```sh
go install github.com/a-poor/fun-run@latest
```

Or pre-built binaries can be downloaded from the [releases](https://github.com/a-poor/fun-run/releases/latest).


## CLI Usage

<details>
<summary><b>CLI help:</b>

```sh
fun-run --help
```
</summary>

```
Fun Run is a tool for running multiple processes simultaneously.

It is designed to be used in development environments where you want
to run multiple processes (e.g. a web server and a database server)
simultaneously. It is similar to the 'docker-compose' tool, but for
running shell commands rather than containers.

Usage:
  fun-run [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize a new fun-run config file.
  run         Start running your commands.
  validate    Validate the configuration file

Flags:
  -h, --help      help for fun-run
  -V, --verbose   Help message for toggle
  -v, --version   version for fun-run

Use "fun-run [command] --help" for more information about a command.
```
</details>


<details>
<summary><b>Initialize a fun-run config file</b>

```sh
fun-run init -
```
</summary>

```yaml
procs:
    - name: say-hello
        cmds:
        - echo hello...
        - sleep 1
        - echo ...world
        - sleep 1
        restart: always
    - name: print-the-date
        cmd: date
        restart: never
    - name: greet-fun-run
        cmd: echo
        args:
        - Hello, ${NAME}!
        envs:
        NAME: fun-run
        restart: on-fail
```
</details>

<details>
<summary><b>Validate a fun-run config file</b>

```sh
fun-run validate fun-run.example.yaml
```
</summary>

```
Config file is valid!
```
</details>

<details>
<summary><b>Run processes in a fun-run config file</b>

```sh
fun-run run fun-run.example.yaml
```
</summary>

```
print-the-date Starting...
proc-0 Starting...
proc-2 Starting...
proc-3 Starting...
print-the-date | Tue Nov 22 13:37:59 PST 2022
print-the-date Finished
proc-3 | hello, fun-run!
proc-3 Finished
proc-2 | starting...
proc-0 Finished
proc-2 | continuing...
proc-2 | done
proc-2 Finished
```
</details>

## Config File Format

The config file should have a root key `procs` which is a list of
process configs.

Each process config minimally needs a `cmd` (for a single command)
or `cmds` (for multiple commands).

| Key | Type | Description |
| --- | --- | --- |
| `cmd` | `string` | (Required unless `cmds` is set) Executable to run |
| `cmds` | `[]string` | (Required unless `cmd` is set) Array of commands to run (`sh -c "{{ cmds joined with ';' }}"`) |
| `name` | `string` | Name of the process (defaults to `proc-{{ index }}`) |
| `args` | `[]string` | Arguments to pass to the `cmd` or `cmds` |
| `envs` | `map[string]string` | Environment variables to pass to the command |
| `clear_envs` | `bool` | Should the command get the env vars in addition to `envs`? |
| `workdir` | `string` | Working directory from which to run the command (default: `.`) |
| `restart` | `string` | `never`: never restart, `on-fail`: only restart on failure, `always`: always restart when stopped |
