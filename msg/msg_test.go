package msg

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestMsg(t *testing.T) {
	msg := NewMsg()
	msg.MagicNum = 'P'
	msg.Payload = []byte("hello world")
	msg.Meta = MsgMeta{
		Service:   "snow",
		Method:    "test",
		Id:        1,
		Timestamp: time.Now().Unix(),
		Type:      TypeRequest,
	}

	msgStr, err := json.Marshal(msg)
	if err != nil {
		t.Errorf("marshal msg failed:%s", err.Error())
		return
	}

	t.Logf("msg:%s", string(msgStr))

	msgContent, err := msg.Encode()
	if err != nil {
		t.Errorf("msg encode failed:%s", err.Error())
		return
	}

	newMsg := NewMsg()
	buffer := bytes.NewBuffer(msgContent)
	err = newMsg.Decode(buffer)
	if err != nil {
		t.Errorf("decode failed:%s", err.Error())
		return
	}

	newMsgStr, err := json.Marshal(newMsg)
	if err != nil {
		t.Errorf("marshal msg failed:%s", err.Error())
		return
	}
	t.Logf("new msg:%s, payload:%s", string(newMsgStr),
		string(newMsg.Payload))
}

func BenchmarkEncode(b *testing.B) {
	msg := NewMsg()
	msg.MagicNum = 'P'
	msg.Payload = []byte("hello world")
	msg.Meta = MsgMeta{
		Service:   "snow",
		Method:    "test",
		Id:        1,
		Timestamp: time.Now().Unix(),
		Type:      TypeRequest,
	}

	for i := 0; i < b.N; i++ {
		msg.Encode()
	}
}

func BenchmarkDecode(b *testing.B) {
	msg := NewMsg()
	msg.MagicNum = 'P'
	msg.Payload = []byte("hello world")
	msg.Meta = MsgMeta{
		Service:   "snow",
		Method:    "test",
		Id:        1,
		Timestamp: time.Now().Unix(),
		Type:      TypeRequest,
	}

	content, _ := msg.Encode()
	buffer := bytes.NewBuffer(content)
	newMsg := NewMsg()
	for i := 0; i < b.N; i++ {
		newMsg.Decode(buffer)
	}
}
