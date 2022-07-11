
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
```
ROLLING_PERIOD is rolling period for deposit wallet in second unit (180 = 2 minutes)
THRESHOLD is deposit threshold within rolling period

- Then run this command (Development Issues)
```
Give the example
...
$ make install
$ make run-dev
```

### Running the tests

Explain how to run the automated tests for this system
```sh
Give the example
...
$ make test-dev
```

### Running the tests (With coverage appear on)

Explain how to run the automated tests for this system
```sh
Give the example
...
$ make cover
```

