package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenUploadKey(ext string) string {
	// 8 digital decimal
	return fmt.Sprintf("%s/%s.%s", getDate(), randNum(8), ext)
}

func getDate() string {
	return time.Now().Format("06-01-02")
}

func randNum(length int) string {
	rand.Seed(time.Now().Unix())
	min, max := 1, 10
	for i := 1; i < length; i++ {
		min, max = min*10, max*10
	}
	max -= 1
	return strconv.Itoa(rand.Intn(max-min) + min)
}
