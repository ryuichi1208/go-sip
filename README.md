# Go-SIP

A simple SIP server implemented in Go.

## Features

- Reception and processing of SIP messages
- Handling of REGISTER, INVITE, and BYE requests
- Generation of SIP responses
- Configuration via config file

## Usage

### Starting the Server

Start the server with the following command:

```
go run main.go
```

By default, the server listens on UDP port 5060.

### Configuration File

The server reads a configuration file named `config.json` by default. You can specify a different configuration file:

```
go run main.go -config myconfig.json
```

To generate a default configuration file:

```
go run main.go -generate-config
```

Example configuration file:

```json
{
  "server": {
    "port": "5060",
    "log_level": "info",
    "bind_addr": "0.0.0.0"
  }
}
```

### Command Line Options

Override configuration file values with command line options:

```
go run main.go -port 5080 -bind 127.0.0.1
```

All options:

- `-config <file>` - Path to configuration file
- `-generate-config` - Generate default config file and exit
- `-port <port>` - Override port number
- `-bind <addr>` - Override bind address

## Test Client

A simple SIP client is included for testing:

```
go run examples/sip_client.go
```

The client supports the following commands:

- register: Register a user
- invite: Initiate a call
- bye: End a call
- exit: Exit the client

You can specify the server address and username:

```
go run examples/sip_client.go -server 127.0.0.1:5060 -user alice
```

## Supported SIP Methods

- REGISTER: User registration
- INVITE: Call initiation
- BYE: Call termination
- ACK: Acknowledgment handling

## Disclaimer

This is a simple SIP server for demonstration purposes and is not recommended for production use. A real SIP server requires additional features such as authentication, security, and full SIP specification support.
