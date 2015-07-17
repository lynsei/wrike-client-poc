# Integrating with Wrike!

Talking to the Wrike API.

## Goals

* Make an AMI with the Go app in it.
* Go app talks to Wrike every five minutes then relays task updates to Slack.
* Be able to test provisioning and the toolchain with Docker instead of waiting for AMIs to build.

## Packer via Docker via Vagrant

`vagrant up` starts a Vagrant box that will install packer.  From there, environment
variables will provide the secret keys.  This is done in the Vagrant VM and packer reads
them in and supplies them to the image.

Supply environment vars by using `wrike-creds.sh` *in your home dir*.  Vagrant will copy that
into the VM and run it to supply the variables.

## Using envconfig to supply tokens

```bash
export WRIKECLIENTPOC_WRIKEBEARER="token_here"
export WRIKECLIENTPOC_WRIKECLIENTID="id"
export WRIKECLIENTPOC_WRIKECLIENTSECRET="secret"
export WRIKECLIENTPOC_WRIKEREFRESHTOKEN="token"
export WRIKECLIENTPOC_SLACKURL="slack url"
```

The bearer token will be added to the Authorization HTTP header with the format "bearer token_here".

Other items will be used to refresh the token if needed.

## Wrike OAuth2 interaction

https://developers.wrike.com/documentation/api/overview

https://developers.wrike.com/faq/
