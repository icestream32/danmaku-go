package bilibili

import (
	"fmt"
	"testing"
)

func TestGenWebTicket(t *testing.T) {

	fmt.Println(GetCookie())
	GenWebTicket()
}
