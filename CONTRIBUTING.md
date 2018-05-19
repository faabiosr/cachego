# Contributing

By participating to this project, you agree to abide our [code of conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

`cachego` is written in [Go](https://golang.org/).

Prerequisites:

* `make`
* [Go 1.9+](https://golang.org/doc/install)

Clone `cachego` from source into `$GOPATH`:

```sh
$ mkdir -p $GOPATH/src/github.com/fabiorphp
$ cd $_
$ git clone git@github.com:fabiorphp/cachego.git
$ cd cachego
```

Install the build and lint dependencies:
```console
$ make depend
```

A good way of making sure everything is all right is running the test suite:
```console
$ make test
```

## Formatting the code
Format the code running:
```console
$ make fmt
```

## Create a commit

Commit messages should be well formatted.
Start your commit message with the type. Choose one of the following:
`feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `revert`, `add`, `remove`, `move`, `bump`, `update`, `release`

After a colon, you should give the message a title, starting with uppercase and ending without a dot.
Keep the width of the text at 72 chars.
The title must be followed with a newline, then a more detailed description.

Please reference any GitHub issues on the last line of the commit message (e.g. `See #123`, `Closes #123`, `Fixes #123`).

An example:

```
docs: Add example for --release-notes flag

I added an example to the docs of the `--release-notes` flag to make
the usage more clear.  The example is an realistic use case and might
help others to generate their own changelog.

See #284
```

## Submit a pull request

Push your branch to your `cachego` fork and open a pull request against the
master branch.
