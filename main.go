package main

import "fmt"
import "./server"
func main()  {
	fmt.Println("Airdisk-CMS Server Start")
	server.Run()
	fmt.Println("Airdisk-CMS Server Exit")
}