package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/Amniversary/wedding-logic-redpacket/config"
	"math/big"
	mathrand "math/rand"
	"strconv"
	"time"
)

type RedFlash struct {
	Money               float64   //总金额
	ResultRedPacketData []float64 //分裂出来各个红包的大小
	IndexMin            int64     //最小红包下标
	IndexMax            int64     //最大红包下标
}

/*
传入红包个数和红包总金额 红包类型
类型 1 3 为随机红包 num 为红包分裂数量  money 为总金额
类型 2   为普通红包 num 为红包数量     money 为单个红包金额  红包总金额= num* money
*/
func GenRedPacket(redPacketType int64, num, money float64) (RedFlash, error) {
	var resultPacketData []float64
	var indexMin int
	var min float64
	redFlash := RedFlash{IndexMax: 0, IndexMin: 0}
	//普通红包
	if redPacketType == config.GENERAL_RED_PACKET {
		redFlash.Money = num * money
		for i := 0; i < int(num); i++ {
			redFlash.ResultRedPacketData = append(redFlash.ResultRedPacketData, money)
		}
		return redFlash, nil
	}
	//其他红包
	redFlash.Money = money
	money = money * 100
	average := int64(money / num)
	getAve := float64(RandInt64(1, average))
	if getAve > config.RED_PACKET_MAX {
		getAve = float64(RandInt64(1, config.RED_PACKET_MAX))
	}
	subNum := getAve
	resultPacketData = append(resultPacketData, float64(subNum/100))

	for i := 2; i <= int(num); i++ {
		average = int64(money-subNum) / int64(int(num)-i+1)
		getAve = float64(RandInt64(1, average))
		if getAve > config.RED_PACKET_MAX {
			getAve = float64(RandInt64(1, config.RED_PACKET_MAX))
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
		if number > (config.RED_PACKET_MAX / 100) {
			average = RandInt64(1, int64(config.RED_PACKET_MAX-resultPacketData[indexMin]*100))
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
	redFlash.ResultRedPacketData = resultPacketData
	redFlash.IndexMin = int64(iMin)
	redFlash.IndexMax = int64(iMax)
	return redFlash, nil
}

//TODO: 随机取值 内部函数没做大小输入判断
func RandInt64(min, max int64) int64 {
	i, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return i.Int64() + min
}

//TODO: 数组乱序
func RandArray(arr []float64) []float64 {
	rands := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	l := len(arr)
	for i := l - 1; i > 0; i-- {
		r := rands.Intn(i)
		arr[r], arr[i] = arr[i], arr[r]
	}
	return arr
}
