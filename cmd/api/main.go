package main

import "SchoolManagerApi/internal/server"

// @title		   SchoolManager Api
// @version		 1.0
// @description This is an API for managing marks
// @termsOfService  http://swagger.io/terms/

// @contact.name   Jose Armando Coronel Vasquez
// @contact.email  joseacvz81@gmail.com
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	server.ListenServer()
}
