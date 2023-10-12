package main

import (
	"fmt"
	"strings"

	"github.com/zhiyin2021/zysms/smpp"
)

func main() {
	msg := `id:R7e1qWDqT2 sub:001 dlvrd:001 submit date:2310121502 done date:2310121502 stat:DELIVRD err:000 text:000`
	report := &smpp.DeliverReport{}
	report.MsgId, msg = splitReport(msg, "id:")
	report.Sub, msg = splitReport(msg, "sub:")
	report.Dlvrd, msg = splitReport(msg, "dlvrd:")
	report.SubmitDate, msg = splitReport(msg, "submit date:")
	report.DoneDate, msg = splitReport(msg, "done date:")
	report.Stat, msg = splitReport(msg, "stat:")
	report.Err, msg = splitReport(msg, "err:")
	report.Text = strings.TrimSpace(strings.Replace(msg, "text:", "", 1))
	fmt.Println(report)
}

func splitReport(content, sub1 string) (retSub string, retContent string) {
	content = strings.TrimSpace(content)
	n := strings.Index(content, sub1)
	if n == -1 {
		return "", content
	}
	n += len(sub1)
	m := strings.Index(content[n:], " ")
	if m == -1 {
		return content, ""
	}
	return content[n : m+n], content[n+m:]
}
