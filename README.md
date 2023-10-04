# PayWise

```cmd
rm -f ./internal/transport/pb/*.go
```

To generate go code from protobuf
```cmd
 protoc --proto_path=./internal/transport/grpc/proto --go_out=./internal/transport/grpc/pb --go_opt=paths=source_relative --go-grpc_out=./internal/transport/grpc/pb --go-grpc_opt=paths=source_relative ./internal/transport/grpc/proto/*.proto
```

To start working with Evans to explore your grpc services
```cmd
➜ PayWise git:(main) ✗ docker run -it -v "${pwd}:/mount:ro" ghcr.io/ktr0731/evans:latest --path ./internal/transport/grpc/proto --proto service_paywise.proto --host localhost --port 9000 repl

  ______
 |  ____|
 | |__    __   __   __ _   _ __    ___
 |  __|   \ \ / /  / _. | | '_ \  / __|
 | |____   \ V /  | (_| | | | | | \__ \
 |______|   \_/    \__,_| |_| |_| |___/

 more expressive universal gRPC client


pb.Paywise@localhost:9000> show service
+---------+------------+---------------+----------------+
| SERVICE |    RPC     | REQUEST TYPE  | RESPONSE TYPE  |
+---------+------------+---------------+----------------+
| Paywise | SignupUser | SignupRequest | SignupResponse |
| Paywise | LoginUser  | LoginRequest  | LoginResponse  |
+---------+------------+---------------+----------------+

```