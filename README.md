# README
[![CircleCI](https://circleci.com/gh/tmtk75/weque.svg?style=svg)](https://circleci.com/gh/tmtk75/weque)

Weque is a tool to handle webhooks and notifications.
It's a replacement of [hoko](https://github.com/tmtk75/hoko).


## Features
Weque supports next functions.
* Trigger processe with some environment variables
  when receiving webhooks of [GitHub](https://developer.github.com/webhooks/)
  and [Bitbucket](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html).
* Kick process with some environmet variables when receiving [notifications of docker registry](https://docs.docker.com/registry/notifications/).
* TLS communication
  - with cert and key files.
  - via ACME ([Let's Encrypt](https://letsencrypt.org/)).
* Notification to slack with rich attachment.
* Helper commands to develop.
  - Maintain webhook settings of GitHub and Bitbucket.
  - Run handler scripts with payload without receiving webhooks actually.
  - Send notification to slack manually.


## Installation
```
go get -u github.com/tmtk75/weque
```


## Getting Started
Prerequisites
* [ngrok](https://ngrok.com/)

```
export GITHUB_TOKEN=<your github personal access token>
```
Take it here, <https://github.com/settings/tokens>

### Preparing
```
$ ngrok http 9981
...
Forwarding                    https://df431fc9.ngrok.io -> localhost:9981
...
```
ngrok shows a URL for https and memorize it.

Then start a weque process.
```
SECRET_TOKEN=abc123 go run ./cmd/weque/main.go server
...
```

It's ready to receive webhooks.

### Create a webhook setting and receive ping
```
$ weque github create \
        tmtk75/weque \
        https://df431fc9.ngrok.io \
	abc123
```
Repalce `tmtk75/weque` with a repository you have.

You will see some logs appear for receiving a ping just after
you create a webhook if webhook is created properly.


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

