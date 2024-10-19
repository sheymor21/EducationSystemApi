# EducationSystemApi

Simple API for Booking management using .Net core and Docker

# Dependencies

- Go 1.22.4
- MongoDriver 1.16.0

# How to Run

If you don't have Golang download from  [Go](https://go.dev)

If you have Golang , download the project in your PC:

~~~
git clone https://github.com/sheymor21/EducationSystemApi.git
~~~

## Parameters

- DB_U (Required)  is the db user
- DB_P (Required)  is the db password
- DB_URI (Optional)  it the mongo db Uri, by default it is "mongodb://localhost:27017"
- DB_NAME (Optional) is the name of the database, by default it is "EducationSystem"

For run the project use:

~~~
go run cmd/api/main.go -DB_U (your db user) -DB_P (your db password)
~~~

You can access the API UI by opening the root URL (`/`) of the running server.

# Author

- [sheymor21](https://github.com/sheymor21)