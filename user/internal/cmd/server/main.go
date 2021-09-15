package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	cconfig "github.com/thalysonr/poc_go/common/config"
	"github.com/thalysonr/poc_go/common/log"
	"github.com/thalysonr/poc_go/user/internal/app/repository"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"github.com/thalysonr/poc_go/user/internal/config"
	"github.com/thalysonr/poc_go/user/internal/datasources"
	"github.com/thalysonr/poc_go/user/internal/transport"
	tgrpc "github.com/thalysonr/poc_go/user/internal/transport/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	config             *config.Config
	configSubscription *cconfig.ConfigSubscription
	ctx                context.Context
	db                 *gorm.DB
	errChan            chan error
	logger             *ZapLogger
	repositories       repository.Repositories
	services           sync.Map
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
	if s.db != nil {
		if sqlDB, _ := s.db.DB(); sqlDB != nil {
			sqlDB.Close()
		}
	}
	s.services.Range(func(key, value interface{}) bool {
		value.(transport.Service).Stop()
		return true
	})
	s.logger.Info(s.ctx, "Server Stopped")
}

func (s *Server) run() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s.services.Range(func(key, value interface{}) bool {
		go runService(*s.config, s.errChan, value.(transport.Service))
		return true
	})

out:
	for {
		select {
		case <-done:
			s.logger.Info(s.ctx, "Termination signal received")
			break out
		case err := <-s.errChan:
			if err != nil {
				s.logger.Error(s.ctx, "Server error: %s", err)
				break out
			}
		}
	}
}

func runService(cfg config.Config, errC chan<- error, svc transport.Service) {
	err := svc.Start(cfg)
	if err != nil {
		errC <- err
	}
}

func (s *Server) setup() {
	s.errChan = make(chan error)
	s.ctx = context.Background()
	s.setupLogger()
	s.setupConfig()
	s.setupDB(s.logger)
	s.setupRepositories(s.db)
	s.setupServices()
}

func (s *Server) setupConfig() {
	co, err := cconfig.GetConfigObservable()
	if err != nil {
		panic("could not get config")
	}

	cfg := config.Config{}
	s.config = &cfg
	subscription, err := co.Subscribe(&cfg, func(err error) {
		s.logger.Info(s.ctx, "config changed")
		if err != nil {
			log.GetLogger().Error(s.ctx, "could not refresh config")
		} else {
			s.services.Range(func(key, value interface{}) bool {
				svc := value.(transport.Service)
				if svc.ConfigChanged(*s.config) {
					svc.Stop()
					go runService(*s.config, s.errChan, value.(transport.Service))
				}
				return true
			})
		}
	})

	if err != nil {
		panic("could not get config")
	}
	s.configSubscription = subscription
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

	s.services.Store(uuid.New(), transport.NewHttpServer(*userService))
	s.services.Store(uuid.New(), tgrpc.NewGrpcServer(*userService))
}
