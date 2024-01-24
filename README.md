# Proof of Work sample

This project demonstrate how to build and use proof of work in your projects in order to prevent a DDOS or fraud or
bruteforce your endpoints.

In general - your client should generale some puzzle that require too much resources on client side (to generate it),
but it is almost free to check the solution on server side.

In this case we use `hashcash`.

## How to run

```shell
git clone git@github.com:kazhuravlev/example-pow-guard.git pow-guard
cd pow-guard

# run the server
go run ./cmd/server

# connect to the server
telnet 127.0.0.1 8888

# generate a hashcash (it will work only once and can be used in next 30 sec)
go run ./cmd/gen-token
```

- You should put generated solution to telnet as first message.
- You can generate hashes as many as you want.
- After successful authorization - just hit the `Enter` to interact with server.
