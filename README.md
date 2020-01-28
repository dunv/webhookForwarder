# webhookForwarder

A simple implementation which consists of two (and a third for testing) components
- a server which needs to be deployed and available in the internet for receiving webhook calls
- a client which can run anywhere and which connects to the server via grpc

Whenever a webhook call on the server is received, it will be forwarded to all connected clients. These clients in turn have a http(s) destination specified where that call is forwarded to.

## Example
![diagram](https://raw.githubusercontent.com/dunv/webhookForwarder/master/diagram.png)

### server
* `-i` `--incomingSocket` (default: `0.0.0.0:8080`) where to listen to for webhook calls
* `-p` `--incomingPath` (default: `/`) limit paths for forwarding
* `-o` `--outgoingSocket` (default: `0.0.0.0:50051`) where to listen to for client connections
* `-d` `--printHttpDump` (default: `false`) dump every webhook call to stdout
* e.g. `./webhookForwarder server -i 0.0.0.0:8080 -o 0.0.0.0:50051 -d true`

### client
* `-i` `--incomingSocket` (default: `0.0.0.0:50051`) which server to connect to
* `-o` `--outgoingURI` (default: `http://0.0.0.0:8081`) where to forward webhook calls to
* `-p` `--overridePath` (default: ``) override path-segment (this way the server can be called with whatever path, here we can override it)
* e.g. `./webhookForwarder client -i webhookServer.deployment.com:50051 -o http://localhost:8050/staticPathPrefix`

### dummy
* the dummy just listens to everything and dumps all http-calls to stdout
* `-o` `--incomingSocket` (default: `0.0.0.0:8081`) which socket to listen to
* e.g. `./webhookForwarder dummy -i 0.0.0.0:8050`
