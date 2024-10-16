package danmaku

import (
	"fmt"
	"testing"
)

func TestGetDanmaku(t *testing.T) {

	data, err := GetDanmaku("26128418182")
	if err != nil {

		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(data)
}
