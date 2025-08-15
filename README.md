# Tyk oauth flow sample

### Setup

Start pre-configured tyk server by `docker-compose up`

Run sample server using `go run ./cmds/server/`

### Oauth flow

1. Server will print out the auth url 
2. Enter the above user to your browser
3. Simulate login process by clicking `Submit` button
4. Browser will display oauth access token
5. Using the above access token to access protected api by 

`curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/someapi`

6. Observe the server log

### Config explained

Tyk is configured with 2 set of api: `myapp` and `myapp_oauth2` on `./conf/tyk/apps/`:
- `myapp`: actual api protected app
- `myapp_oauth2`: dummy app with no backend, only to serve `oauth` endpoints

`myapp` is configured to translate key metadata and forward to upstream service by

```
{
    ...
    "global_headers": {
      "x-anduin-userid": "$tyk_meta.userid"
    }
    ...
}
```

At the beginning of server started, server will automatically create an oauth app in function `createOauthApp` and printout auth url

Server `loginHandler` will handle the pseudo login process and calling tyk `/tyk/oauth/authorize-client` to generate auth code. This auth code will be injected to the redirect url for the user. This process should be implemented within `gondor`

Server `redirectHandler` simulate user/client side to receive auth code and exchange it to an oauth access token by calling tyk `/oauth/token` api. After that, user can use the access token to request protected api at `/api/v1/xxx`

During the request proxy, tyk will check and validate api key given by using and translate it into `x-anduin-userid` header. `gondor` server should be able to handle this header
