# Back-end To Do List App

To do list app is a simple app that allows user to write down their daily tasks.

## Features
- Register
- Email verification
- Login
- Forget password
- Logout
- Update user
- Read user
- Create a to-do list by date
- Delete to-do list
- Update to-do list
- Read to-do list by date
- Read all to-do list

## Tech

Dillinger uses a number of open source projects to work properly:

- Language: Golang
- Database: MongoDB
- Packages: HTTP Router, Go-Playground Validator, Golang-JWT, Google UUID, GoMail V2, Bcrypt, MongoDB

## Installation

This To Do List App is built with Go 1.19 version

Install the packages and start the server.

```sh
go get github.com/julienschmidt/httprouter
go get github.com/go-playground/validator
go get github.com/golang-jwt/jwt
go get github.com/google/uuid
go get gopkg.in/gomail.v2
go get golang.org/x/crypto/bcrypt
go get go.mongodb.org/mongo-driver/bson
go get go.mongodb.org/mongo-driver/mongo

cd cmd
go run main.go
```

For development environments...

```sh
SECRET_KEY=yoursecretkey
CONFIG_AUTH_EMAIL=youremail@gmail.com
CONFIG_AUTH_PASSWORD=youremailpassword
```

## API Documentation

Full API Documentation in Postman

[![N|Solid](https://logosdownload.com/logo/postman-logo-512.png)](https://documenter.getpostman.com/view/13066205/2s84DrP1cA)

Check out my journey creating this project in Medium!

[![N|Solid](https://miro.medium.com/max/8978/1*s986xIGqhfsN8U--09_AdA.png)](https://medium.com/@pratamafarhan10/back-end-to-do-list-app-planning-my-first-project-with-go-1-be36647df691)
