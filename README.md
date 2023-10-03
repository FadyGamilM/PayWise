# PayWise

```cmd
rm -f ./internal/transport/pb/*.go
```

```cmd
 protoc --proto_path=./internal/transport/grpc/proto --go_out=./internal/transport/grpc/pb --go_opt=paths=source_relative --go-grpc_out=./internal/transport/grpc/pb --go-grpc_opt=paths=source_relative ./internal/transport/grpc/proto/*.proto
```

How to list the files 
```cmd
➜ PayWise git:(main) ✗ docker run --rm -it -v "${pwd}:/mount:ro" alpine ls /mount/internal/transport/grpc/proto/
Unable to find image 'alpine:latest' locally
latest: Pulling from library/alpine
96526aa774ef: Pull complete
Digest: sha256:eece025e432126ce23f223450a0326fbebde39cdf496a85d8c016293fc851978
Status: Downloaded newer image for alpine:latest
rpc_login_user.proto   rpc_signup_user.proto  service_paywise.proto  user.proto
```