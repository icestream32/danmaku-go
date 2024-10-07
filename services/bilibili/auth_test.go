package bilibili

import (
	"fmt"
	"testing"
)

func TestGenWebTicket(t *testing.T) {

	fmt.Println(GetCookie())
	fmt.Println(GenWbiKeysFromTicket())
	fmt.Println(GenWbiKeysFromNav())
	fmt.Println(GenWbi("https://api.bilibili.com/x/space/wbi/acc/info?mid=1850091"))
}
