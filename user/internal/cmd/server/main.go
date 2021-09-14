package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thalysonr/poc_go/common/log"
	"github.com/thalysonr/poc_go/user/internal/app/repository"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"github.com/thalysonr/poc_go/user/internal/datasources"
	"github.com/thalysonr/poc_go/user/internal/transport"
	tgrpc "github.com/thalysonr/poc_go/user/internal/transport/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	ctx          context.Context
	db           *gorm.DB
	logger       *ZapLogger
	repositories repository.Repositories
	services     []transport.Service
}

func main() {
	server := Server{}
	defer server.close()
	server.setup()
	server.run()
}

////////////////////////////////////////////////////////////////////////////////
///////                       AUXILIARY FUNCTIONS                        ///////
////////////////////////////////////////////////////////////////////////////////

func (s *Server) close() {
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
	for _, service := range s.services {
		service.Stop()
	}
	s.logger.Info(s.ctx, "Server Stopped")
}

func (s *Server) run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	for _, service := range s.services {
		go func(errC chan<- error, svc transport.Service) {
			err := svc.Start()
			if err != nil {
				errC <- err
			}
		}(errChan, service)
	}

	select {
	case <-done:
		s.logger.Info(s.ctx, "Termination signal received")
	case err := <-errChan:
		s.logger.Error(s.ctx, "Server error: %s", err)
	}
}

func (s *Server) setup() {
	s.ctx = context.Background()
	s.setupLogger()
	s.setupDB(s.logger)
	s.setupRepositories(s.db)
	s.setupServices()
}

func (s *Server) setupDB(logger *ZapLogger) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic("failed to connect db")
	}
	s.db = db
}

func (s *Server) setupLogger() {
	logger := NewZapLogger(log.DEBUG, nil)
	log.SetLogger(logger)
	s.logger = logger
}

func (s *Server) setupRepositories(db *gorm.DB) {
	s.repositories = repository.Repositories{
		User: datasources.NewUserDBDatasource(db),
	}
}

func (s *Server) setupServices() {
	userService := service.NewUserService(s.repositories.User)

	s.services = []transport.Service{
		transport.NewHttpServer(*userService),
		tgrpc.NewGrpcServer(*userService),
	}
}
