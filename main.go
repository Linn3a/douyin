package main

import (
	"douyin/bootstrap"
	"log"
)

func main() {

	// go service.RunMessageServer()

	//app.Use(logger.New())
	//app.Use(logger.New(logger.Config{
	//	Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	//	Output: os.Stdout,
	//}))

	app, err := bootstrap.Init()
	if err != nil {
		log.Println(err)
		return
	}

	err = app.Listen(":8080")
	if err != nil {
		return
	}
}
