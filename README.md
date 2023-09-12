# invaentory

Quickly find all IPv6 and IPv4 hosts in a LAN.

[![hydrun CI](https://github.com/pojntfx/invaentory/actions/workflows/hydrun.yaml/badge.svg)](https://github.com/pojntfx/invaentory/actions/workflows/hydrun.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/invaentory.svg)](https://pkg.go.dev/github.com/pojntfx/invaentory)
[![Matrix](https://img.shields.io/matrix/invaentory:matrix.org)](https://matrix.to/#/#invaentory:matrix.org?via=matrix.org)
[![Binary Downloads](https://img.shields.io/github/downloads/pojntfx/invaentory/total?label=binary%20downloads)](https://github.com/pojntfx/invaentory/releases)

## Overview

`invaentory` is a very fast IPv6 and IPv4 LAN scanner, with the intention of allowing the user to quickly take inventory of hosts on their network.

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

To take inventory of all the IPv6 and IPv4 hosts in your LAN, simply run `invaentory`. Note that while Linux allows for rootless execution, macOS requires the use of `sudo` and Windows the use of PowerShell as administrator:

```shell
$ invaentory
172.17.0.1
Pinging   7% [>                           ] (4612/65553, 200 host/s) [23s:5m4s]
```

Once the scan is finished, all found hosts will be listed:

```shell
$ invaentory
172.17.0.1
100.64.154.241
2001:7c7:2121:8d00:40e3:c0ea:d71c:db75
100.64.154.242
2001:7c7:2121:8d00:da47:32ff:fec9:62a0
100.64.154.244
2001:7c7:2121:8d00::3
100.64.154.243
fe80::b2a8:6eff:fe0c:ed1a%enp0s13f0u1u2u2
100.64.154.254
100.64.154.246
100.64.154.250
2001:7c7:2121:8d00:8c77:a50a:e3a9:d284
100.64.154.250
2001:7c7:2121:8d00:d125:4d82:b9f7:5a00
100.64.154.247
```

It is also possible to exclude certain IPs by supplying a regular expression to `--exclude`:

```shell
$ invaentory --exclude '100..*'
172.17.0.1
2001:7c7:2121:8d00:40e3:c0ea:d71c:db75
2001:7c7:2121:8d00:da47:32ff:fec9:62a0
2001:7c7:2121:8d00::3
fe80::b2a8:6eff:fe0c:ed1a%enp0s13f0u1u2u2
2001:7c7:2121:8d00:8c77:a50a:e3a9:d284
2001:7c7:2121:8d00:d125:4d82:b9f7:5a00
```

The progress bar and logging output is written to `STDERR`, so you can further process just the IP addresses as usual, for example to scan all the IPv6 nodes using `nmap`:

```shell
$ for host in $(invaentory -4=false); do nmap -6 ${host}; done
Starting Nmap 7.91 ( https://nmap.org ) at 2022-01-15 02:42 CET
Nmap scan report for felicias-xps13 (2001:7c7:2121:8d00:40e3:c0ea:d71c:db75)
Host is up (0.000089s latency).
Not shown: 999 closed ports
PORT   STATE SERVICE
80/tcp open  http

Nmap done: 1 IP address (1 host up) scanned in 0.09 seconds
Starting Nmap 7.91 ( https://nmap.org ) at 2022-01-15 02:42 CET
Nmap scan report for felicitass-proton.user.selfnet.de (2001:7c7:2121:8d00:da47:32ff:fec9:62a0)
Host is up (0.00027s latency).
Not shown: 996 closed ports
PORT     STATE SERVICE
22/tcp   open  ssh
80/tcp   open  http
443/tcp  open  https
9090/tcp open  zeus-admin

Nmap done: 1 IP address (1 host up) scanned in 0.08 seconds
Starting Nmap 7.91 ( https://nmap.org ) at 2022-01-15 02:42 CET
# ...
```

🚀 **That's it!** You can now take inventory of your network.

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

invaentory (c) 2023 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
