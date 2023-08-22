package main

import (
	"container/heap"
)

// HuffmanNode 表示Huffman树的节点
type HuffmanNode struct {
	char  rune
	freq  int
	left  *HuffmanNode
	right *HuffmanNode
}

// HuffmanHeap 表示Huffman树的最小堆
type HuffmanHeap []*HuffmanNode

// 实现heap.Interface接口的Len方法
func (h HuffmanHeap) Len() int {
	return len(h)
}

// 实现heap.Interface接口的Less方法
func (h HuffmanHeap) Less(i, j int) bool {
	return h[i].freq < h[j].freq
}

// 实现heap.Interface接口的Swap方法
func (h HuffmanHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// 实现heap.Interface接口的Push方法
func (h *HuffmanHeap) Push(x interface{}) {
	*h = append(*h, x.(*HuffmanNode))
}

// 实现heap.Interface接口的Pop方法
func (h *HuffmanHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// 构建Huffman树
func buildHuffmanTree(freq map[rune]int) *HuffmanNode {
	h := &HuffmanHeap{}
	for char, f := range freq {
		heap.Push(h, &HuffmanNode{char: char, freq: f})
	}

	for h.Len() > 1 {
		node1 := heap.Pop(h).(*HuffmanNode)
		node2 := heap.Pop(h).(*HuffmanNode)
		parent := &HuffmanNode{
			freq:  node1.freq + node2.freq,
			left:  node1,
			right: node2,
		}
		heap.Push(h, parent)
	}

	return heap.Pop(h).(*HuffmanNode)
}

// 生成Huffman编码
func generateHuffmanCode(root *HuffmanNode, code string, codes map[rune]string) {
	if root == nil {
		return
	}

	if root.left == nil && root.right == nil {
		codes[root.char] = code
		return
	}

	generateHuffmanCode(root.left, code+"0", codes)
	generateHuffmanCode(root.right, code+"1", codes)
}
