# Request Logger Forward Proxy

## Description

This is a forward proxy server that can log all requests and responses to the target server and print them to the console or in browser.

## Usage

__Note:__ For all examples below we assume that the target server is https://jsonplaceholder.typicode.com.

### Executable usage

Proxy format is `<from>::<to>`, where `from` is the address to listen to and `to` is the target address.
By default, the proxy listens on `0.0.0.0:21001`.

With environment variables:
```bash
PROXY_ADDR=https://jsonplaceholder.typicode.com \
./requestlogger
```

Or use flags:
```bash
./requestlogger \
    --proxy=https://jsonplaceholder.typicode.com
```

### Docker compose usage

```yaml
services:
  request_logger:
    image: ghcr.io/revenkroz/request-logger:main
    container_name: proxy
    environment:
      PROXY_ADDR: https://jsonplaceholder.typicode.com
    ports:
        - "21000:21000" # frontend
        - "21001:21001" # proxy
```

## Run

After start, you can open http://localhost:21000 in your browser to see the logs.

Send requests to http://localhost:21001 to see them in the logs.
```shell
curl https://localhost:21001/todos/1
```

## Multiple listen addresses

```bash
./requestlogger \
    --proxy=0.0.0.0:21001::https://jsonplaceholder.typicode.com \
    --proxy=0.0.0.0:21002::https://example.com
```

or with environment variables:
```bash
PROXY_ADDR=0.0.0.0:21001::https://jsonplaceholder.typicode.com,0.0.0.0:21002::https://example.com \
./requestlogger
```

## List of all flags

* `--proxy` - Multiple values, proxy listen address and target address (if empty, `env:PROXY_ADDR` will be used).
* `-faddr` - Frontend listen address (if empty, `env:FRONTEND_ADDR` will be used, default `0.0.0.0:21000`).
* `-maxlogs` - Maximum number of logs to keep in memory (if empty, `env:MAX_LOGS` will be used, default `20`).
* `-stdout` - Print logs to stdout (if empty, `env:USE_STDOUT` will be used, default `false`).
