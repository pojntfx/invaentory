# invaentory

Quickly find all IPv6 and IPv4 hosts in a LAN.

[![hydrun CI](https://github.com/pojntfx/invaentory/actions/workflows/hydrun.yaml/badge.svg)](https://github.com/pojntfx/invaentory/actions/workflows/hydrun.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/invaentory.svg)](https://pkg.go.dev/github.com/pojntfx/invaentory)
[![Matrix](https://img.shields.io/matrix/invaentory:matrix.org)](https://matrix.to/#/#invaentory:matrix.org?via=matrix.org)
[![Binary Downloads](https://img.shields.io/github/downloads/pojntfx/invaentory/total?label=binary%20downloads)](https://github.com/pojntfx/invaentory/releases)

## Overview

🚧 This project is a work-in-progress! Instructions will be added as soon as it is usable. 🚧

## Installation

Static binaries are available on [GitHub releases](https://github.com/pojntfx/invaentory/releases).

On Linux, you can install them like so:

```shell
$ curl -L -o /tmp/invaentory "https://github.com/pojntfx/invaentory/releases/latest/download/invaentory.linux-$(uname -m)"
$ sudo install /tmp/invaentory /usr/local/bin
$ sudo setcap cap_net_raw+ep /usr/local/bin/invaentory # This allows rootless execution
```

On macOS, you can use the following:

```shell
$ curl -L -o /tmp/invaentory "https://github.com/pojntfx/invaentory/releases/latest/download/invaentory.darwin-$(uname -m)"
$ sudo install /tmp/invaentory /usr/local/bin
```

On Windows, the following should work (using PowerShell as administrator):

```shell
PS> Invoke-WebRequest https://github.com/pojntfx/invaentory/releases/latest/download/invaentory.windows-x86_64.exe -OutFile \Windows\System32\invaentory.exe
```

You can find binaries for more operating systems and architectures on [GitHub releases](https://github.com/pojntfx/invaentory/releases).

## Usage

🚧 This project is a work-in-progress! Instructions will be added as soon as it is usable. 🚧

Be sure to check out the [reference](#reference) for more information.

## Reference

```bash
$ invaentory --help
Usage of invaentory:
  -4    Ping using ICMPv4 (default true)
  -6    Ping using ICMPv6 (default true)
  -exclude string
        Regex of addresses to exclude
  -multicast-timeout int
        Time in milliseconds to wait for responses for multicast (IPv6) pings (default 2000)
  -parallel int
        Amount of pings to run in parallel (default 100)
  -progress
        Show progress bar on STDERR (default true)
  -unicast-timeout int
        Time in milliseconds to wait for responses for unicast (IPv4) pings (default 500)
  -verbose
        Enable verbose logging to STDERR
```

## Acknowledgements

- This project would not have been possible were it not for [@digineo](https://github.com/digineo)'s [go-ping](https://github.com/digineo/go-ping) package; be sure to check it out too!
- All the rest of the authors who worked on the dependencies used! Thanks a lot!

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

To build invaentory locally, run:

```shell
$ git clone https://github.com/pojntfx/invaentory.git
$ cd invaentory
$ make depend
$ make
$ sudo setcap cap_net_raw+ep out/invaentory
$ out/invaentory
```

Have any questions or need help? Chat with us [on Matrix](https://matrix.to/#/#invaentory:matrix.org?via=matrix.org)!

## License

invaentory (c) 2022 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
