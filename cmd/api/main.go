package main

import "SchoolManagerApi/internal/server"

// @title		   Marks Api
// @version		 1.0
// @description This is an API for managing marks
// @termsOfService  http://swagger.io/terms/

// @contact.name   Jose Armando Coronel Vasquez
// @contact.email  joseacvz81@gmail.com
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description hi , how are you?
func main() {
	server.ListenServer()
}
