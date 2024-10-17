package main

import (
	"danmaku-go/router"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	fmt.Println("Hello, World!")

	e := echo.New()

	router.Route(e)

	e.GET("/", func(c echo.Context) error {

		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
