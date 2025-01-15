package smpp

import (
	"encoding/hex"
	"testing"

	"github.com/zhiyin2021/zysms/codec"
)

func TestDeliverSMReport(t *testing.T) {
	data := `id:2d4d4563-1cea-4646-a2c5-852dadc36c381997 sub:001 dlvrd:001 submit date:2407301056 done date:2407301056 stat:DELIVRD err:000 text:[kvgame]Mã OTP 6792.`
	var report DeliverReport
	if err := report.Unmarshal(data); err != nil {
		t.Fatal(err)
	}
	t.Log(report)

	data = `000000dc00000005000000000000000100010132333438303637383537303336000500536f6b616265740004000000000000f1007269b20e36330c1683c162adf2ad0c6bd1c4e2700b772bd35ae5f02d2726cb70b431d90c20f3ba580730180c2064b65d4ed60130180c20f3bab89da683c8617a5907321a0c27bbc56cb01a20e4b7bb0c2287e9651d321a0c27bbc56c301b20737a98ae03c422336995121b2065b95c0730180c001e002536663030313030312d653765302d346262612d383765342d65613739626432383463646600042700010204230003030000`
	buf, _ := hex.DecodeString(data)

	deliver := NewDeliverSM()
	reader := codec.NewReader(buf)
	deliver.Unmarshal(reader)

	t.Log(deliver)
}
