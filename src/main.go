package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	likeRepo "github.com/isd-sgcu/johnjud-backend/src/app/repository/like"
	petRepo "github.com/isd-sgcu/johnjud-backend/src/app/repository/pet"
	imageSrv "github.com/isd-sgcu/johnjud-backend/src/app/service/image"
	likeSrv "github.com/isd-sgcu/johnjud-backend/src/app/service/like"
	petSrv "github.com/isd-sgcu/johnjud-backend/src/app/service/pet"
	"github.com/isd-sgcu/johnjud-backend/src/config"
	"github.com/isd-sgcu/johnjud-backend/src/database"
	likePb "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
	petPb "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	imagePb "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		sig := <-s

		log.Info().
			Str("service", "graceful shutdown").
			Msgf("got signal \"%v\" shutting down service", sig)

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Error().
				Str("service", "graceful shutdown").
				Msgf("timeout %v ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("cleaning up: %v", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Error().
						Str("service", "graceful shutdown").
						Err(err).
						Msgf("%v: clean up failed: %v", innerKey, err.Error())
					return
				}

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("%v was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend").
			Msg("Failed to load config")
	}

	db, err := database.InitPostgresDatabase(&conf.Database, conf.App.IsDevelopment())
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend").
			Msg("Failed to init postgres connection")
	}

	fileConn, err := grpc.Dial(conf.Service.File, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "johnjud-file").
			Msg("Cannot connect to service")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", conf.App.Port))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend").
			Msg("Failed to start service")
	}

	grpcServer := grpc.NewServer()

	likeRepo := likeRepo.NewRepository(db)
	likeService := likeSrv.NewService(likeRepo)

	imageClient := imagePb.NewImageServiceClient(fileConn)
	imageService := imageSrv.NewService(imageClient)
	petRepo := petRepo.NewRepository(db)
	petService := petSrv.NewService(petRepo, imageService)

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	likePb.RegisterLikeServiceServer(grpcServer, likeService)
	petPb.RegisterPetServiceServer(grpcServer, petService)

	reflection.Register(grpcServer)
	go func() {
		log.Info().
			Str("service", "backend").
			Msgf("JohnJud backend starting at port %v", conf.App.Port)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().
				Err(err).
				Str("service", "backend").
				Msg("Failed to start service")
		}
	}()

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"server": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})

	<-wait

	grpcServer.GracefulStop()
	log.Info().
		Str("service", "backend").
		Msg("Closing the listener")
	lis.Close()
	log.Info().
		Str("service", "backend").
		Msg("End the program")
}
