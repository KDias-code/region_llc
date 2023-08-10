package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"syscall"
	"time"

	"go.uber.org/zap"

	"product-service/internal/config"
	"product-service/internal/handler"
	"product-service/internal/repository"
	"product-service/internal/service/catalogue"
	"product-service/pkg/log"
	"product-service/pkg/server"
)

const (
	schema      = "product"
	version     = "1.0.0"
	description = "product-service"
)

// Run инициализирует все приложение.
func Run() {
	logger := log.New(version, description)

	configs, err := config.New()
	if err != nil {
		logger.Error("ERR_INIT_CONFIG", zap.Error(err))
		return
	}

	repositories, err := repository.New(
		repository.WithPostgresStore(schema, configs.POSTGRES.DSN))
	//repository.WithMemoryStore())
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORY", zap.Error(err))
		return
	}
	defer repositories.Close()

	catalogueService, err := catalogue.New(
		catalogue.WithTasksRepository(repositories.Tasks),
		catalogue.WithTasksCache(repositories.Tasks),
	)
	if err != nil {
		logger.Error("ERR_INIT_CATALOGUE_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:          configs,
			CatalogueService: catalogueService,
		},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLER", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.HTTP.Port))
	if err != nil {
		logger.Error("ERR_INIT_SERVER", zap.Error(err))
		return
	}

	// Запускаем наш сервер в горутине, чтобы он не блокировался.
	if err = servers.Run(logger); err != nil {
		logger.Error("ERR_RUN_SERVER", zap.Error(err))
		return
	}

	// Мягкое завершение работы
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1) // create channel to signify a signal being sent

	// Мы допустим корректное завершение работы при выходе с помощью SIGINT (Ctrl+C)
	//SIGKILL, SIGQUIT или SIGTERM (Ctrl+/) не будут перехвачены.

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) //Когда посылается сигнал прерывания или завершения, уведомляем канал
	<-quit                                             // Это блокирует основной поток до тех пор, пока не будет получено прерывание.
	fmt.Println("Gracefully shutting down...")

	// установить крайний срок ожидания.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Не блокирует, если нет подключений, в противном случае будет ждать
	// до истечения срока дедлайна.
	if err = servers.Stop(ctx); err != nil {
		panic(err) // failure/timeout изящно отключает httpServer
	}

	fmt.Println("Running cleanup tasks...")
	// Сюда попадают ваши задачи по очистке

	fmt.Println("Server was successful shutdown.")
}
