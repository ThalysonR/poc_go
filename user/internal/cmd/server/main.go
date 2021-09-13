package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thalysonr/poc_go/common/log"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"github.com/thalysonr/poc_go/user/internal/datasources"
	"github.com/thalysonr/poc_go/user/internal/transport"
	tgrpc "github.com/thalysonr/poc_go/user/internal/transport/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()

	logger := NewZapLogger(log.DEBUG)
	log.SetLogger(logger)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic("failed to connect db")
	}

	logger.Info(ctx, "Starting services...")
	services := getServices(db)
	errChan := runServices(services)

	select {
	case <-done:
		logger.Info(ctx, "Termination signal received")
	case err := <-errChan:
		logger.Info(ctx, "Server error: %s", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()
	for _, service := range services {
		service.Stop()
	}

	logger.Info(ctx, "Server Stopped")
}

////////////////////////////////////////////////////////////////////////////////
///////                       AUXILIARY FUNCTIONS                        ///////
////////////////////////////////////////////////////////////////////////////////

func getServices(db *gorm.DB) []transport.Service {
	userRepository := datasources.NewUserDBDatasource(db)
	userService := service.NewUserService(userRepository)

	return []transport.Service{
		transport.NewHttpServer(*userService),
		tgrpc.NewGrpcServer(*userService),
	}
}

func runServices(services []transport.Service) <-chan error {
	errChan := make(chan error, 1)

	for _, service := range services {
		go func(errC chan<- error, svc transport.Service) {
			err := svc.Start()
			if err != nil {
				errC <- err
			}
		}(errChan, service)
	}
	return errChan
}
