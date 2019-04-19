# License - Manager
The following is a microservice which features a REST controller
A connection to a sqlite database

The REST endpoint is served on http://localhost:8080/

* POST /{user}
* POST /{user}/licenses

The endpoints will attempt to autenticate the user from the database as the request comes through. 

With the correct password the user will get the licenese they have attached to thier account.

## Building

`make static`

### Building a docker container
`make DOCKERTAG=pair-man:latest docker`

## Running

To run locally you can run with
`go run *.go`

### Running with a docker container

`docker run -it sevren/license-manager:latest`