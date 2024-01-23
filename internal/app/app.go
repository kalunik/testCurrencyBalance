package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kalunik/testCurrencyBalance/config"
	"github.com/kalunik/testCurrencyBalance/internal/api"
	repo "github.com/kalunik/testCurrencyBalance/internal/repository"
	"github.com/kalunik/testCurrencyBalance/internal/usecase"
	"github.com/kalunik/testCurrencyBalance/pkg/log"
	"github.com/nats-io/stan.go"
	"net/http"
	"os"
	"os/signal"
)

type App struct {
	r      *api.Router
	pgpool *pgxpool.Pool
	connW  stan.Conn
	connR  stan.Conn
	log    log.Logger
	conf   config.AppConfig
}

func NewApp(pgpool *pgxpool.Pool, connW stan.Conn, connR stan.Conn, logger log.Logger, config config.AppConfig) *App {
	return &App{r: nil, pgpool: pgpool, connW: connW, connR: connR, log: logger, conf: config}
}

func (a *App) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	pgRepo := repo.NewPgRepository(a.pgpool)
	natsRepo := repo.NewNatsRepository(a.connW, a.connR, a.conf)

	urlService := usecase.NewWalletUsecase(pgRepo, natsRepo)

	walletHandlers := api.NewWalletHandlers(urlService, a.log)

	a.r = api.NewRouter()

	a.r.PathMetaRoutes(walletHandlers)

	go func() {
		for {
			urlService.ActivateTransaction()
		}
	}()

	a.log.Infof("api server will start on %s port", a.conf.Server.Port)
	go http.ListenAndServe(a.conf.Server.Port, a.r.Mux)

	<-ctx.Done()
	a.log.Infof("programm stopped")
	cancel()
}
