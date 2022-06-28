# Log Shipper Service

It listens to the given directory. Once there is a new file created in the directory, it reads the logs in the content, parse the logs into json and call [log receiver service](https://github.com/wanghantao11/log-receiver) to create them.

### Author: Hantao Wang

### Getting started
- Create .env file with the configure params in config.go
- Start server by `go run cmd/server/main.go`

