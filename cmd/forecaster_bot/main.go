package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	bot "github.com/babadro/forecaster/internal/core/forecaster_bot"
	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/caarlos0/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 15 * time.Second
	shutdownTimeout = 5 * time.Second
)

type envVars struct {
	HTTPAddr       string `env:"HTTP_ADDR,required"`
	TelegramToken  string `env:"TELEGRAM_TOKEN,required"`
	DBConn         string `env:"DB_CONN,required"`
	NgrokAuthToken string `env:"NGROK_AUTH_TOKEN,required"`
}

func main() {
	var envs envVars

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	if err := env.Parse(&envs); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	dbPool, err := pgxpool.Connect(context.Background(), envs.DBConn)
	if err != nil {
		log.Fatalf("Unable to connection to database :%v\n", err)
	}
	defer dbPool.Close()

	forecastDB := postgres.NewForecastDB(dbPool)

	_ = bot.NewService(forecastDB)

	tunnel, err := ngrokRun(context.Background(), envs.NgrokAuthToken)
	if err != nil {
		log.Printf("Unable to run ngrok: %v\n", err)
		return
	}

	tgBot, err := initBot(tunnel.URL(), envs.TelegramToken)
	if err != nil {
		log.Printf("Unable to init bot: %v\n", err)
		return
	}

	var wg sync.WaitGroup

	server := &http.Server{
		Addr: envs.HTTPAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			update, updateErr := tgBot.HandleUpdate(r)
			if updateErr != nil {
				errMsg, _ := json.Marshal(map[string]string{"error": updateErr.Error()})
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(errMsg)
				return
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				go processUpdate(update, tgBot)
			}()
		}),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	errChan := make(chan error)

	go func() {
		if serveErr := server.Serve(tunnel); serveErr != nil {
			errChan <- fmt.Errorf("unable to serve: %w", serveErr)
		}
	}()

	select {
	case <-c:
		log.Println("Shutting down...")
	case err = <-errChan:
		log.Printf("Error: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Unable to shutdown: %v\n", err)
	}

	wg.Wait()
}

func processUpdate(upd *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "hello from bot")
	if _, sendErr := bot.Send(msg); sendErr != nil {
		log.Printf("Unable to send message: %v\n", sendErr)
	}
}

func ngrokRun(ctx context.Context, token string) (ngrok.Tunnel, error) {
	tun, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(token),
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
