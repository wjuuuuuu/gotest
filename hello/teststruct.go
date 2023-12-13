package main

import "fmt"

type config struct {
	Key string
	Val interface{}
}

func main() {

	var configs []config
	var low config
	var up config

	low.Key = "server_name"
	low.Val = "server01"
	up.Key = "common"
	up.Val = low

	configs = append(configs, up)

	var mini2 config
	var low2 config
	var up2 config

	mini2.Key = "file_path"
	mini2.Val = "/Users/codjs/dev/go/src/github.com/chaechaep/kt_dx_platform/bin/log"
	low2.Key = "log"
	low2.Val = mini2
	up2.Key = "common"
	up2.Val = low2

	configs = append(configs, up2)

	fmt.Println(configs)
}
