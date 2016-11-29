package message

import (
	"bytes"
	"encoding/binary"
)

type MessageHead struct{
	id		uint16
	len		int32
	compress	int8
}

func NewHead(a_id uint16, a_len int32) *MessageHead {
	return &MessageHead{
		id: a_id,
		len: a_len,
	}
}

func (head MessageHead) GetMsgId() uint16 {
	return head.id
}

func (head MessageHead) GetDataLen() int32 {
	return head.len
}

func (head *MessageHead) Encode() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, head.id)
	binary.Write(buf, binary.LittleEndian, head.len)
	binary.Write(buf, binary.LittleEndian, head.compress)
	return buf.Bytes()
}

func (head *MessageHead) Decode(data []byte)  {
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &head.id)
	binary.Read(buf, binary.LittleEndian, &head.len)
	binary.Read(buf, binary.LittleEndian, &head.compress)
}
