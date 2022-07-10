# go-codebase

This Project is used to create Golang Codebase to be embbed to MyPertamina

### Prerequisites

What things you need to install the software and how to install them

```
Give examples
Golang v1.15
Go Mod
....
```

### Installing

A step by step series of examples that tell you have to get a development env running

Say what the step will be
- Create ENV file (.env) with this configuration:
```
APP_NAME=go-codebase
PORT=9091
MARIADB_HOST=localhost
MARIADB_PORT=3306
MARIADB_USERNAME=username
MARIADB_PASSWORD=password
MARIADB_DATABASE=database
MARIADB_MAX_OPEN_CONNECTIONS=25
MARIADB_MAX_IDLE_CONNECTIONS=25
KAFKA_BROKERS=localhost
KAFKA_SSL_ENABLE=true
KAFKA_USERNAME=username
KAFKA_PASSWORD=password
KAFKA_SSL_ENABLE=false
BASIC_AUTH_USERNAME=username
BASIC_AUTH_PASSWORD=password
AES_SECRET=yoursecretkey
AES_IV=yoursalt
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

### Authors

* **TelkomDev** - *Initial work* - [Gitlab](https://gitlab.playcourt.id/telkomdev)


### License

This project is licensed under Telkom Indonesia License - see the [LICENSE.md](LICENSE.md) file for details

### Acknowledgments

* For sample file README.md, see [WIKI](https://gitlab.playcourt.id/telkomdev/codebase-backend/wikis/Readme.md-Sample) in this repository.