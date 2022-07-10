**Run **
```shell
docker-compose up
```

**Tests**
```shell
go test ./... -v
```

* **Gateway,** exposes api for accounts management
* **User,** account management grpc service is responsible for CRUD operations against user accounts
* **Listener,** will listen for account requests/events (CRUD)