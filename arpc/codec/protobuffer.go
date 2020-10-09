package codec

import pb "github.com/gogo/protobuf/proto"

type ProtoBuffer struct{}

func (c *ProtoBuffer) Marshal(v interface{}) ([]byte, error) {
	msg, _ := v.(pb.Message)
	return pb.Marshal(msg)
}

func (c *ProtoBuffer) Unmarshal(data []byte, v interface{}) error {
	msg, _ := v.(pb.Message)
	return pb.Unmarshal(data, msg)
}
