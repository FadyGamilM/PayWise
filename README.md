# PayWise

```cmd
rm -f ./internal/transport/pb/*.go
```

```cmd
 protoc --proto_path=./internal/transport/grpc/proto --go_out=./internal/transport/grpc/pb --go_opt=paths=source_relative --go-grpc_out=./internal/transport/grpc/pb --go-grpc_opt=paths=source_relative ./internal/transport/grpc/proto/*.proto
```