# Request Logger Reverse Proxy

## Description

This is a reverse proxy server that can log all incoming requests to the target server and print them to the console or wev view.

## Usage

__Note:__ For all examples below we assume that the target server is https://jsonplaceholder.typicode.com.

### Executable usage

Add environment variables:
```bash
TARGET_URL=https://jsonplaceholder.typicode.com
```

And run.
```bash
./requestlogger
```

Or use flags:
```bash
./requestlogger \
    -target=https://jsonplaceholder.typicode.com
```

#### List of all flags

* `-target` - Target URL (if empty, `env:TARGET_URL` will be used).
* `-addr` - Proxy listen address (if empty, `env:LISTEN_ADDR` will be used, default `0.0.0.0:21000`).
* `-faddr` - Frontend listen address (if empty, `env:FRONTEND_ADDR` will be used, default `0.0.0.0:21001`).
* `-maxlogs` - Maximum number of logs to keep in memory (if empty, `env:MAX_LOGS` will be used, default `20`).
* `-stdout` - Print logs to stdout (if empty, `env:USE_STDOUT` will be used, default `false`).

### Docker compose usage

```yaml
services:
  request_logger:
    image: ghcr.io/revenkroz/request-logger:main
    container_name: proxy
    environment:
      TARGET_URL: https://jsonplaceholder.typicode.com
    ports:
        - "21000:21000"
        - "21001:21001"
```

## Run

After start, you can open http://localhost:21001 in your browser to see the logs.

Send requests to http://localhost:21000 to see them in the logs.
```shell
curl https://localhost:21000/todos/1
```
