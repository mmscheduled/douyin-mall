package main

import (
    "github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
    h := server.Default()
    h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
        ctx.JSON(200, map[string]string{"message": "pong"})
    })
    h.Spin()
}