package business

import (
	"math/big"
	"crypto/rand"
	"fmt"
	"strconv"
)

const RED_PACKET_MAX = 20000 //TODO: 单位 (分)

type RedFlash struct {
	ResultRedPacketData []float64
	IndexMin int64
	IndexMax int64
}

func GenRedPacket(num float64, money float64) (RedFlash, bool, string) {
	var resultPacketData []float64
	var indexMin int
	var min float64
	var err string
	RedFlash := RedFlash{IndexMax: 0, IndexMin: 0}
	if num > 200 {
		err = "一次最多可发200个红包"
		return RedFlash, false, err
	}
	if money < (num * 0.01) {
		err = "单个红包不能小于0.01元"
		return RedFlash, false, err
	}
	money = money * 100
	if money > (num * RED_PACKET_MAX) {
		err = "单个红包不能大于200元"
		return RedFlash, false, err
	}
	average := int64(money / num)
	getAve := float64(RandInt64(1, average))
	if getAve > RED_PACKET_MAX {
		getAve = float64(RandInt64(1, RED_PACKET_MAX))
	}
	subNum := getAve
	resultPacketData = append(resultPacketData, float64(subNum/100))

	for i := 2; i <= int(num); i++ {
		average = int64(money-subNum) / int64(int(num)-i+1)
		getAve = float64(RandInt64(1, average))
		if getAve > RED_PACKET_MAX {
			getAve = float64(RandInt64(1, RED_PACKET_MAX))
		}
		resultPacketData = append(resultPacketData, float64(getAve/100))
		subNum += getAve
	}
	left := money - subNum
	for {
		if left <= 0 {
			break
		}
		indexMin = 0
		min = resultPacketData[0]
		for i := 1; i < int(num); i++ {
			if resultPacketData[i] < min {
				indexMin = i
				min = resultPacketData[i]
			}
		}
		number, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", resultPacketData[indexMin]+left/100), 64)
		if number > (RED_PACKET_MAX / 100) {
			average = RandInt64(1, int64(RED_PACKET_MAX-resultPacketData[indexMin]*100))
			leftNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", resultPacketData[indexMin]+float64(average)/100), 64)
			resultPacketData[indexMin] = leftNum
			left -= float64(average)
		} else {
			leftNum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", resultPacketData[indexMin]+left/100), 64)
			resultPacketData[indexMin] = leftNum
			left -= left
		}
	}
	resultPacketData = RandArray(resultPacketData) //TODO: 红包乱序
	//TODO: 查找最大值 最小值
	iMax := 0
	iMin := 0
	max := resultPacketData[0]
	min = resultPacketData[0]
	for i := 1; i < int(num); i++ {
		if resultPacketData[i] > max {
			iMax = i
			max = resultPacketData[i]
		}
		if resultPacketData[i] < min {
			iMin = i
			min = resultPacketData[i]
		}
	}
	RedFlash.ResultRedPacketData = resultPacketData
	RedFlash.IndexMin = int64(iMin)
	RedFlash.IndexMax = int64(iMax)
	return RedFlash, true, ""
}

//TODO: 随机取值
func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}
