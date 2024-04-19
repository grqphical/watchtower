# watchtower - A programmable TCP server

watchtower is a program that hosts a TCP server and can be programmed to search for specific data within incoming traffic

## Functional Goals

1. Host a TCP Server
2. Allow users to filter incoming requests for certain data points
   - Implement regular expression and string searches
3. Add CLI options to control behaviour of server
   - Include control of network protocol
   - Include control of host address and port
4. Add option for user to output matched packets to a file
5. Allow user to customize response of server with a file

## Install

```bash
$ go install github.com/grqphical/watchtower@latest
```

## Basic Usage

Running the base command hosts a basic TCP server on localhost port 8000

```bash
$ watchtower
```

You can specify the host and port with the `-a` and `-p` flags (address and port respectively)

```bash
$ watchtower -p 2000 -a 0.0.0.0
```

### Filtering

To add terms to search for use the `-s` flag

```bash
$ watchtower -s foo
```

To use regex expressions set the envrionment variable `WATCHTOWER_USE_REGEX` to `1`

```bash
export WATCHTOWER_USE_REGEX=1
# NOTE: Regex statements should be in quotes due to how shells process backslashes and other
# special characters
$ watchtower -s "\d*"
```

You can set watchtower to write all matches to a file with the `-f` flag

```bash
$ watchtower -s "\d*" -f output.txt
```

### Responses

You can control how watchtower responds to incoming connections by providing a file to the `r` flag

```bash
$ watchtwoer -s "\d*" -r response.html
```

## Changelog

### 0.1.1

- Removed buffer size option, buffer size is now automatically determined

## License

watchtower is licensed under the MIT license
