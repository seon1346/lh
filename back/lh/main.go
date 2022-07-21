package main

import (
	"fmt"
	"lh/fabric"
	"lh/http"
	"os"
)

func main() {
	os.Setenv("TESTTEST", "xx")

	fmt.Println(os.Getenv("TESTTEST"))
	//DISCOVERY_AS_LOCALHOST=TRUE  필수[env]
	os.Setenv("DISCOVERY_AS_LOCALHOST", "TRUE")
	fabric.Init_Wallet()
	fabric.ConTest()
	http.Run_HttpServer("", "")

}
