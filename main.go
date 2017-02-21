package main

import "./server"
import (
	"fmt"
)


func main()  {
	fmt.Println("Airdisk-CMS Server Start")
	server.Run()
	fmt.Println("Airdisk-CMS Server Exit")
}