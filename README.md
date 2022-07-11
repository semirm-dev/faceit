**Run**
* Please run RabbitMQ and PostgresDb first so services can have ready connection
```shell
docker-compose up rmq db
```
* After rmq and pg are ready, run services
```shell
docker-compose up gateway account_listener
```

**Tests**
```shell
go test ./... -v
```

* **Gateway,** exposes api for accounts management
* **User,** account management grpc service, responsible for CRUD operations on user accounts
* **Listener,** demo accounts event listener, current implementation only logs event type and affected entity id

> Note: I am the author of rmq package/library

**TODO**
- [ ] finish tests: TestAccountService_GetAccountsByFilter
- [ ] finish tests: handlers_test
- [ ] implement healthcheck for gateway and users service
- [ ] handle rmq connection loss, reconnection (all code is already ready in rmq package/lib, just use it)
- [ ] add some more logs to gateway handlers