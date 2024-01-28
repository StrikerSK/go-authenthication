# User authentication using Go

This repository is looking for solutions for creating, configuring and resolving user sessions.

#Prerequisites:
* `config.yaml` file should be present as in provided example, currently it has to be in root of the project to detect configuration file properly, therefore please do not move this file 

## Usage:
From **root** directory run:
1. Run `docker-compose up` to prepare the infrastructure
2. Run `go run main.go`

Application currently configured to run on localhost port `8080`, to see example use this [collection](./postman/UserAuthRequests.json) in Postman.

## Resources: 
* [Using JWT for Authentication in a Golang Application](https://learn.vonage.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr/)
* [Session based authentication in Go](https://www.sohamkamani.com/golang/session-based-authentication/)
* [Working with Redis in Go](https://www.alexedwards.net/blog/working-with-redis)