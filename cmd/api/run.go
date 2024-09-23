package main

import (
	"context"
	"go-web-api-starter/internal/apiutils"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func run(
	ctx context.Context,
	getEnv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer cancel()

	//Uncomment for postgres config
	//sslEnabled := common.BoolEnv(getEnv, "SSL_ENABLED", false)
	//dbConfig := database.NewDatabase(sslEnabled)
	//db, err := dbConfig.OpenDB("postgres")
	//if db != nil {
	//	defer func() {
	//		dbErr := db.Close()
	//		if dbErr != nil {
	//			stderr.Write([]byte(dbErr.Error()))
	//		}
	//	}()
	//}
	//if err != nil {
	//	return err
	//}
	//
	//err = db.RunMigrations(database.DialectPostgres)
	//if err != nil {
	//	return err
	//}

	app := &application{
		config: apiutils.NewApiConfig(getEnv, "API_PORT"),
	}

	httpServer := newServer(
		app.config.Logger,
	)

	err := apiutils.Serve(
		httpServer,
		app.config.Logger,
		app.config.Env,
		app.config.Port,
		app.config.Version,
		ctx,
	)
	if err != nil {
		return err
	}
	return nil
}
