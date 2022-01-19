# Contributing

If you would like to contribute to the DPK CLI Runtime, you're very welcome to open a GitHub pull request (PR).

## Prerequisites

The make targets will install most needed tools at the appropriate version as-needed. They are installed via
[`asdf`](https://asdf-vm.com/) directory and referenced there so as to avoid interfering with any tools you might
already have installed. If you wish to pre-install these tools, a convenience make target can be used:
`make install-tools`.

*NOTE*: On MacOS, please use current [GNU Coreutils](https://www.gnu.org/software/coreutils/coreutils.html), installed
via [homebrew](https://brew.sh/) (see [here](https://formulae.brew.sh/formula/coreutils) for details, especially the
instructions to correctly configure your `PATH` to include the `gnubin` directory) and [GNU Findutils](https://www.gnu.org/software/findutils/),
installed via [homebrew](https://brew.sh/) (see [here](https://formulae.brew.sh/formula/findutils) for details).

## Installing pre-commit hook

This project comes with a [pre-commit](https://pre-commit.com) configuration and it is highly recommended to install the
respective git hook:

```sh
make install-tool.pre-commit
pre-commit install
pre-commit install -t commit-msg
```

## Running Make targets in a Docker container

You can run any of the Make targets of this project inside of a Docker container using the `ci.docker.run` target:

```sh
make ci.docker.run RUN_WHAT="make build-snapshot"
```

To run an interactive bash shell in the Docker container used in CI, run:

```sh
make ci.docker.run
```

## Upgrading tools

To upgrade all tools used in the repo, run:

```sh
make upgrade-tools
```

## Publishing a new release

Create and push a Git tag to publish a new release. This will start the [Release workflow](https://github.com/mesosphere/dkp-cli-runtime/actions/workflows/release.yaml)
which will create the GitHub release.
