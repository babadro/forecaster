package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	bot "github.com/babadro/forecaster/internal/core/forecaster_bot"
	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/caarlos0/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

var envVars = struct {
	HTTPAddr         string `env:"HTTP_ADDR" envDefault:":8080"`
	TelegramToken    string `env:"TELEGRAM_TOKEN"`
	DBConn           string `env:"DB_CONN"`
	NGROCK_AUTHTOKEN string `env:"NGROCK_AUTHTOKEN"`
}{}

func main() {
	// listen to os signals
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	if err := env.Parse(&envVars); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	dbPool, err := pgxpool.Connect(context.Background(), envVars.DBConn)
	if err != nil {
		log.Fatalf("Unable to connection to database :%v\n", err)
	}
	defer dbPool.Close()

	forecastDB := postgres.NewForecastDB(dbPool)

	_ = bot.NewService(forecastDB)

	tgbot, err := tgbotapi.NewBotAPI(envVars.TelegramToken)
	if err != nil {
		log.Fatalf("Unable to create telegram bot: %v\n", err)
	}

	log.Printf("Authorized on account %s", tgbot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// wait for os signal
	<-c
}

func ngrokRun(ctx context.Context) error {
	tun, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(envVars.NGROCK_AUTHTOKEN),
	)

	if err != nil {
		return err
	}

	log.Println("tunnel created:", tun.URL())

	return http.Serve(tun, http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello from ngrok-go.</h1>")
}

func initBot() {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := tgbotapi.NewWebhookWithCert("https://www.example.com:8443/"+bot.Token, "cert.pem")

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
