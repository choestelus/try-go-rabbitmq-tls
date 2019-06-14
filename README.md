# try-go-rabbitmq-tls

---

very quick and dirty example of configuring rabbitmq and go client to connect via TLS connection

for demonstration purpose only

### Certificates generation

Create Cerfificate Authority first, then create server and client certificates by signing with created CA certificate, example certificates in `cert` directory

### RabbitMQ configuration

see `rabbitmq.config` and `rabbitmq-env.conf` examples in this repository.
note that `rabbitmq.config` uses [Erlang term configuration format](http://erlang.org/doc/man/config.html)

### running main executable

This project uses go module as dependency manager and go version 1.12

export environment variables first, example in `.env`
then run with `go run main.go`
