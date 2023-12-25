package main

import (
	"context"
	"database/sql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/okoroemeka/simple_bank/api"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	_ "github.com/okoroemeka/simple_bank/doc/statik"
	"github.com/okoroemeka/simple_bank/gapi"
	"github.com/okoroemeka/simple_bank/pb"
	"github.com/okoroemeka/simple_bank/util"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config variables", err)
		return
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create grpc server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener", err)
	}

	log.Printf("gRPC server is running on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("Cannot start grpc server", err)
	}

}
func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create grpc server", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("Cannot create gateway server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New() // statikFS implements http.FileSystem

	if err != nil {
		log.Fatal("Cannot create statik file system:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot create http listener", err)
	}

	log.Printf("gateway server is running on %s", listener.Addr().String())
	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("Cannot start http gateway server", err)
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot server instance", err)
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
