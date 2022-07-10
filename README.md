**Run**
* Please run RabbitMQ first so services can have ready connection
```shell
docker-compose up rmq
```
* After rmq is ready, run services
```shell
docker-compose up gateway account_listener
```

**Tests**
```shell
go test ./... -v
```

* **Gateway,** exposes api for accounts management
* **User,** account management grpc service is responsible for CRUD operations against user accounts
* **Listener,** demo listener that will listen for account events (CRUD)

> Note: I am the author of rmq package/library