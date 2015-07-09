# Integrating with Wrike!

Talking to the Wrike API.

## Using envconfig to supply tokens

```bash
export WRIKECLIENTPOC_WRIKEBEARER="token_here"
export WRIKECLIENTPOC_WRIKECLIENTID="id"
export WRIKECLIENTPOC_WRIKECLIENTSECRET="secret"
```

The bearer token will be added to the Authorization HTTP header with the format "bearer token_here".

Other items will be used to refresh the token if needed.

## Wrike OAuth2 interaction

https://developers.wrike.com/documentation/api/overview

https://developers.wrike.com/faq/
