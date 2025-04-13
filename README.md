# Go-SIP

A simple SIP server implemented in Go.

## Features

- Reception and processing of SIP messages
- Handling of REGISTER, INVITE, and BYE requests
- Generation of SIP responses
- Configuration via config file
- Comprehensive test suite
- CI/CD with GitHub Actions

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

## Testing

Run the test suite with:

```
go test ./...
```

This will run all tests in all packages. The test suite includes:

- Unit tests for SIP message parsing and generation
- Mock-based tests for the SIP server
- Configuration file handling tests
- Command-line argument parsing tests

For more verbose output, add the `-v` flag:

```
go test -v ./...
```

For race condition detection:

```
go test -race ./...
```

### Testing Requirements

Tests require a `testdata` directory with test configuration files. This is automatically created when running tests via GitHub Actions, but for local testing you might need to create it:

```bash
mkdir -p testdata config/testdata
echo '{"server":{"port":"5070","log_level":"info","bind_addr":"127.0.0.1"}}' > testdata/test_config.json
echo '{"server":{"port":"5070","log_level":"info","bind_addr":"127.0.0.1"}}' > config/testdata/test_config.json
```

## Continuous Integration

This project uses GitHub Actions for continuous integration:

- Builds and tests the code on multiple Go versions (1.18, 1.19, 1.20)
- Runs linting with golangci-lint
- Generates test coverage reports and uploads to Codecov
- Performs detailed package-by-package testing on failures

You can see the current build status and test coverage at the top of this README.

## Troubleshooting

### Running Tests

If you encounter issues with running tests:

1. Make sure the required test directories exist:

   ```bash
   mkdir -p testdata config/testdata
   ```

2. Create test configuration files:

   ```bash
   echo '{"server":{"port":"5070","log_level":"info","bind_addr":"127.0.0.1"}}' > testdata/test_config.json
   echo '{"server":{"port":"5070","log_level":"info","bind_addr":"127.0.0.1"}}' > config/testdata/test_config.json
   ```

3. Run tests with verbose output:

   ```bash
   go test -v ./...
   ```

4. For more detailed debugging, run tests one package at a time:
   ```bash
   go test -v github.com/user/go-sip/sip
   ```

## Disclaimer

This is a simple SIP server for demonstration purposes and is not recommended for production use. A real SIP server requires additional features such as authentication, security, and full SIP specification support.
