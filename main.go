package main

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/okoroemeka/simple_bank/api"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	_ "github.com/okoroemeka/simple_bank/doc/statik"
	"github.com/okoroemeka/simple_bank/gapi"
	"github.com/okoroemeka/simple_bank/mail"
	"github.com/okoroemeka/simple_bank/pb"
	"github.com/okoroemeka/simple_bank/util"
	"github.com/okoroemeka/simple_bank/worker"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
)

func main() {

	config, err := util.LoadConfig(".")

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Could not load config variables")
		return
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to database")
	}

	runDBMigration(config.MigrationUrl, config.DBSource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runGatewayServer(config, store, taskDistributor)
	go runTaskProcessor(config, redisOpt, store)
	runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationUrl, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create migration instance:")
	}

	if err := migration.Up(); !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msgf("Cannot migrate database:%s", err.Error())
	}

	log.Info().Msg("db migration completed successfully")

}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create grpc server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create listener")
	}

	log.Info().Msgf("gRPC server is running on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start grpc server")
	}

}
func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	processor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("Starting task processor")

	if err := processor.Start(); err != nil {
		log.Fatal().Err(err).Msg("Cannot start task processor")
	}

}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create grpc server")
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
		log.Fatal().Err(err).Msg("Cannot create gateway server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New() // statikFS implements http.FileSystem

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create statik file system:")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create http listener")
	}

	log.Info().Msgf("gateway server is running on %s", listener.Addr().String())

	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start http gateway server")
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot server instance")
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server")
	}
}
