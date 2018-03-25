# README

## Getting Started
* consul
* ngrok

### Preparing
```
$ ngrok http 3000
...

$ make consul
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