# SchoolManagerApi

Simple API for School management using Golang

# How to Run

If you don't have Golang download from  [Go](https://go.dev)

If you have Golang , download the project in your PC:

~~~
git clone https://github.com/sheymor21/EducationSystemApi.git
~~~

if you use [Docker](https://www.docker.com), you can run docker compose up for run the project, you need to set the .env file before.
## Parameters

- DB_U (Required)  is the db user
- DB_P (Required)  is the db password
- DB_URI (Optional)  it the mongo db Uri, by default it is "mongodb://localhost:27017"
- DB_NAME (Optional) is the name of the database, by default it is "EducationSystem"
- port (Optional) is the application port , by default 8080
- env (Optional) is the application environment, by default dev

For run the project use:

~~~
go run cmd/api/main.go -DB_U (your db user) -DB_P (your db password)
~~~

If you prefer create a **.env** file, you already have a **.env.example** at the project with that only need to run
~~~
go run cmd/api/main.go
~~~

You can access the API UI by opening the root URL (`/`) of the running server.

## Login

### Create a User

You need to create a **teacher** for use all other endpoints, then you go to **login** , insert the credencial and you will receive a bearer token

When you create a teacher or a student, automatically you create a user, the users credentials are the carnet and the password.

### How Credentials Works

The **password** is created by the last 3 characters of the carnet and the lastname, if your name is **John Smith** and your carnet is **F421A34D21D** then you password will be 21D-smith


# Author

- [sheymor21](https://github.com/sheymor21)