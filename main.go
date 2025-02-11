package main

import (
	"context"
	"errors"
	db "main/db/sqlc"
	"main/gapi"
	"main/pb"
	"main/token"
	"main/utils"
	"main/ws"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGINT,
	syscall.SIGTERM,
}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	waitGroup, waitGroupContext := errgroup.WithContext(ctx)

	store := db.NewStore(connPool)

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create token maker")
	}

	// Initialize WebSocket manager
	wsManager := ws.NewManager()
	go wsManager.Start()

	runGPRCServer(waitGroupContext, waitGroup, config, store, tokenMaker)
	runWebSocketServer(waitGroupContext, waitGroup, config, wsManager, tokenMaker)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error in wait group")
	}
}

func runGPRCServer(ctx context.Context, waitGroup *errgroup.Group, config utils.Config, store db.Store, tokenMaker token.Maker) {
	server, err := gapi.NewServer(config, store, tokenMaker)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterTicTacToeServer(grpcServer, server)
	reflection.Register(grpcServer)

	listner, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Str("port", listner.Addr().String()).Msg("started gRPC server")

		err = grpcServer.Serve(listner)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully stopping gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gracefully stopped gRPC server")
		return nil
	})
}

func runWebSocketServer(ctx context.Context, waitGroup *errgroup.Group, config utils.Config, wsManager *ws.Manager, tokenMaker token.Maker) {
	// Create WebSocket auth middleware
	wsHandler := ws.NewHandler(wsManager, tokenMaker)

	// Create a new HTTP server for WebSocket
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler.HandleConnection)

	wsServer := &http.Server{
		Addr:    config.WebSocketServerAddress,
		Handler: mux,
	}

	waitGroup.Go(func() error {
		log.Info().Str("address", wsServer.Addr).Msg("started WebSocket server")

		if err := wsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("WebSocket server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully stopping WebSocket server")

		if err := wsServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown WebSocket server gracefully")
			return err
		}

		log.Info().Msg("gracefully stopped WebSocket server")
		return nil
	})
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration")
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("cannot run migration")
	}

	log.Info().Msg("migration completed successfully")
}
