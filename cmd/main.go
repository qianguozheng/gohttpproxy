package main

import "log"
import "github.com/qianguozheng/gohttpproxy"

func main() {
	log.Println("Gohttpproxy start")

	p := gohttpproxy.New()
	p.Start(":8888")
	log.Println("Exit")
}
