package sgip

import (
	"strconv"
	"time"
)

type SgipHeader struct {
	SeqId [3]uint32 // 源节点编号 + 月日时分秒 + 流水序号
}

func getTm() uint32 {
	tm, _ := strconv.ParseUint(time.Now().Format("0215040506"), 10, 10)
	return uint32(tm)
}
