# README
[![CircleCI](https://circleci.com/gh/tmtk75/weque.svg?style=svg)](https://circleci.com/gh/tmtk75/weque)

Weque is a tool to handle webhooks and notifications.
It's a replacement of [hoko](https://github.com/tmtk75/hoko).

Weque supports next functions.
* Trigger processe when receiving webhooks of [GitHub](https://developer.github.com/webhooks/) and [Bitbucket](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html).
* Kick process when receiving [notifications of docker registry](https://docs.docker.com/registry/notifications/).
* TLS communication
  - with cert and key files.
  - via ACME ([Let's Encrypt](https://letsencrypt.org/)).
* Notification to slack


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
```


## Credit
<div>The GitHub icon made by <a href="https://www.flaticon.com/authors/dave-gandy" title="Dave Gandy">Dave Gandy</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>

<div>The Bitbucket icon made by <a href="https://www.flaticon.com/authors/swifticons" title="Swifticons">Swifticons</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>

