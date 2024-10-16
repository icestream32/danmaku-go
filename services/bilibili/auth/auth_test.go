package auth

import (
	"fmt"
	"testing"
)

func TestGenWebTicket(t *testing.T) {

	fmt.Println(GetCookie())
	fmt.Println(GenWbiKeysFromTicket())
	fmt.Println(GenWbiKeysFromNav())
	fmt.Println(GenWbi("https://search.bilibili.com/all?keyword=%E5%8E%9F%E7%A5%9E"))
}
