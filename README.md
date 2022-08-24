# Member Dashboard

## Introduction

Member Dashboard is the source of truth for who has access to what at the makerspace.  
Membership statuses will be pulled down from Payment Providers on a daily basis.  
If a member has made a payment in the past 30 days, they will be considered an active member.

## High level

- The server pulls payment information from paypal (and stores in the db) so we can tell who is currently an active member
- the server will maintain access lists and periodically push those access lists to the microcontrollers on the network
- The microcontroller (aka a resource) stores its access list locally so it's not dependant on the network when someone wants to access the space

## Documentation
Documentation can be found [here](https://hackrva.github.io/memberdashboard/development/setup.html)