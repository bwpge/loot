# loot

A tool for storing and organizing loot collected during offensive security operations.

## Installation

```sh
go install github.com/bwpge/loot@latest
```

## Usage

### Basics

Create a new loot file in the current directory:

```sh
loot init
```

The resulting `loot.json` is discovered automatically from the current working directory and any parent. Use `-L <path>` to point at a specific file.

Check the status:

```sh
loot status
```

### Add, update, or delete entries

Add an entry:

```sh
loot add hunter2 -t password -c 'web admin'
loot add alice@corp.local -t username -H dc01
```

Common formats (e.g. `user@domain`, `DOMAIN\user`, `user:pass`, etc.) are split into additional entries automatically (use `-n` to disable).

Add each line in a file as an entry:

```sh
loot add -l creds.txt
```

Add file contents as an entry:

```sh
loot add -i id_rsa -t ssh-key
```

Update or remove entries:

```sh
loot update <id> -c 'new comment' -t password,reused
loot rm <id>
loot rm --all
loot rm -f <flag>
```

### View entries

List and filter:

```sh
loot ls
loot ls -t password
loot ls -H dc01
```

Show raw values (useful for piping into other tools):

```sh
loot show <id>
loot show -t password -s ','
```

### Flags

Capture a flag (e.g., HTB machine):

```sh
loot capture 6367c48dd193d56ea7b0baad25b19455e529f5ee --user alice -H 10.10.10.5
loot capture 3da541559918a808c2402bba5012f6c60b27661c --root
```

Use `--admin`/`-a` as a shorthand for `--root Administrator`.

View flags (same as `list --flags`):

```sh
loot flags
```

## Shell completion

Shell completion is available with `loot completion`. Example with zsh:

```sh
source <(loot completion zsh)
```

Use `loot completion <shell> -h` for instructions for your shell.
