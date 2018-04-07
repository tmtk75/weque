# README
Weque is a server to receive webhooks written with modern libraries.

## Getting Started
* ngrok

### Preparing
```
$ ngrok http 3000
...

SECRET_TOKEN=abc123 go run ./cmd/weque/main.go server
...
```
TBD

### Create and list
```
$ weque github create \
        tmtk75/weque \
        https://3def21d4.ngrok.io/github \
	abc123
```
TBD


### Receive



## Debug
### tcpflow on MacOS
```
$ tcpflow -i lo0 -c 'port 3000'
```


## Development
```
[0]$ go run ./cmd/weque/main.go serve
...

[1]$ ngrok http 9981
