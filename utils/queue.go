package utils

import (
	"sync"
)

// 队列 切片实现
type Queue struct {
	data  []int64
	max   int
	size  int
	mutex sync.Mutex
}

// 创建
func NewQueue(max int) *Queue {
	return &Queue{max: max, size: 0, data: make([]int64, max)}
}

// 入队
func (s *Queue) Push(val int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.size == s.max {
		s.data = append(s.data[1:], val)
	} else {
		s.data = append(s.data, val)
		s.size++
	}
}

// 出队
func (s *Queue) Pop() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.size--
	val := s.data[0]
	s.data = s.data[1:]
	return val
}

// 显示队列元素
func (s *Queue) Data() []int64 {
	return s.data
}
