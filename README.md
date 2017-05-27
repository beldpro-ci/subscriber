<h1 align="center">subscriber 📃  </h1>

<h5 align="center">Subscribe Emails with MailChimp without confirmation</h5>

<br/>

### Overview

MailChimp doesn't provide a way of subscribing emails to a list directly from a form element without confirmation.
This package (wrapped as a Docker image as well) allows one to have an endpoint and subscribe emails just by providing a MailChimp list id and API key.


<br />

#### Usage 

There are two ways of using `subscriber` out of the box.

1. Using the binary directly

Fetch the binary for your platform:
- `curl -Ls https://my.filla.be/beld/subscriber-linux-amd64.tgz | tar xvz`
- `curl -Ls https://my.filla.be/beld/subscriber-darwin-amd64.tgz | tar xvz`

Or use use `go get`:
- `go get -u github.com/beldpro-ci/subscriber/subscriber` will install `subscriber` in your `$GOPATH/bin`

2. Using the Docker image

- `docker pull beldpro/subscriber`

<br />

Having it from either of these methods, use it by passing the required environment variables (all of them have `SUBSCRIBER` as prefix):

| Environment Variable  | Description |
| ------------- | ------------- |
| `SUBSCRIBER_MAILCHIMP_API_KEY`  | API Key from MailChimp  |
| `SUBSCRIBER_MAILCHIMP_LIST_ID`  | ID of the MailChimp list to subscribe users to  |
| `SUBSCRIBER_MAILCHIMP_URL`  | URL of the MailChimp Datacenter |
| `SUBSCRIBER_PORT`  | Port that the server listens to (default: `8080`)  |


ps.: `SUBSCRIBER_MAILCHIMP_URL` should not contain the `api` version - i.e, it should look like `https://us15.api.mailchimp.com`.


<br />

For more information, run `subscriber --help`:


<br />

```sh
$ subscriber --help
NAME:
   subscriber - Subscribes user to mailchimp

USAGE:
   subscriber [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value     Port to listen to HTTP requests (default: 8080) [$SUBSCRIBER_PORT]
   --api-key value  MailChimp API Key [$SUBSCRIBER_MAILCHIMP_API_KEY]
   --url value      MailChimp's datacenter URL [$SUBSCRIBER_MAILCHIMP_URL]
   --list-id value  Id of the list which users should subscribe to [$SUBSCRIBER_MAILCHIMP_LIST_ID]
   --help, -h       show help
   --version, -v    print the version
```


<br/>


#### Example


```yml
version: '3.2'
services:
  subscriber:
    image: 'beldpro/subscriber'
    environment:
      - SUBSCRIBER_MAILCHIMP_API_KEY=<redacted> 
      - SUBSCRIBER_MAILCHIMP_LIST_ID=<redacted> 
      - SUBSCRIBER_MAILCHIMP_URL=https://us15.api.mailchimp.com 
      - SUBSCRIBER_PORT=8000
```


### LICENSE

``` 
Copyright (C) Beld PRO - All Rights Reserved
Unauthorized copying of any files, via any medium is strictly prohibited
Proprietary and confidential
Written by Ciro S. Costa <ciro9758@gmail.com>, 2017

```

