package main

import (
	_ "grxx_query/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

