/*
用户只需要调用初始化Transfer(NewTransfer)，
然后调用Read方法读取请求数据，调用Write方法发送数据，
而不需要关注封装和解数据包的细节
数据包的格式（dataLen|dataID|data），用户可将data数据加密后在进行发送
*/
package cnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

var HeaderLength uint32 = 8
var MaxSize uint32 = 4092

type DataPackage struct {
	ID      uint32
	Len     uint32
	Content []byte
}

type Transfer struct {
	// 连接
	Conn net.Conn

	// 包头长度
	HeaderLength uint32

	// 数据包
	DataPackage

	// 数据包最大长度
	MaxSize uint32
}

// 初始化Transfer并使用默认参数
func NewTransfer(conn net.Conn) *Transfer {
	return &Transfer{
		Conn:         conn,
		HeaderLength: HeaderLength,
		MaxSize:      MaxSize,
	}
}

// 读数据ID和数据Content
func (t *Transfer) Read() (*DataPackage, error) {
	// 1. 读取包头
	buf := make([]byte, t.HeaderLength)
	_, err := io.ReadFull(t.Conn, buf)
	if err != nil {
		return nil, err
	}
	// 2. 解包头
	if err := t.unpack(buf); err != nil {
		return nil, err
	}
	// 3. 读取包数据
	if t.DataPackage.Len > 0 {
		t.DataPackage.Content = make([]byte, t.DataPackage.Len)
		_, err := io.ReadFull(t.Conn, t.DataPackage.Content)
		if err != nil {
			return nil, err
		}
	}
	return &t.DataPackage, nil
}

// 写数据（数据ID和数据Content,Len由Content计算所得）
func (t *Transfer) Write(dataPackage DataPackage) error {
	dataPackage.Len = uint32(len(dataPackage.Content))
	t.DataPackage = dataPackage
	buf, err := t.pack()
	if err != nil {
		return err
	}
	cnt, err := t.Conn.Write(buf)
	if err != nil {
		return errors.New(fmt.Sprintf("write data error: %s\n", err))
	}
	if cnt != len(buf) {
		return errors.New(fmt.Sprintf("write data length error: %s\n", err))
	}
	return nil
}

// 解包头获取数据的长度以进一步获取数据Content
func (t *Transfer) unpack(data []byte) error {
	// 创建一个二进制数据的io.Reader
	buf := bytes.NewReader(data)
	// 读取DataLength
	if err := binary.Read(buf, binary.BigEndian, &t.DataPackage.Len); err != nil {
		return errors.New(fmt.Sprintf("read DataLen error: %s\n", err))
	}
	// 这里进行判断包的长度是否超过设定的最大长度
	if t.DataPackage.Len > t.MaxSize {
		return errors.New("too large msg data recv")
	}
	// 读取DataID
	if err := binary.Read(buf, binary.BigEndian, &t.DataPackage.ID); err != nil {
		return errors.New(fmt.Sprintf("read DataID error: %s\n", err))
	}
	return nil
}

// 打包数据Len、ID和Content
func (t *Transfer) pack() ([]byte, error) {
	// 封装创建一个存放bytes字节的buf
	buf := bytes.NewBuffer([]byte{})
	// 写入dataLen
	if err := binary.Write(buf, binary.BigEndian, t.DataPackage.Len); err != nil {
		return nil, errors.New(fmt.Sprintf("write dataLen error: %s\n", err))
	}
	// 写入dataId
	if err := binary.Write(buf, binary.BigEndian, t.DataPackage.ID); err != nil {
		return nil, errors.New(fmt.Sprintf("write dataID error: %s\n", err))
	}
	// 写入data
	if err := binary.Write(buf, binary.BigEndian, t.DataPackage.Content); err != nil {
		return nil, errors.New(fmt.Sprintf("write data error: %s\n", err))
	}
	return buf.Bytes(), nil
}
