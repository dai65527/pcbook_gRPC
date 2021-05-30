package serializer

import (
	"math/rand"
	"pcbook/pb"
	"pcbook/sample"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(laptop1, binaryFile)
	if err != nil {
		t.Fatal(err)
	}

	laptop2 := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(binaryFile, laptop2)
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(laptop1, laptop2) {
		t.Fatal("wrote file and read file not match")
	}

	err = WriteProtobufToJSONFile(laptop1, jsonFile)
	if err != nil {
		t.Fatal(err)
	}
}
