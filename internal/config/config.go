package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultMaxSize    = 100
	defaultMaxBackups = 2
	defaultMaxAge     = 10
)

type Config struct {
	botToken   string        // token to connect telegram api
	botAddress string        // boot ip ore server name from hooks
	botPort    int           // boot port from hooks
	apiAddr    string        // address to API from nginx
	certPath   string        // path to certificate
	keyPath    string        // path to key
	logLev     zerolog.Level // log Level
	logFile    string
	answerText // TemplateText
}

type answerText struct {
	help        string
	start       string
	startAdmin  string
	youID       string
	error       string
	integration string
}

func (c *answerText) TextHelp() string {
	return c.help
}

func (c *answerText) TextStart() string {
	return c.start
}

func (c *answerText) TextStartAdmin() string {
	return c.startAdmin
}

func (c *answerText) TextYouID() string {
	return c.youID
}

func (c *answerText) TextError() string {
	return c.error
}
func (c *answerText) TextIntegration() string {
	return c.integration
}

func (c Config) String() (result string) {
	return fmt.Sprintf("%#v", c)
}

func (c *Config) BotToken() string {
	return c.botToken
}

func (c *Config) BotAddress() string {
	return c.botAddress
}

func (c *Config) BotPort() int {
	return c.botPort
}

func (c *Config) APIAddress() string {
	return c.apiAddr
}

func (c *Config) LogLevel() zerolog.Level {
	return c.logLev
}

func (c *Config) LogFileName() string {
	return c.logFile
}

func (c *Config) Cert() string {
	return c.certPath
}

func (c *Config) Key() string {
	return c.keyPath
}

// NewConfig - конструктор конфигурации программы
func NewConfig() *Config {
	cfg := new(Config)
	cfg.loadFromViper()

	return cfg
}

// loadFromViper - метод загрузки конфигурации из конфигурационного файла
func (c *Config) loadFromViper() {
	viper.SetConfigName("smart48bot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		err := CreateConfigFile()
		if err != nil {
			log.Fatal().Msgf("Fatal error Create config file: %v \n", err)
		}
		log.Fatal().Msg("Config file crated - smart48bot.yaml.example")
	}
	if viper.GetString("telegram_token") == "" {
		log.Error().Msgf("Не указано telegram_token в конфигурационном файле")

		return
	}
	c.botToken = viper.GetString("telegram_token")
	if viper.GetString("bot_addr") == "" {
		log.Error().Msgf("Не указано bot_addr в конфигурационном файле")

		return
	}
	c.botAddress = viper.GetString("bot_addr")
	if viper.GetInt("bot_port") == 0 {
		log.Error().Msgf("Не указано bot_port в конфигурационном файле")

		return
	}
	c.botPort = viper.GetInt("bot_port")
	if viper.GetString("api_addr") == "" {
		log.Error().Msgf("Не указано api_addr в конфигурационном файле")

		return
	}
	c.apiAddr = viper.GetString("api_addr")
	if viper.GetString("cert_path") == "" {
		log.Error().Msgf("Не указано cert_path в конфигурационном файле")

		return
	}
	c.certPath = viper.GetString("cert_path")
	if viper.GetString("key_path") == "" {
		log.Error().Msgf("Не указано key_path в конфигурационном файле")

		return
	}
	c.keyPath = viper.GetString("key_path")
	if viper.GetString("log.level") == "" {
		log.Error().Msgf("Не указано log.level в конфигурационном файле")

		return
	}
	ll, err := zerolog.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		log.Fatal().Msgf("Fatal error parse LogLevel: %v \n", err)
	}
	c.logLev = ll
	if viper.GetString("log.file_name") != "" {
		c.logFile = viper.GetString("log.file_name")
	}
	if viper.GetString("text.help") == "" {
		log.Error().Msgf("Не указано text.help в конфигурационном файле")

		return
	}
	c.answerText.help = viper.GetString("text.help")
	if viper.GetString("text.start") == "" {
		log.Error().Msgf("Не указано text.start в конфигурационном файле")

		return
	}
	c.answerText.start = viper.GetString("text.start")
	if viper.GetString("text.start_admin") == "" {
		log.Error().Msgf("Не указано text.start_admin в конфигурационном файле")

		return
	}
	c.answerText.startAdmin = viper.GetString("text.start_admin")
	if viper.GetString("text.you_id") == "" {
		log.Error().Msgf("Не указано text.you_id в конфигурационном файле")

		return
	}
	c.answerText.youID = viper.GetString("text.you_id")
	if viper.GetString("text.error") == "" {
		log.Error().Msgf("Не указано text.error в конфигурационном файле")

		return
	}
	c.answerText.error = viper.GetString("text.error")
	if viper.GetString("text.integration") == "" {
		log.Error().Msgf("Не указано text.integration в конфигурационном файле")

		return
	}
	c.answerText.integration = viper.GetString("text.integration")
	c.ConfigureLogger()
}

// ConfigureLogger - метод конфигуратор логера и ротации лог-файлов
func (c *Config) ConfigureLogger() {
	zerolog.SetGlobalLevel(c.LogLevel())
	z := zerolog.New(&lumberjack.Logger{
		Filename:   c.logFile,         // Имя файла
		MaxSize:    defaultMaxSize,    // Размер в МБ до ротации файла
		MaxBackups: defaultMaxBackups, // Максимальное количество файлов, сохраненных до перезаписи
		MaxAge:     defaultMaxAge,     // Максимальное количество дней для хранения файлов
		Compress:   true,              // Следует ли сжимать файлы логов с помощью gzip
	}).With().Timestamp().Logger()
	log.Logger = z
}

func CreateConfigFile() (err error) {
	configExample := `telegram_token: "you bot_token"
bot_addr: "you bot webhook address"
bot_port: 8443
api_addr: "127.0.0.1"
cert_path: "cert.pem"
key_path: "privkey.pem"
log:
  file_name: "smart48bot.log"
  level: "error"
text:
  help: "Для отправки сообщения (видео, фото) конкретному пользователю понадобится его ID, который автоматически присваивается всем подключенным к боту.
  Также получения/установка статуса устройства.
  В меню Интеграция приведен пример как послать сообщение. ID пользователя постоянный и не меняется пока бот не будет удален из списка контактов."
  start: "Добро пожаловать!\n\n
Наш Телеграм-бот создан для отправки сообщений пользователям умных домов MimiSmart.
Выберите действие на клавиатуре ниже."
  start_admin: "Запущен бот у нового пользователя. ID: "
  you_id: "Ваш ID: "
  error: "Команда не распознана
  Выберите действие на клавиатуре ниже."
  integration: "Версия 0.0.1\n\n
Умеет отправлять еще фото и видеофайлы. Разработана совместно со скриптом интеграции 0.0.1\n\n
Данные в адрес скрипта https://serv.smart48.ru/api/smart48/ передаются методом POST.\n
Для отправки image: https://serv.smart48.ru/api/smart48/image/?file=имя_файла&chat_id=id_получателя&text=текст_сообщения\n
Для отправки video: https://serv.smart48.ru/api/smart48/video/?file=имя_файла&chat_id=id_получателя&text=текст_сообщения\n
Для отправки file: https://serv.smart48.ru/api/smart48/file/?file=имя_файла&chat_id=id_получателя&text=текст_сообщения\n
Для отправки текста: https://serv.smart48.ru/api/smart48/msg/?chat_id=id_получателя&text=текст_сообщения\n

\n\n
Для Cuarm5 доступны только endpoint hex и tg_hex.php в формате HEX с разделителем ||
id_получателя||текст_сообщения
https://serv.smart48.ru/hex/?hex=сообщение
"`
	file, err := os.Create("smart48bot.yaml.example")
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(configExample)

	return nil
}
