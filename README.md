# Kore
#### Put your bot in *all* the chats

A pluggable bot that allows you to easily put your bot on multiple chat platforms using a single session. Written in Go.

## Building from source

To build from source:
```
go get -u github.com/hegemone/kore
cd $GOPATH/src/github.com/hegemone/kore
make
```

To run, you can execute `make run` in the root directory of the project.

While there is a Dockerfile in the repo, it is not currently updated to reflect the current state of the project.

**IMPORTANT**: When building plugins, go will attempt to install a number of
utilities to your $GOROOT location. It must be writable by your user to allow
for those tools to be installed on a first run. Also, plugins require you use
Go v1.8 or higher.

Important make commands:
* `make plugins` builds the example bacon plugin as a `.so` libs in `build/`
* `make adapters` builds the example adapters as `.so` libs in `build/`
* `make kore` builds the executable as `korecomm` in `build/`, and depends on
`plugins` and `adapters` targets.
* `make run` sets up extension load paths via env vars and runs the executable.
* `make clean` cleans the `build/` directory.
* `make image` will build a Docker image from source in your local registry.
* `make` by default will run `make build`, which is an alias for `kore`.

Change for testing
