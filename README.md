
### Prerequisites

What things you need to install the software and how to install them

```
Give examples
Golang v1.16
Go Mod
....
```

### Installing

A step by step series of examples that tell you have to get a development env running

Say what the step will be
- Create ENV file (.env) with this configuration:
```
APP_NAME=wallet-service
PORT=9000
ROLLING_PERIOD=180
THRESHOLD=10000
KAFKA_BROKERS=localhost:9092
KAFKA_SSL_ENABLE=false
KAFKA_USERNAME=
KAFKA_PASSWORD=
```
- Then run this command (Development Issues)
```
Give the example
...
$ make run-dev
```

- Then run this command (Production Issues)
```
Give the example
...
$ make install
$ make test
$ make build
$ ./app
```

### Running the tests

Explain how to run the automated tests for this system
```sh
Give the example
...
$ make test
```

### Running the tests (With coverage appear on)

Explain how to run the automated tests for this system
```sh
Give the example
...
$ make cover
```

### Deployment

Add additional notes about how to deploy this on a live system

### Built With

* [Gorilla/Mux] The rest framework used
* [Mockery] Mock Up Generator
* [GoMod] - Dependency Management
* [Docker] - Container Management

