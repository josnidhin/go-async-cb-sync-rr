The sequence diagrams are generated using https://sequencediagram.org/

* Source for `API Service with Callbacks` Diagram
```
title API Service with Callbacks

participant "Client" as c
participant "API Service" as api

activate c
c->api:Do Task
activate api
activate api
api-->c:Order Created (ID)
deactivate api
deactivate c

api->c: Callback
activate c
c-->api:Ok
deactivate c
deactivate api
```

* Source for `API Service with Callback to Request-Reply`
```
title API Service with Callback to Request-Reply

participant "Client" as c
participant "Adapter Service" as a
participant "API Service" as api

activate c
c->a:Do Task
activate a
a->api: Do Task with callback
activate api
activate api
api-->a:Order Created (ID)
deactivate api


api->a: Callback
activate a
a-->api:Ok
deactivate a
deactivate api
a-->c: Task done
deactivate a
deactivate c
```
