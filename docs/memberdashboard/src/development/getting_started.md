# Getting Started
## HighLevel
For Dev purposes, we use a fake "InMemory" db.
We can seed fake content with generators (see `./test/generators`) that we have built.
The frontend gets embeded into the binary. So, we don't need docker to host any database.

* Install [golang](https://go.dev/doc/install)
* Install [Node](https://nodejs.org/en/)
* build and run with `make`
* Install dependencies (see: [Install Dependencies](#install-dependencies))

e.g.
```bash
make build-ui
make run
```

Hopefully that makes the app fairly easy to stand up and get to hacking on.


## Install Dependencies
Download go modules
> note: go modules might download automatically the first time you run the app
```bash
go get ./...
```

Download node_modules
```bash
cd web
npm ci
```
