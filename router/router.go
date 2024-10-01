package router

import (
	"danmu-go/controller"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo) error {

	e.GET("/danmaku", controller.GetDanmaku)

	return nil
}
