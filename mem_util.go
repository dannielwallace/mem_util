// mem_util project mem_util.go
package mem_util

import (
	"errors"
	"fmt"
	"runtime"
)

//
// 内存池需要实现的接口
//
type MemPool interface {
	Alloc(size int) []byte
}

//
type RawMemPool struct {
	maxPackSize int
}

//
// 简单的内存池实现，用于避免频繁的零散内存申请
//
type SimpleMemPool struct {
	memPool     []byte
	memPoolSize int
	maxPackSize int
}

func NewRawMemPool(maxPackSize int) (*RawMemPool, error) {
	if maxPackSize <= 0 {
		return nil, errors.New("maxPackSize <= 0")
	}

	return &RawMemPool{
		maxPackSize: maxPackSize,
	}, nil
}

func (this *RawMemPool) Alloc(size int) (result []byte) {
	result = make([]byte, size)
	return
}

//
// 创建一个简单内存池，预先申请'memPoolSize'大小的内存，每次分配内存时从中切割出来，直到剩余空间不够分配，再重新申请一块。
// 参数'maxPackSize'用于限制外部申请内存允许的最大长度，所以这个值必须大于等于'memPoolSize'。
//
func NewSimpleMemPool(memPoolSize, maxPackSize int) (*SimpleMemPool, error) {
	if maxPackSize > memPoolSize {
		return nil, errors.New("maxPackSize > memPoolSize")
	}

	return &SimpleMemPool{
		memPool:     make([]byte, memPoolSize),
		memPoolSize: memPoolSize,
		maxPackSize: maxPackSize,
	}, nil
}

//
// 申请一块内存，如果'size'超过'maxPackSize'设置将返回nil
//
func (this *SimpleMemPool) Alloc(size int) (result []byte) {
	if size > this.maxPackSize {
		return nil
	}

	if len(this.memPool) < size {
		this.memPool = make([]byte, this.memPoolSize)
	}

	result = this.memPool[0:size]
	this.memPool = this.memPool[size:]
	return
}

// func Allocated returns a string of current memory usage such as "8KB" or "16MB"
func Allocated() string {
	s := new(runtime.MemStats)
	runtime.ReadMemStats(s)

	MemAllocated := s.Alloc

	var arrayMemUnit [4]uint64
	for i := 0; i < 4; i++ {
		arrayMemUnit[i] = MemAllocated % 1024
		MemAllocated = MemAllocated / 1024
	}

	return fmt.Sprintf("%d GB, %d MB, %d KB, %d Byte", arrayMemUnit[3], arrayMemUnit[2], arrayMemUnit[1], arrayMemUnit[0])
}
