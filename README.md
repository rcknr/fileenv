# About

`fileenv` is a simple program to aid running [twelve-factor apps](https://12factor.net/config) in Docker (Swarm).

# TL;DR

* You need `YOURENVVAR` to hold your Docker secret, `<secret-name>`
* Set `YOURENVVAR_FILE=/run/secrets/<secret-name>`
* Set your Dockerfile CMD as `/path/to/fileenv <your cmd> <your args...>`
* Profit!

# Why

When using Docker Swarm, it's desirable to use [Docker secrets](https://docs.docker.com/engine/swarm/secrets/) to manage secrets [instead of passing secrets in environment variables](https://github.com/moby/moby/issues/13490).

Many twelve-factor apps don't support reading from Docker secrets because they only read directly from environment variables.

`fileenv` is a stop-gap solution to read Docker secrets (or any file) and set environment variables for your app, until you can update your app to read from Docker secrets.

# Usage

```
Usage: fileenv [flags...] /path/to/program [program arguments...]
Flags:
  -debug
    	print debug info messages
  -fail
    	immediately exit on warning
```

**Note: `-debug` will print the contents of variables it sets. If you don't want secrets printed in your logs, don't use `-debug`!**

`fileenv` will read through each environment variable. If the name of the environment variable ends with `_FILE` when upper-cased, `fileenv` will open the file referenced by the variable, read its contents, and set the enviroment variable without `_FILE` at the end to the [TrimSpace](https://golang.org/pkg/strings/#TrimSpace)'d contents of the file.

If any of these steps fails, a warning will be printed to `stderr`. If `-fail` is set, `fileenv` will immediately exit with status `1`. Otherwise `fileenv` will continue on.

Once all environment variables are set, `fileenv` will execute the given program and arguments, passing through `stdin`, `stdout`, and `stderr`. If the program fails to start, or returns a non-zero exit status, the error will be logged to `stderr`.
