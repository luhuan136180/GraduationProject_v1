package utils

import (
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"time"
)

var sf *sonyflake.Sonyflake

const max = 1<<16 - 1

func init() {
	// 有一定的几率会产生相同的machineID
	// 比如两个引用这个包的实例，在同一个微秒被创建...
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	st := sonyflake.Settings{
		StartTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (uint16, error) {
			return uint16(r.Intn(max)), nil
		},
		CheckMachineID: nil,
	}

	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		zap.L().Panic("sonyflake init fails")
	}
}

func NextID() string {
	ui64id, _ := sf.NextID()
	return strconv.FormatUint(ui64id, 10)
}
