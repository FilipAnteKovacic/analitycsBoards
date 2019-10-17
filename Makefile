all:
	protoc -I/usr/local/include -I.  --go_out=plugins=grpc:.  chart.proto
	protoc -I/usr/local/include -I.  --grpc-gateway_out=logtostderr=true:. chart.proto
	protoc -I/usr/local/include -I.  --swagger_out=logtostderr=true:.  chart.proto
	protoc-go-inject-tag -input=./chart.pb.go
	go generate .