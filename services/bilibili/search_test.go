package bilibili

import (
	"fmt"
	"testing"
)

func TestSearchAll(t *testing.T) {

	body, err := SearchAll("原神")
	if err != nil {

		t.Fail()
	}
	fmt.Println(body)
}

func TestSearchByType(t *testing.T) {

	body, err := SearchByType("修复4K放映厅", BiliUser)
	if err != nil {

		t.Fail()
	}
	fmt.Println(body)
}
