package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/minhajuddinkhan/huffman"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/smpp"
)

func main() {

	var x = []byte{0x00, 0x00, 0x00, 0x1a, 0x80, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x67, 0x00, 0x00, 0x03, 0x0f, 0x36, 0x32, 0x37, 0x34, 0x31, 0x36, 0x33, 0x33, 0x36, 0x00}
	r := bytes.NewReader(x)
	pdu, err := smpp.Parse(r)
	// resp := &smpp.SubmitSMResp{}
	// resp.Unmarshal(codec.NewReader(x))
	log.Println(pdu, err)
}

func main1() {
	//data := "您的验证码为:%d, 请在5分钟内使用,请勿泄露给他人.【创瑞短信】"

	text := "hello world"
	freq := make(map[rune]int)
	for _, char := range text {
		freq[char]++
	}

	root := buildHuffmanTree(freq)

	codes := make(map[rune]string)
	generateHuffmanCode(root, "", codes)

	fmt.Println("Huffman Codes:", codes)
	for char, code := range codes {
		fmt.Printf("%c: %s\n", char, code)
	}
	return
	tree := huffman.NewHuffmanTree("test")

	var encoded string
	err := tree.Encode(&encoded)
	if err != nil {
		panic(err)
	}

	fmt.Println(encoded) //00001000111111010110010011110111110010110000100110111000101

	decoded, err := tree.Decode(encoded)
	if err != nil {
		panic(err)
	}
	fmt.Println(decoded) //Are you a gopher?

	return
	tm := time.Now()
	for i := 0; i < 1000000; i++ {
		// cmpp.CreatePDUFromCmdID(cmpp.CMPP_ACTIVE_TEST, cmpp.V30)
	}
	s1 := time.Since(tm)
	header := cmpp.Header{CommandID: cmpp.CMPP_ACTIVE_TEST}
	tm = time.Now()
	for i := 0; i < 1000000; i++ {
		cmpp.CreatePDUHeader(header, cmpp.V30)
	}
	s2 := time.Since(tm)
	log.Println("time1:", s1.Microseconds(), s1.Milliseconds(), "time2:", s2.Microseconds(), s2.Milliseconds())
}
