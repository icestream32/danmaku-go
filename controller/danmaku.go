package controller

import (
	"danmaku-go/services/bilibili/danmaku"

	"github.com/labstack/echo/v4"
)

func GetDanmaku(c echo.Context) error {

	danmaku.GetDanmaku("1")

	return c.String(200, "OK")
}
