package controller

import (
	"danmu-go/services"

	"github.com/labstack/echo/v4"
)

func GetDanmaku(c echo.Context) error {

	services.GetDanmaku()

	return c.String(200, "OK")
}
