package main

import (
	"os"
	"os/signal"
	"smart48-telegram-bot/internal/config"
	"smart48-telegram-bot/internal/http_API"
	"smart48-telegram-bot/internal/telegram_bot"
	"time"

	"github.com/rs/zerolog/log"

	"syscall"
)

func main() {
	// создаем новую конфигурацию приложения
	cfg := config.NewConfig()
	// создаем нового бота
	bot := telegram_bot.NewBot(cfg)
	// создаем новый HTTP Server
	http_API.NewHTTP(cfg, bot)
	// запускаем restarter в отдельной горутине
	bot.WaitG().Add(1)
	go restarted()
	go listenChannels()
	bot.WaitG().Wait()
}

// restarted - функция автоматического выхода из программы если файл существует
func restarted() {
	file := "restart.fl"
	for {
		time.Sleep(30 * time.Second)
		if _, err := os.Stat(file); err == nil {
			err := os.Remove(file)
			if err != nil {
				log.Err(err)
			}
			os.Exit(1)
		}
	}
}

func listenChannels() {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT)
	log.Debug().Msgf("My PID: %d", os.Getpid())
	for sig := range stopCh {
		switch sig {
		case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
			log.Info().Msgf("Exit program from signal %v", sig)
			os.Exit(1)
		default:
			log.Info().Msg("Unknown signal.")
		}
	}
}
