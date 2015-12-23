# gopath

A wrapper for go related command which set GOPATH dynamically.

# Overview

For some reason, one needs to build/get/run go programs out of his
main $GOPATH occasionally, or do not want to set an all-in-one static
$GOPATH at all.

This program can set $GOPATH to the first directory it found from
current working directory towards root which contains a subdirectory
called *src*. If there is already a $GOPATH in environment, gopath
will respect it and do nothing. The following table explains the
effect of $GOPATH according to $PWD.

| $PWD                 | $GOPATH          |
| -------------------- | ---------------- |
| /path/to/project/src | /path/to/project |
| /tmp/proj/src/foo    | /tmp/proj        |

# Installation

## Build

```shell
go get go.papla.net/gopath
```
gopath is go-getable, you just need to put it into your $PATH after
build.

## Replace the original binary

gopath will treat argv[0] as the original name if it not named
"gopath", otherwise argv[1]. So you can either make a symlink or use
shell alias.

e.g., using symlink:
```shell
ln -s /path/to/gopath /path/go
```

and using shell alias:
```shell
alias go = 'gopath go'
```

## Setup origin

Two methods are supported currently, using a config file(recommended)
or rename the origin.

1. Modify gopathrc.sample and save it into one of the following location:
   - $XDG_CONFIG_HOME/gopathrc
   - $HOME/.config/gopathrc
   - $HOME/.gopathrc

2. rename the origin *foo* to *foo.bin*, e.g. `mv go go.bin`.


# Issue

Do not support windows yet.
