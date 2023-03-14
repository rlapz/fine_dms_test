package delivery

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"enigmacamp.com/fine_dms/config"
	"enigmacamp.com/fine_dms/controller"
	"enigmacamp.com/fine_dms/manager"
	"github.com/gin-gonic/gin"
)

type AppServer struct {
	infra  manager.InfraManager
	ucMgr  manager.UsecaseManager
	srv    *http.Server
	engine *gin.Engine
	secret config.Secret
}

func NewAppServer() AppServer {
	cfg := config.NewAppConfig()

	infrMgr := manager.NewInfraManager(cfg)
	rpMgr := manager.NewRepoManager(infrMgr)
	ucMgr := manager.NewUsecaseManager(rpMgr)

	engine := gin.Default()
	return AppServer{
		infra: infrMgr,
		ucMgr: ucMgr,
		srv: &http.Server{
			Addr: fmt.Sprintf("%s:%s", cfg.ApiConfig.Host,
				cfg.ApiConfig.Port),
			Handler: engine,
		},
		engine: engine,
		secret: cfg.Secret,
	}
}

func (self *AppServer) Run() error {
	if err := self.infra.Init(); err != nil {
		return err
	}

	defer self.infra.Deinit()

	// router version 1
	self.v1()

	go func() {
		if err := self.srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalln(err)
			}
		}

		log.Println("Server closed")
	}()

	self.waitSignal()
	return nil
}

// private
func (self *AppServer) v1() {
	baseRg := self.engine.Group("/v1")
	controller.NewUserController(baseRg, self.ucMgr.UserUsecase(), &self.secret)
	controller.NewTagsController(baseRg, self.ucMgr.TagsUsecase())
	controller.NewFileController(baseRg, self.ucMgr.FileUseCase(), &self.secret)
}

func (self *AppServer) waitSignal() {
	// gracefully exit this program
	qChan := make(chan os.Signal, 1)

	signal.Notify(qChan, os.Interrupt, syscall.SIGTERM)
	<-qChan

	var timeout = 3 * time.Second

	fmt.Println()
	log.Println("Shutdown Server ...")
	log.Printf("Add timeout %d seconds\n", timeout/time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := self.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()

	log.Println("Server exiting")
}
