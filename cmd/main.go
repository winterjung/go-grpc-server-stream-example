package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	pb "example.com/image"
)

func main() {
	ctx := context.Background()
	grpcServer := NewGRPCServer()
	httpServer, err := NewHTTPServer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatal(err)
		}

		log.WithField("port", "8080").Info("starting gRPC server")
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			log.Fatal(err)
		}
	}()

	go func() {
		log.WithField("port", "8081").Info("starting HTTP server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit

	log.Info("stopping HTTP server")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Info("stopping gRPC server")
	grpcServer.GracefulStop()
}

type imageServer struct{}

var _ pb.ImageServer = (*imageServer)(nil)

func NewGRPCServer() *grpc.Server {
	imgServer := &imageServer{}
	grpcServer := grpc.NewServer()
	pb.RegisterImageServer(grpcServer, imgServer)
	return grpcServer
}

func NewHTTPServer(ctx context.Context) (*http.Server, error) {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		),
	)
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := pb.RegisterImageHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:8080",
		options,
	); err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    ":8081",
		Handler: cors(mux),
	}

	return server, nil
}

func (s *imageServer) GetImage(req *pb.GetImageRequest, stream pb.Image_GetImageServer) error {
	f, err := os.Open("images/" + req.Path)
	if err != nil {
		return status.Error(codes.NotFound, "file not found")
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.WithError(err).Error("failed to close file")
		}
	}()

	// Maximum 16KB size per stream.
	buf := make([]byte, 16*2<<10)
	for {
		num, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := stream.Send(&pb.GetImageResponse{Data: buf[:num]}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
