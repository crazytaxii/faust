package main

import (
	"fmt"
	"testing"
)

func TestRandNum(t *testing.T) {
	fmt.Println("rand num:", randNum(10000000, 100000000))
}

func TestImgConvert(t *testing.T) {
	srcPath := "./test/Go-Logo_Aqua.png"
	buf, err := imgConvert(srcPath, JPG)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println("img size:", len(buf))
}
