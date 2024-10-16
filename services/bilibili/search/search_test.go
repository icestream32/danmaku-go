package search

import (
	"fmt"
	"testing"
)

func TestGetPlayerPageList(t *testing.T) {

	body, err := GetPlayerPageList("113244261843988")
	if err != nil {

		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(body)
}

func TestSearchArchivesByKeywords(t *testing.T) {

	body, err := SearchArchivesByKeywords("修复4K放映厅", "败犬女主太多了")
	if err != nil {

		fmt.Println(err)
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

func TestSearchAll(t *testing.T) {

	body, err := SearchAll("原神")
	if err != nil {

		t.Fail()
	}
	fmt.Println(body)
}
