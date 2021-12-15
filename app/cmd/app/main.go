package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"pehks1980/be2_hw81/internal/app/endpoint"
	"pehks1980/be2_hw81/internal/pkg/repository"
	"syscall"
	"time"
)



func main() {
	log.Print("Starting the app")
	// настройка порта, настроек хранилища, таймаут при закрытии сервиса
	port := flag.String("port", "8000", "Port")
	shutdownTimeout := flag.Int64("shutdown_timeout", 3, "shutdown timeout")
	//storageName := flag.String("storage name", "postgres://postuser:postpassword@192.168.1.204:5432/app1",
	//	"pg: 'postgres://dbuser:dbpasswd@ip_address:port/dbname'")
	storageName := flag.String("storage name", "http://192.168.1.210:9200",
		"es: 'http://IP:PORT'")
	storageName1 := flag.String("elsatic index name", "be2_hw8",
		"es: 'mysuperelasticsearchindex'")

	var repoif endpoint.RepoIf
	// подстановка в интерфейс соотвествующего хранилища
	//repoif = new(repository.PgRepo)
	repoif = new(repository.EsRepo)
	ctx := context.Background()
	repoif = repoif.New(ctx, *storageName, *storageName1)
	defer repoif.CloseConn()
	//init app struct
	app := endpoint.App{}
	app.CTX = ctx
	app.Repository = repoif

	//создание сервера с таким портом, и обработчиком интерфейс которого связывается а файлохранилищем
	// т.к. инициализация происходит (RegisterPublicHTTP)- в интерфейс endpoint подаетмя структура из file.go
	serv := http.Server{
		Addr:    net.JoinHostPort("", *port),
		Handler: app.RegisterPublicHTTP(),
	}
	// запуск сервера
	go func() {
		if err := serv.ListenAndServe(); err != nil {
			log.Fatalf("listen and serve err: %v", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	log.Print("Started app, REST API port = ", *port)
	// ждет сигнала
	sig := <-interrupt

	log.Printf("Sig: %v, stopping app", sig)
	// шат даун по контексту с тайм аутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*shutdownTimeout)*time.Second)
	defer cancel()
	if err := serv.Shutdown(ctx); err != nil {
		log.Printf("shutdown err: %v", err)
	}
}