package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/*
	负责二进制流的读和写
*/

type BinaryBuffer struct {
	buffer *bytes.Buffer
}

var order binary.ByteOrder

func init() {
	order = binary.LittleEndian
}

func NewBinaryBuffer() *BinaryBuffer {
	return &BinaryBuffer{buffer : new(bytes.Buffer)}
}
func (buffer *BinaryBuffer) WriteInt(v int32) {
	binary.Write(buffer.buffer, order, v)
}
func (buffer *BinaryBuffer) ReadInt() int32 {
	var v int32 = 0
	binary.Read(buffer.buffer, order, &v)
	return v
}
func (buffer *BinaryBuffer) WriteFloat(v float32) {
	binary.Write(buffer.buffer, order, v)
}
func (buffer *BinaryBuffer) ReadFloat() float32 {
	var v float32 = 0
	binary.Read(buffer.buffer, order, &v)
	return v
}
func (buffer *BinaryBuffer) WriteBytes(v []byte) {
	cnt := int32(len(v))
	buffer.WriteInt(cnt)
	binary.Write(buffer.buffer, order, v)
}
func (buffer *BinaryBuffer) ReadBytes()([]byte) {
	cnt := buffer.ReadInt()
	v := make([]byte, cnt)
	buffer.buffer.Read(v)
	return v
}
func (buffer *BinaryBuffer) WriteString(v string) {
	buffer.WriteBytes([]byte(v))
}
func (buffer *BinaryBuffer) ReadString() string {
	v := buffer.ReadBytes()
	return string(v)
}

func testBinaryBuffer() {
	buffer := NewBinaryBuffer()

	var a int32 = 30
	var b float32 = 4.2
	str := "123你好"

	buffer.WriteInt(a)
	buffer.WriteFloat(b)
	buffer.WriteString(str)

	a = 0
	b = 4.2
	str = ""
	a = buffer.ReadInt()
	b = buffer.ReadFloat()
	str = buffer.ReadString()

	fmt.Println(a, b, str)
}