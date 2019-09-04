# golang-grpc-plugin-test

## Use

### Start a simple (non-HTTPS) webserver

```
$ python3 -m http.server 8080 --bind 0.0.0.0
```
Leave this process running - it is the endpoint for the broker

### Start the broker

- Go into `app/` folder
- Run `./run.sh` to start the broker

### Retrieve data through the broker using curl

```
$ curl localhost:9090/
```
