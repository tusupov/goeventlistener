# Implement Events Listeners in Go

## Params
* `PORT` - address for server listen, default 8080

## Run width docker
``` bash
$ docker-compose up
```
## API
POST [`/listener`](#api-create) - Create new listener\
DELETE [`/listener/{listenerName}`](#api-delete) - Delete listener\
POST [`/publish/{eventName}`](#api-publish) - Publish listeners by event name

### <a name="api-create"></a>Create new listener
```
POST /listener
{
  "event": "eventName",
  "name": "listenerName",
  "address": "listenerHttpAddress"
}
```

### <a name="api-delete"></a> Delete listener
```
DELETE /listener/{listenerName}
```

### <a name="api-delete"></a> Publish listeners by event name
```
POST /publish/{eventName}
```

## Example
```

==> POST /listener
{
    "event": "event1",
    "name": "listener1",
    "address": "https://golang.org"
}
curl -i http://localhost:8080/listener -X POST -d '{"event": "event1", "name": "listener1","address": "https://golang.org"}'
<== HTTP 200 OK
{"success": "ok"}

==> POST /listener
{
    "event": "event2",
    "name": "listener2",
    "address": "https://golang.org"
}
curl -i http://localhost:8080/listener -X POST -d '{"event": "event2", "name": "listener2","address": "https://golang.org"}'
<== HTTP 200 OK
{"success": "ok"}

==> POST /listener
{
    "event": "event2",
    "name": "listener1",
    "address": "https://golang.org"
}
curl -i http://localhost:8080/listener -X POST -d '{"event": "event2", "name": "listener1","address": "https://golang.org"}'
<== HTTP 400 Bad Request
{"error": "Listener exists."}

==> DELETE /listener/listener1
curl -i http://localhost:8080/listener/listener1 -X DELETE
<== HTTP 200 OK
{"success": "ok"}

==> DELETE /listener/listener1
curl -i http://localhost:8080/listener/listener1 -X DELETE
<== HTTP 400 Bad Request
{"error": "Listener not exists."}

==> POST /publish/event2
curl -i http://localhost:8080/publish/event2 -X POST
<== HTTP 200 OK
{"success": "ok"}

==> POST /publish/event3
curl -i http://localhost:8080/publish/event3 -X POST
<== HTTP 400 Bad Request
{"error": "Event not exists."}

==> POST /publish/event1
curl -i http://localhost:8080/publish/event1 -X POST
<== HTTP 400 Bad Request
{"error": "Event listeners not exists."}
```
