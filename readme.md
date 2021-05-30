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
