package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJson converts protocol buffer to JSON string
func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		Indent:          "\t",
		UseProtoNames:   false,
	}
	data, err := marshaler.Marshal(message)
	return string(data), err
}
