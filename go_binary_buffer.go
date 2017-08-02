package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

/*
	负责二进制流的读和写
*/

const (
	errBinaryRead  = "err binary buffer read"
	errBinaryWrite = "err binary buffer write"
	errBinaryBytes = "err binary buffer Bytes" // 二进制数组最大50k
	errBinaryArray = "err binary array"        // 数组大小最大1k
	maxBytes       = 51200                     // 50 * 1024
	maxArray       = 1024
)

type BinaryBuffer struct {
	buffer *bytes.Buffer
}

var order binary.ByteOrder

func init() {
	order = binary.LittleEndian
}

func NewBinaryBuffer(data []byte) *BinaryBuffer {
	if data == nil {
		return &BinaryBuffer{buffer: new(bytes.Buffer)}
	}
	return &BinaryBuffer{buffer: bytes.NewBuffer(data)}
}
func (buffer *BinaryBuffer) GetBytes() []byte {
	return buffer.buffer.Bytes()
}
func (buffer *BinaryBuffer) WriteInt(v int32) {
	err := binary.Write(buffer.buffer, order, v)
	if err != nil {
		panic(errors.New(errBinaryWrite))
	}
}
func (buffer *BinaryBuffer) ReadInt() int32 {
	var v int32 = 0
	err := binary.Read(buffer.buffer, order, &v)
	if err != nil {
		panic(errors.New(errBinaryRead))
	}
	return v
}
func (buffer *BinaryBuffer) WriteFloat(v float32) {
	err := binary.Write(buffer.buffer, order, v)

	if err != nil {
		panic(errors.New(errBinaryRead))
	}
}
func (buffer *BinaryBuffer) ReadFloat() float32 {
	var v float32 = 0
	err := binary.Read(buffer.buffer, order, &v)

	if err != nil {
		panic(errors.New(errBinaryRead))
	}
	return v
}
func (buffer *BinaryBuffer) WriteBytes(v []byte) {
	cnt := int32(len(v))
	if cnt > maxBytes {
		panic(errors.New(errBinaryBytes))
	}

	buffer.WriteInt(cnt)
	err := binary.Write(buffer.buffer, order, v)
	if err != nil {
		panic(errors.New(errBinaryBytes))
	}
}
func (buffer *BinaryBuffer) ReadBytes() []byte {
	cnt := buffer.ReadInt()

	if cnt > maxBytes {
		panic(errors.New(errBinaryBytes))
	}

	v := make([]byte, cnt)
	n, err := buffer.buffer.Read(v)

	if err != nil || n != int(cnt) {
		panic(errors.New(errBinaryBytes))
	}

	return v
}

// 扩展部分
func (buffer *BinaryBuffer) WriteString(v string) {
	buffer.WriteBytes([]byte(v))
}
func (buffer *BinaryBuffer) ReadString() string {
	v := buffer.ReadBytes()
	return string(v)
}
func (buffer *BinaryBuffer) WriteIntArray(v []int32) {
	var size int32 = int32(len(v))
	if size > maxArray {
		panic(errors.New(errBinaryArray))
	}
	buffer.WriteInt(size)
	for i := int32(0); i < size; i++ {
		buffer.WriteInt(v[i])
	}
}
func (buffer *BinaryBuffer) ReadIntArray() []int32 {
	var size int32 = buffer.ReadInt()
	if size > maxArray {
		panic(errors.New(errBinaryArray))
	}
	ret := make([]int32, size)
	for i := int32(0); i < size; i++ {
		ret[i] = buffer.ReadInt()
	}
	return ret
}
func (buffer *BinaryBuffer) WriteFloatArray(v []float32) {
	var size int32 = int32(len(v))
	if size > maxArray {
		panic(errors.New(errBinaryArray))
	}

	buffer.WriteInt(size)
	for i := int32(0); i < size; i++ {
		buffer.WriteFloat(v[i])
	}
}
func (buffer *BinaryBuffer) ReadFloatArray() []float32 {
	var size int32 = buffer.ReadInt()

	if size > maxArray {
		panic(errors.New(errBinaryArray))
	}

	ret := make([]float32, size)
	for i := int32(0); i < size; i++ {
		ret[i] = buffer.ReadFloat()
	}
	return ret
}
func (buffer *BinaryBuffer) WriteStringArray(v []string) {
	var size int32 = int32(len(v))

	if size > maxArray {
		panic(errors.New(errBinaryArray))
	}

	buffer.WriteInt(size)
	for i := int32(0); i < size; i++ {
		buffer.WriteString(v[i])
	}
}
func (buffer *BinaryBuffer) ReadStringArray() []string {
	var size int32 = buffer.ReadInt()
	if size > maxArray {
		panic(errors.New(errBinaryArray))
	}

	ret := make([]string, size)
	for i := int32(0); i < size; i++ {
		ret[i] = buffer.ReadString()
	}
	return ret
}

func testBinaryBuffer() {
	buffer := NewBinaryBuffer(nil)

	var a int32 = 30
	var b float32 = 4.2
	str := "123你好"
	cc := []int32{1, 2, 3}
	dd := []float32{1.1, 2.2, 3.3}

	buffer.WriteInt(a)
	buffer.WriteFloat(b)
	buffer.WriteString(str)
	buffer.WriteIntArray(cc)
	buffer.WriteFloatArray(dd)

	a = 0
	b = 4.2
	str = ""
	a = buffer.ReadInt()
	b = buffer.ReadFloat()
	str = buffer.ReadString()
	cc = buffer.ReadIntArray()
	dd = buffer.ReadFloatArray()

	fmt.Println(a, b, str, cc, dd)
}
