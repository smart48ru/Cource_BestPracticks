package telegram_bot

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// BotStruct - внутренняя структура бота
type BotStruct struct {
	bot        *tgBotAPI.BotAPI
	keyboard   tgBotAPI.ReplyKeyboardMarkup
	updateChan tgBotAPI.UpdatesChannel
	wg         sync.WaitGroup
	token      string
	addr       string
	port       int
	cert       string
	logLevel   zerolog.Level
	answerText
}

// answerText - структура ответов на команды бота. Настраивается из конфигурационного файла
type answerText struct {
	help        string
	start       string
	startAdmin  string
	youID       string
	error       string
	integration string
}

type Configurator interface {
	BotToken() string
	BotAddress() string
	BotPort() int
	Cert() string
	LogLevel() zerolog.Level
	TextYouID() string
	TextHelp() string
	TextStart() string
	TextStartAdmin() string
	TextError() string
	TextIntegration() string
}

// NewBot - создает нового бота и настраивает его

func NewBot(conf Configurator) *BotStruct {
	bot, err := tgBotAPI.NewBotAPI(conf.BotToken())
	if err != nil {
		log.Error().Err(err)
	}
	if conf.LogLevel() <= 0 {
		bot.Debug = true
	}

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)
	tgBot := new(BotStruct)
	tgBot.bot = bot
	tgBot.addr = conf.BotAddress()
	tgBot.port = conf.BotPort()
	tgBot.cert = conf.Cert()
	u := tgBotAPI.NewUpdate(0)
	u.Timeout = 60
	tgBot.updateChan, err = bot.GetUpdatesChan(u)
	if err != nil {
		log.Error().Err(err)
	}
	tgBot.logLevel = conf.LogLevel()
	tgBot.SetWebHook()
	tgBot.NewKeyboard()
	tgBot.youID = conf.TextYouID()
	tgBot.start = conf.TextStart()
	tgBot.startAdmin = conf.TextStartAdmin()
	tgBot.help = conf.TextHelp()
	tgBot.integration = conf.TextIntegration()
	tgBot.error = conf.TextError()

	tgBot.wg.Add(1)
	go tgBot.BotUpdater()
	tgBot.updateChan = bot.ListenForWebhook("/" + bot.Token)
	log.Debug().Msgf("Bot : %v", tgBot)

	return tgBot
}

// SetWebHook - метод устанавливает Webhook с сервером telegram

func (b *BotStruct) SetWebHook() {
	time.Sleep(time.Second * 2) //nolint:gomnd    // устанавливаем WebHook c задержкой, что бы не было ошибки
	str := fmt.Sprintf("https://%s/%s", b.addr, b.token)
	tgBotAPI.NewWebhookWithCert(str, b.cert)
	info, err := b.bot.GetWebhookInfo()
	if err != nil {
		log.Error().Err(err)
	}
	if info.LastErrorDate != 0 {
		log.Error().Msgf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	u := tgBotAPI.NewUpdate(0)
	u.Timeout = 5
	info.LastErrorDate = 0
}

func (b *BotStruct) UpdateChan() tgBotAPI.UpdatesChannel {
	return b.updateChan
}
func (b *BotStruct) WaitG() *sync.WaitGroup {
	return &b.wg
}

func (b *BotStruct) BotInfo() *tgBotAPI.BotAPI {
	return b.bot
}

func (b *BotStruct) Keyboard() tgBotAPI.ReplyKeyboardMarkup {
	return b.keyboard
}

// NewKeyboard - метод создающий новую панель клафиш в боте
func (b *BotStruct) NewKeyboard() {
	MyKeyboard := tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton("Мой ID"),
			tgBotAPI.NewKeyboardButton("Помощь"),
			tgBotAPI.NewKeyboardButton("Интеграция"),
		),
	)
	b.keyboard = MyKeyboard
}

// BotUpdater - метод обрабатывающий всё из updatechanel
func (b *BotStruct) BotUpdater() {
	wg := b.WaitG()
	defer wg.Done()
	adminsID := []int64{96500923, 131858248}
	for update := range b.updateChan {
		if update.Message != nil {
			log.Debug().Msgf("id: %d - [%s] %s", update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)
			command := strings.ToLower(update.Message.Text)
			switch command {
			case "ваш id", "id", "мой id", "/id":
				answer := fmt.Sprintf("%s %v", b.answerText.youID, update.Message.Chat.ID)
				msg := tgBotAPI.NewMessage(update.Message.Chat.ID, answer)
				msg.ReplyMarkup = b.keyboard
				log.Debug().Msgf("Ваш ID: %v %v", update.Message.Chat.ID, answer)
				b.bot.Send(msg)

			case "/start":
				answer := fmt.Sprintf("%v", b.answerText.start)
				msg := tgBotAPI.NewMessage(update.Message.Chat.ID, answer)
				adminAnswer1 := fmt.Sprintf("%v %v", b.answerText.startAdmin, update.Message.Chat.ID)
				// отправляем аминам
				for _, id := range adminsID {
					adminMsg := tgBotAPI.NewMessage(id, adminAnswer1)
					b.bot.Send(adminMsg)
				}
				msg.ReplyMarkup = b.keyboard
				log.Debug().Msgf("/start: %v %v", update.Message.Chat.ID, answer)
				b.bot.Send(msg)

			case "интеграция", "integration", "/integration":
				answer := fmt.Sprintf("%v", b.answerText.integration)
				msg := tgBotAPI.NewMessage(update.Message.Chat.ID, answer)
				// msg.ReplyMarkup = tgBotStruct.NewKeyboard // ответить на сообщение
				msg.ReplyMarkup = b.keyboard
				log.Debug().Msgf("Интеграция: %v %v", update.Message.Chat.ID, answer)
				b.bot.Send(msg)

			case "помощь", "help", "/help":
				answer := fmt.Sprintf("%v", b.answerText.help)
				msg := tgBotAPI.NewMessage(update.Message.Chat.ID, answer)
				msg.ReplyMarkup = b.keyboard
				log.Debug().Msgf("Помощь: %v %v", update.Message.Chat.ID, answer)
				b.bot.Send(msg)

			default:
				answer := fmt.Sprintf("%v", b.answerText.error)
				msg := tgBotAPI.NewMessage(update.Message.Chat.ID, answer)
				msg.ReplyMarkup = b.keyboard
				log.Debug().Msgf("Команда не распознана: %v %v", update.Message.Chat.ID, answer)
				b.bot.Send(msg)
			}
		}
	}
}
