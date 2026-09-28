package main

import (
	"jobsite/config"
	"jobsite/driver"
	"jobsite/repository"
	"jobsite/router"
	"jobsite/service"
	"jobsite/utility"
	"net/http"

	"github.com/go-kit/log/level"
)

func main() {

	logFile, logger, err := utility.InitLogger("logger", "system.log")
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize logger", "err", err)
		return
	}
	defer logFile.Close()

	cfg, err := config.LoadConfig()
	if err != nil {
		level.Error(logger).Log("msg", "failed to load config", "err", err)
		return
	}
	db, err := driver.ConnectDB(cfg)
	if err != nil {
		level.Error(logger).Log("msg", "failed to connect database", "err", err)
		return
	}

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	countryRepo := repository.NewCountryRepository(db)
	jobTypeRepo := repository.NewJobTypeRepository(db)

	userService := service.NewUserService(userRepo, logger, cfg.ApiSecret)
	jobSvc := service.NewJobService(jobRepo, db, logger)
	lookupSvc := service.NewLookupService(countryRepo, jobTypeRepo)

	handler := router.NewRouter(logger, userService, jobSvc, lookupSvc, cfg.ApiSecret)
	level.Info(logger).Log("msg", "starting HTTP server", "port", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		level.Error(logger).Log("msg", "server stopped with error", "err", err)
		return
	}
}
