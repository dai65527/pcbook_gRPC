## これやってます
https://youtu.be/YzypniHHU3w?list=PLy_6D98if3UJd5hxWNfAqKMr15HZqFnqf

## 引っかかったところ
### #6 protoファイルからのコードの生成
protocによってコードを生成する際に`protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb`というコマンドを使っていたが、手元の環境（libprotoc 3.15.5）ではエラーとなる。

以下のように修正するとうまくいった
- `processor_message.proto`にオプション`option go_package = "./;pb";`を追記。
- コマンドを`protoc --proto_path=proto proto/*.proto --go_out=pb`とする。
  - `protoc --proto_path=proto proto/*.proto --go-grpc_out=. --go_out=pb`とでも良さそうだったがここでは`--go-grpc_out=.`オプションはいらないよう。（生成ファイルにdiffなし）

https://qiita.com/kitauji/items/bab05cc8215abe8a6431


### #8 sample/generator.go func NewLaptopのtimestamp生成

```
    // 非推奨
		// UpdatedAt:   ptypes.TimestampNow(),
    // 推奨
		UpdatedAt: timestamppb.Now(),
```

### #9 serializer/json.go func ProtobufToJSON内でのシリアライズ化に使用する関数
動画では、marshalerとして、protojson.MarshalOptionsを使うとエラーが出る（`cannot use message (type protoreflect.ProtoMessage) as type protoiface.MessageV1 in argument to marshaler.MarshalToString`)。おそらく、protobufの生成したファイルのバージョン違いによるもの。
こちらのパッケージを使えばOK。`google.golang.org/protobuf/encoding/protojson`。ただし、オプション名と関数名が異なるので注意。以下のようになる。

```
func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		Indent:          "\t",
		UseProtoNames:   true,
	}
	data, err := marshaler.Marshal(message)
	return string(data), err
}
```

### #10 protocによるservice用コード生成
これも動画と最新バージョンでは異なる部分があった。

ここまで使ってきたprotocのコマンド`protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=.`だけでは足りず、`LaptopServiceClient/Server`が定義される`laptop_service_grpc.pb.go`が生成されない。
オプション`--go-grpc_out=pb`を追加すればOK（`pb`は出力先フォルダ）

まとめると、こうなる（Makefile更新した）
```
protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
```

ちなみに、
- helloworld/helloworld.pb.go: メッセージやシリアライズ
- helloworld/helloworld_grpc.pb.go: gRPCのサーバ/クライアント

です。

### #10 pb.RegisterLaptopServiceServer
クライアントからテストをしようとする部分、LaptopServiceServerにLaptopServerを登録する以下のコードでエラーがでた。

```go
func startTestLaptopServer(t *testing.T) (*LaptopServer, string) {
	laptopServer := NewLaptopServer(NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer) // ここでエラー
	（略）
}
```

エラー：`cannot use laptopServer (variable of type *LaptopServer) as pb.LaptopServiceServer value in argument to pb.RegisterLaptopServiceServer: missing method mustEmbedUnimplementedLaptopServiceServer`

これもバージョン違いで、`pb.LaptopServiceServer`インターフェイスに `mustEmbedUnimplementedLaptopServiceServer`メソッドが追加されているから。（future compatibilityのためらしい）

対処法は2つ。

#### 1. `pb.UnimplementedLaptopServiceServer`を追加する（推奨）
`pb.UnimplementedLaptopServiceServer`に`pb.mustEmbedUnimplementedLaptopServiceServer`実装されているので、これをEmbedすればOKでした。

```
type LaptopServer struct {
	Store LaptopStore
	pb.UnimplementedLaptopServiceServer
}
```

#### 2. protocのオプションに`--go-grpc_opt=require_unimplemented_servers=false`を追加する（非推奨）
protocでコンパイルする際に`--go-grpc_out=require_unimplemented_servers=false`をつける。

```
protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb \
       --go-grpc_out=require_unimplemented_servers=false
```

後方互換性のために用意されているオプションなので新規作成の際は使うべきではない。

参考：https://github.com/grpc/grpc-go/blob/master/cmd/protoc-gen-go-grpc/README.md
参考：https://note.com/dd_techblog/n/nb8b925d21118
