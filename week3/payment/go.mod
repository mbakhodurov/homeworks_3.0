module github.com/mbakhodurov/homeworks2/week3/payment

replace github.com/mbakhodurov/homeworks2/week3/shared => ../shared

go 1.25.4

require (
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/mbakhodurov/homeworks2/week3/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.82.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260709200747-435963d16310.1 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260630182238-925bb5da69e7 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
