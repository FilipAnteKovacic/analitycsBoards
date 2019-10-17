package main

import (
	context "context"
	fmt "fmt"
	"log"
	"net"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	status "google.golang.org/grpc/status"
)

type server struct {
}

func (*server) CreateSource(ctx context.Context, req *CreateSourceRequest) (*CreateSourceResponse, error) {

	// Request data to item data
	data := req.GetData()

	// Create data in DB
	err := createSource(*data)
	if err != nil {

		// Return response error
		return &CreateSourceResponse{
			Status: "error",
		}, err

	}

	dataUpdateRestart()

	// Push msg to websocket clients
	sendWSMessage(WSMessage{
		Sender: "server",
		Action: "create",
		Type:   "chart",
	})

	// Return response success
	return &CreateSourceResponse{
		Status: "success",
	}, nil

}

func (*server) ReadSource(ctx context.Context, req *ReadSourceRequest) (*ReadSourceResponse, error) {

	// Read by ID
	data, err := readByUniqueSource(req.GetUnique())
	if err != nil {

		// Return response error
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find  with specified ID: %v", err),
		)
	}

	// Return response data
	return &ReadSourceResponse{
		Data: &data,
	}, nil
}

func (*server) UpdateSource(ctx context.Context, req *UpdateSourceRequest) (*UpdateSourceResponse, error) {

	// Get request data
	reqData := req.GetData()

	// Read request ID
	Unique := reqData.GetUnique()

	// Update in DB
	err := updateByUniqueSource(Unique, *reqData)
	if err != nil {

		// Return response error
		return &UpdateSourceResponse{
			Status: "error",
		}, err

	}

	dataUpdateRestart()

	// Push msg to websocket clients
	sendWSMessage(WSMessage{
		Sender: "server",
		Action: "update",
		Type:   "chart",
		Source: *reqData,
	})

	// Return response success
	return &UpdateSourceResponse{
		Status: "success",
	}, nil

}

func (*server) DeleteSource(ctx context.Context, req *DeleteSourceRequest) (*DeleteSourceResponse, error) {

	// Read request Unique
	Unique := req.GetUnique()

	// Delete
	err := deleteByUniqueSource(Unique)
	if err != nil {

		// Return response error
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)

	}

	dataUpdateRestart()

	// Push msg to websocket clients
	sendWSMessage(WSMessage{
		Sender: "server",
		Action: "update",
		Type:   "chart",
		Source: Source{
			Unique: Unique,
		},
	})
	// Return response success
	return &DeleteSourceResponse{
		Status: "success",
	}, nil
}

func (*server) ListSource(ctx context.Context, req *ListSourceRequest) (*ListSourceResponse, error) {

	listRes, err := listSource()
	if err != nil {

		// Return response error
		return nil, status.Errorf(
			codes.DataLoss,
			err.Error(),
		)
	}

	// Return response success
	return &ListSourceResponse{
		Data: &Sources{
			Source: listRes,
		},
	}, nil
}

func (*server) CreateChart(ctx context.Context, req *CreateChartRequest) (*CreateChartResponse, error) {

	// Request data to item data
	data := req.GetData()

	// Create data in DB
	unique, err := createChart(*data)
	if err != nil {

		// Return response error
		return &CreateChartResponse{
			Status: "error",
		}, err

	}

	dataUpdateRestart()

	updateChartDataByUnique(unique)

	// Push msg to websocket clients
	sendWSMessage(WSMessage{
		Sender: "server",
		Action: "create",
		Type:   "chart",
	})

	// Return response success
	return &CreateChartResponse{
		Status: "success",
	}, nil

}

func (*server) ReadChart(ctx context.Context, req *ReadChartRequest) (*ReadChartResponse, error) {

	// Read by ID
	data, err := readByUniqueChart(req.GetUnique())
	if err != nil {

		// Return response error
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find  with specified ID: %v", err),
		)
	}

	// Return response data
	return &ReadChartResponse{
		Data: &data,
	}, nil
}

func (*server) UpdateChart(ctx context.Context, req *UpdateChartRequest) (*UpdateChartResponse, error) {

	// Get request data
	reqData := req.GetData()

	// Read request ID
	Unique := reqData.GetUnique()

	// Update in DB
	err := updateByUniqueChart(Unique, *reqData)
	if err != nil {

		// Return response error
		return &UpdateChartResponse{
			Status: "error",
		}, err

	}

	dataUpdateRestart()

	updateChartDataByUnique(Unique)

	// Push msg to websocket clients
	sendWSMessage(WSMessage{
		Sender: "server",
		Action: "update",
		Type:   "chart",
		Chart:  *reqData,
	})

	// Return response success
	return &UpdateChartResponse{
		Status: "success",
	}, nil

}

func (*server) DeleteChart(ctx context.Context, req *DeleteChartRequest) (*DeleteChartResponse, error) {

	// Read request Unique
	Unique := req.GetUnique()

	// Delete
	err := deleteByUniqueChart(Unique)
	if err != nil {

		// Return response error
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)

	}

	dataUpdateRestart()

	// Push msg to websocket clients
	sendWSMessage(WSMessage{
		Sender: "server",
		Action: "update",
		Type:   "chart",
		Chart: Chart{
			Unique: Unique,
		},
	})
	// Return response success
	return &DeleteChartResponse{
		Status: "success",
	}, nil
}

func (*server) ListChart(ctx context.Context, req *ListChartRequest) (*ListChartResponse, error) {

	listRes, err := listChart()
	if err != nil {

		// Return response error
		return nil, status.Errorf(
			codes.DataLoss,
			err.Error(),
		)
	}

	// Return response success
	return &ListChartResponse{
		Data: &Charts{
			Charts: listRes,
		},
	}, nil
}

func startGRPCServer(address, certFile, keyFile string) error {
	// create a listener on TCP port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("could not load TLS keys: %s", err)
	}

	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	// create a gRPC server object
	grpcServer := grpc.NewServer(opts...)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	RegisterDataServiceServer(grpcServer, &server{})

	// start the server
	log.Printf("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}

	return nil
}
