package main

import "github.com/derekdowling/gardenswap-api/api"

func main() {
	server := api.BuildServer()
	server.Run(":3000")
}
