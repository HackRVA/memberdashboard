# Member Dashboard

<div style="display: inline-flex; gap: 8px">
    <a href="https://pkg.go.dev/github.com/HackRVA/memberserver"><img src="https://pkg.go.dev/badge/github.com/HackRVA/memberserver.svg" alt="Go Reference"></a>
    <img alt="Test passing" src="https://github.com/Hackrva/memberdashboard/workflows/Test/badge.svg" />
    <img alt="Test UI passing" src="https://github.com/Hackrva/memberdashboard/workflows/Test%20UI/badge.svg" />
    <a href="https://goreportcard.com/report/github.com/HackRVA/memberserver">
    <img alt="Go report" src="https://goreportcard.com/badge/github.com/HackRVA/memberserver">
    </a>
</div>


## Introduction

The Member Dashboard allows us to register members for use with HackRVA's RFID system.  It will track whether they are active by pulling their subscription status from our payment provider.


## High level Orchestration

- The server pulls payment information from paypal (and stores in the db) so we can tell who is currently an active member
- The server will maintain access lists and periodically push those access lists to the RFID Readers on the network
- The RFID reader stores its access list locally so it's not dependant on the network when someone wants to access the space

## Run locally

```bash
# install frontend deps
make deps-frontend
make
```

> note that this doesn't setup the mqtt server and it doesn't communicate with any mqtt devices

## Documentation

Documentation can be found [here](https://hackrva.github.io/memberdashboard/development/setup.html)
