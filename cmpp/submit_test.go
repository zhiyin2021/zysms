package cmpp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zhiyin2021/zysms/codec"
)

func BenchmarkSubmit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			p := NewSubmitReq(V30).(*SubmitReq)
			p.PkTotal = 1
			p.PkNumber = 1
			p.RegisteredDelivery = 1
			p.MsgLevel = 1
			p.ServiceId = "test"
			p.FeeUserType = 2
			p.FeeTerminalId = "13500002696"
			// FeeTerminalType:    0
			p.MsgFmt = 8
			p.MsgSrc = "900001"
			p.FeeType = "02"
			p.FeeCode = "10"
			p.ValidTime = "151105131555101+"
			p.AtTime = ""
			p.SrcId = "900001"
			p.DestUsrTl = 1
			p.TpUdhi = 1
			p.DestTerminalId = []string{"+8613500002696"}
			p.Message.SetMessage("你的验证码为:283919,如非本人操作,请忽略.【百度网盘】", codec.UCS2)

			w := codec.NewWriter()
			p.Marshal(w)

			reader := codec.NewReader(w.Bytes())
			p1 := NewSubmitReq(V30).(*SubmitReq)
			p1.Unmarshal(reader)
		}()
	}
	b.StopTimer()
}

func BenchmarkDeliver(b *testing.B) {
	tm := time.Now().Format("0601021504")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			deliver := NewDeliverReq(V30).(*DeliverReq)
			deliver.DestId = "900001"
			deliver.SrcTerminalId = "8613500002696"
			deliver.RegisterDelivery = 1
			deliver.MsgId = 8613500002696
			deliver.MsgFmt = 0
			deliver.Report = &DeliverReport{
				MsgId:        deliver.MsgId,
				SubmitTime:   tm,
				DoneTime:     tm,
				Stat:         "DELIVED",
				SmscSequence: 0,
			}
			w := codec.NewWriter()
			deliver.Marshal(w)

			reader := codec.NewReader(w.Bytes())

			p1 := NewDeliverReq(V30).(*DeliverReq)
			p1.Unmarshal(reader)
		}()
	}
	b.StopTimer()
}

// func BenchmarkSubmitResp(b *testing.B) {
// 	data := `000001010000000400000001000000000000000001010101746573740000000000000231333530303030323639360000000000000000000000000000000000000000000000010839303030303130323130000000003135313130353133313535353130312B000000000000000000000000000000000000393030303031000000000000000000000000000000012B38363133353030303032363936000000000000000000000000000000000000003E4F6076849A8C8BC178014E3A003A003200380033003900310039002C5982975E672C4EBA64CD4F5C002C8BF75FFD7565002E3010767E5EA67F5176D830110000000000000000000000000000000000000000`
// 	buf, _ := hex.DecodeString(data)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		func() {
// 			reader := codec.ReaderPool.Get().(*codec.BytesReader)
// 			reader.Init(buf)
// 			defer codec.ReaderPool.Put(reader)
// 			p := NewSubmitReq(V30).(*SubmitReq)
// 			p.Unmarshal(reader)
// 		}()
// 	}
// 	b.StopTimer()
// }

func TestSubmitSm(t *testing.T) {

	p := NewSubmitReq(V30).(*SubmitReq)
	p.PkTotal = 1
	p.PkNumber = 1
	p.RegisteredDelivery = 1
	p.MsgLevel = 1
	p.ServiceId = "test"
	p.FeeUserType = 2
	p.FeeTerminalId = "13500002696"
	// FeeTerminalType:    0
	p.MsgFmt = 8
	p.MsgSrc = "900001"
	p.FeeType = "02"
	p.FeeCode = "10"
	p.ValidTime = "151105131555101+"
	p.AtTime = ""
	p.SrcId = "900001"
	p.DestUsrTl = 1
	p.TpUdhi = 1
	p.DestTerminalId = []string{"+8613500002696"}
	p.Message.SetMessage("你的验证码为:283919,如非本人操作,请忽略.【百度网盘】", codec.UCS2)
	w := codec.NewWriter()
	p.Marshal(w)
	t.Logf("%X", w.Bytes())
}

func TestDeliverSm(t *testing.T) {

	tm := time.Now().Format("0601021504")
	deliver := NewDeliverReq(V30).(*DeliverReq)
	deliver.DestId = "900001"
	deliver.SrcTerminalId = "8613500002696"
	deliver.RegisterDelivery = 1
	deliver.MsgId = 8613500002696
	deliver.MsgFmt = 0
	deliver.Report = &DeliverReport{
		MsgId:        deliver.MsgId,
		SubmitTime:   tm,
		DoneTime:     tm,
		Stat:         "DELIVED",
		SmscSequence: 0,
	}
	w := codec.NewWriter()

	deliver.Marshal(w)

	t.Logf("%X", w.Bytes())
	reader := codec.NewReader(w.Bytes())
	p1 := NewDeliverReq(V30).(*DeliverReq)
	err := p1.Unmarshal(reader)
	require.Nil(t, err)
	t.Logf("%v => %+v", err, p1)
}
