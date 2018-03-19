package business

import (
	"time"
	"math/rand"
)

//TODO: 数组乱序
func RandArray(arr []float64) []float64 {
	rands := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(arr)
	for i := l - 1; i > 0; i-- {
		r := rands.Intn(i)
		arr[r], arr[i] = arr[i], arr[r]
	}
	return arr
}
