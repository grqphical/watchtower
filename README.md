# watchtower - A programmable TCP/UDP server

watchtower is a program that hosts a server and can be programmed to search for specific data within incoming traffic.

## Functional Goals

1. Host a TCP/UDP Server
2. Allow users to filter incoming requests for certain data points
   - Implement regular expression and string searches
3. Add CLI options to control behaviour of server
   - Include control of network protocol
   - Include control of host address and port

## Install

```bash
$ go install github.com/grqphical/watchtower@latest
```

## Basic Usage

Running the base command hosts a basic TCP server on localhost port 8000

```bash
$ watchtower
```

You can specify the host and port with the `a` and `p` flags (address and port respectively)

```bash
$ watchtower -p 2000 -a 0.0.0.0
```

### Filtering

To add terms to search for use the `f` flag

```bash
$ watchtower -s foo
```

To use regex expressions set the envrionment variable `WATCHTOWER_USE_REGEX` to `1`

```bash
export WATCHTOWER_USE_REGEX=1
# NOTE: Regex statements should be in quotes due to how shells process backslashes and other
# special characters
watchtower -s "\d*"
```

## License

watchtower is licensed under the MIT license
