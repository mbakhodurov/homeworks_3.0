module github.com/mbakhodurov/homeworks2/week2/payment

replace github.com/mbakhodurov/homeworks2/week2/shared => ../shared

go 1.25.4

require (
	github.com/google/uuid v1.6.0
	github.com/mbakhodurov/homeworks2/week2/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.82.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260415201107-50325440f8f2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260630182238-925bb5da69e7 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
