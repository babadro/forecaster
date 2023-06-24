package main

import (
	"context"
	"encoding/json"
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
	HTTPAddr       string `env:"HTTP_ADDR" envDefault:":8080"`
	TelegramToken  string `env:"TELEGRAM_TOKEN"`
	DBConn         string `env:"DB_CONN"`
	NgrokAuthtoken string `env:"NGROCK_AUTHTOKEN"`
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

	tunnel, err := ngrokRun(context.Background())

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	tgBot, err := initBot(tunnel.URL(), envVars.TelegramToken)
	if err != nil {
		log.Fatalf("Unable to init bot: %v\n", err)
	}

	go func() {
		err := http.Serve(tunnel, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			update, err := tgBot.HandleUpdate(r)
			if err != nil {
				errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(errMsg)
				return
			}

			go processUpdate(update, tgBot)
		}))

		if err != nil {
			log.Fatalf("Unable to serve: %v\n", err)
		}
	}()

	// wait for os signal
	<-c
}

func processUpdate(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "hello from bot")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Unable to send message: %v\n", err)
	}
}

func ngrokRun(ctx context.Context) (ngrok.Tunnel, error) {
	tun, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(envVars.NgrokAuthtoken),
	)

	if err != nil {
		return nil, err
	}

	log.Println("tunnel created:", tun.URL())

	return tun, nil
}

func initBot(link, token string) (*tgbotapi.BotAPI, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = true

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	wh, err := tgbotapi.NewWebhook(link)
	if err != nil {
		return nil, err
	}

	_, err = botAPI.Request(wh)
	if err != nil {
		return nil, err
	}

	info, err := botAPI.GetWebhookInfo()
	if err != nil {
		return nil, err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return botAPI, nil
}
