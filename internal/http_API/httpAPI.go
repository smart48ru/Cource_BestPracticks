package http_API

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

const (
	defaultReadTimeout       = 30
	defaultWriteTimeout      = 30
	defaultReadHeaderTimeout = 20
	defaultStatusCodeOK      = 200
	defaultStatusCodeError   = 400
)

//go:generate mockgen -source=httpAPI.go -destination=mocks/mock.go
type Sender interface {
	Send(c tgBotAPI.Chattable) (tgBotAPI.Message, error)
}

type BotInformer interface {
	BotInfo() *tgBotAPI.BotAPI
}

type HTTPServer struct {
	server  string
	port    int
	botSend Sender
	// tgKeyBoard Keyboarder
}

type Configurator interface {
	BotPort() int
	APIAddress() string
}

// Временно отключен . Нужно уточнит нужно ли отправлять клавиатуру если она не менялась
// type Keyboarder interface {
//	NewKeyboard()
//	Keyboard() tgBotAPI.ReplyKeyboardMarkup
//}

// NewHTTP - конструктор HttpServer
func NewHTTP(conf Configurator, b BotInformer) *HTTPServer {
	httpSer := new(HTTPServer)
	httpSer.server = conf.APIAddress()
	httpSer.port = conf.BotPort()
	httpSer.botSend = b.BotInfo()
	httpSer.HandleRequest()

	return httpSer
}

func (h HTTPServer) Server() (server string) {
	return h.server
}

func (h HTTPServer) Port() (port int) {
	return h.port
}

func (h HTTPServer) String() (result string) {
	return fmt.Sprintf("%#v", h)
}

// HandleRequest - метод описывающий все endpoint и handlerFunc к ним
func (h HTTPServer) HandleRequest() {
	http.HandleFunc("/api", h.api)
	http.HandleFunc("/api/", h.api)
	http.HandleFunc("/tg_send.php", h.api)
	http.HandleFunc("/tg_hex.php", h.hexAPI)
	http.HandleFunc("/api/smart48/hex/", h.hexAPI)
	http.HandleFunc("/api/smart48/image/", h.imageAPI)
	http.HandleFunc("/api/smart48/video/", h.videoAPI)
	http.HandleFunc("/api/smart48/file/", h.fileAPI)
	http.HandleFunc("/api/smart48/msg/", h.msgAPI)
	http.HandleFunc("/api/smart48/hex", h.hexAPI)
	http.HandleFunc("/api/smart48/image", h.imageAPI)
	http.HandleFunc("/api/smart48/video", h.videoAPI)
	http.HandleFunc("/api/smart48/file", h.fileAPI)
	http.HandleFunc("/api/smart48/msg", h.msgAPI)

	str := fmt.Sprintf("%s:%d", h.server, h.port)
	srv := &http.Server{
		Addr:              str,
		ReadTimeout:       defaultReadTimeout * time.Second,
		WriteTimeout:      defaultWriteTimeout * time.Second,
		ReadHeaderTimeout: defaultReadHeaderTimeout * time.Second,
	}

	go srv.ListenAndServe()
}

// imageAPI - хендлер обработки отправки изображений
func (h HTTPServer) imageAPI(w http.ResponseWriter, r *http.Request) {
	ch := r.FormValue("chat_id")
	text := r.FormValue("text")
	chatID, err := strconv.ParseInt(ch, 10, 64) //nolint:gomnd
	if err != nil {
		stringErr := fmt.Sprintf("Ошибка конвертации chatID %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, defaultStatusCodeError)

		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		stringErr := fmt.Sprintf("Ошибка чтений из file= %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, defaultStatusCodeError)

		return
	}

	if file == nil {
		return
	}
	// настраиваем файл для отправки
	fReader := tgBotAPI.FileReader{
		Name:   fileHeader.Filename,
		Reader: file,
		Size:   fileHeader.Size,
	}
	msgFile := tgBotAPI.NewPhotoUpload(chatID, fReader)
	log.Debug().Msgf("структура FileReader = %v", fReader)
	h.botSend.Send(msgFile)

	if text == "" {
		http.Error(w, "OK", defaultStatusCodeOK)

		return
	}

	msg := tgBotAPI.NewMessage(chatID, text)
	log.Debug().Msgf("Строка отправки в бот = %v", msg)
	h.botSend.Send(msg)
	http.Error(w, "OK", defaultStatusCodeOK)
}

// videoAPI - хендлер обработки отправки видео

func (h HTTPServer) videoAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200) //nolint:gomnd
	ch := r.FormValue("chat_id")
	text := r.FormValue("text")
	chatID, err := strconv.ParseInt(ch, 10, 64) //nolint:gomnd
	if err != nil {
		stringErr := fmt.Sprintf("Ошибка конвертации chatID %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, defaultStatusCodeError)

		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		stringErr := fmt.Sprintf("Ошибка чтений из file= %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, defaultStatusCodeError)

		return
	}

	if file == nil {
		return
	}
	fReader := tgBotAPI.FileReader{
		Name:   fileHeader.Filename,
		Reader: file,
		Size:   fileHeader.Size,
	}
	msgFile := tgBotAPI.NewVideoUpload(chatID, fReader)
	log.Debug().Msgf("структура FileReader = %v", fReader)
	h.botSend.Send(msgFile)

	if text == "" {
		http.Error(w, "OK", defaultStatusCodeOK)

		return
	}

	msg := tgBotAPI.NewMessage(chatID, text)
	log.Debug().Msgf("Строка отправки в бот = %v", msg)
	h.botSend.Send(msg)
	http.Error(w, "OK", defaultStatusCodeOK)
}

// fileAPI - хендлер обработки отправки файлов
func (h HTTPServer) fileAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200) //nolint:gomnd
	ch := r.FormValue("chat_id")
	text := r.FormValue("text")
	chatID, err := strconv.ParseInt(ch, 10, 64) //nolint:gomnd
	if err != nil {
		stringErr := fmt.Sprintf("Ошибка конвертации chatID %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, 400) //nolint:gomnd

		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		stringErr := fmt.Sprintf("Ошибка чтений из file= %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, 400) //nolint:gomnd

		return
	}

	if file == nil {
		return
	}
	fReader := tgBotAPI.FileReader{
		Name:   fileHeader.Filename,
		Reader: file,
		Size:   fileHeader.Size,
	}
	msgFile := tgBotAPI.NewDocumentUpload(chatID, fReader)
	log.Debug().Msgf("структура FileReader = %v", fReader)
	h.botSend.Send(msgFile)

	if text == "" {
		http.Error(w, "OK", 200) //nolint:gomnd

		return
	}

	msg := tgBotAPI.NewMessage(chatID, text)
	log.Debug().Msgf("Строка отправки в бот = %v", msg)
	h.botSend.Send(msg)
	http.Error(w, "OK", defaultStatusCodeOK)
}

// msgAPI - хендлер обработки отправки сообщений
func (h HTTPServer) msgAPI(w http.ResponseWriter, r *http.Request) {
	ch := r.FormValue("chat_id")
	text := r.FormValue("text")
	if ch == "" {
		http.Error(w, "No chat_id", defaultStatusCodeError)

		return
	}
	chatID, err := strconv.ParseInt(ch, 10, 64) //nolint:gomnd
	if err != nil {
		stringErr := fmt.Sprintf("Ошибка конвертации chatID %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, defaultStatusCodeError)

		return
	}

	if text == "" {
		stringErr := "Text is empty"
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, defaultStatusCodeError)

		return
	}
	msg := tgBotAPI.NewMessage(chatID, text)
	log.Debug().Msgf("Строка отправки в бот = %v", msg)
	h.botSend.Send(msg)
	http.Error(w, "OK", defaultStatusCodeOK)
}

// api - старый хендлер поддерживается только для старых версий клиента
// Не хочу его делать красивым, надеюсь скоро удалим этом хендлер
// Не хочу управлять кодами ответа.
func (h HTTPServer) api(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200) //nolint:gomnd
	ch := r.FormValue("chat_id")
	text := r.FormValue("text")
	bot := r.FormValue("bot")
	image, _, _ := r.FormFile("image")
	video, _, _ := r.FormFile("video")
	files, _, _ := r.FormFile("file")

	if bot != "smart48" {
		return
	}

	chatID, err := strconv.ParseInt(ch, 10, 64) //nolint:gomnd
	if err != nil {
		log.Error().Msgf("Ошибка конвертации chatID %v", err)
	}
	if video != nil {
		file, fileHeader, err := r.FormFile("video")
		if err != nil {
			log.Error().Msgf("Ошибка чтений из video= %v", err)
		}

		fReader := tgBotAPI.FileReader{
			Name:   fileHeader.Filename,
			Reader: file,
			Size:   fileHeader.Size,
		}
		msgFile := tgBotAPI.NewVideoUpload(chatID, fReader)
		log.Debug().Msgf("структура FileReader = %v", fReader)
		h.botSend.Send(msgFile)
		msg := tgBotAPI.NewMessage(chatID, text)
		log.Debug().Msgf("Строка отправки в бот = %v", msg)
		h.botSend.Send(msg)
	} else if image != nil {
		file, fileHeader, err := r.FormFile("image")
		if err != nil {
			log.Error().Msgf("Ошибка чтений из image= %v", err)
		}

		fReader := tgBotAPI.FileReader{
			Name:   fileHeader.Filename,
			Reader: file,
			Size:   fileHeader.Size,
		}
		msgFile := tgBotAPI.NewVideoUpload(chatID, fReader)
		log.Debug().Msgf("структура FileReader = %v", fReader)
		h.botSend.Send(msgFile)
		msg := tgBotAPI.NewMessage(chatID, text)
		h.botSend.Send(msg)
	} else if files != nil {
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Error().Msgf("Ошибка чтений из file= %v", err)
		}

		fReader := tgBotAPI.FileReader{
			Name:   fileHeader.Filename,
			Reader: file,
			Size:   fileHeader.Size,
		}
		msgFile := tgBotAPI.NewVideoUpload(chatID, fReader)
		log.Debug().Msgf("структура FileReader = %v", fReader)
		h.botSend.Send(msgFile)
		msg := tgBotAPI.NewMessage(chatID, text)
		log.Debug().Msgf("Строка отправки в бот = %v", msg)
		h.botSend.Send(msg)
	} else {
		h.msgAPI(w, r)
		// msg := tgBotAPI.NewVideoUpload(chatID, text)
		// h.botSend.Send(msg)
	}
}

// hexAPI - хендлер для Cuarm5 принимает и отправляет только id и текст сообщения через разделитель || в формате HEX
func (h HTTPServer) hexAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200) //nolint:gomnd
	hexDate := r.FormValue("hex")
	log.Debug().Msgf("Строка XEH = %v", hexDate)
	// конвертируем все в string
	decoded, err := hex.DecodeString(hexDate)
	if err != nil {
		stringErr := fmt.Sprintf("Строка не в HEX %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, defaultStatusCodeError)
	}
	// записываем в слайс разделитель ||
	decodedStr := string(decoded)
	splitd := strings.Split(decodedStr, "||")
	// текст сообщения
	text := splitd[1]
	// конвертируем id получателя строку в int64
	chatID, err := strconv.ParseInt(splitd[0], 10, 64) //nolint:gomnd
	if err != nil {
		stringErr := fmt.Sprintf("Ошибка конвертации chatID %v\n", err)
		log.Error().Msg(stringErr)
		http.Error(w, stringErr, defaultStatusCodeError)
	}
	msg := tgBotAPI.NewMessage(chatID, text)
	h.botSend.Send(msg)
	http.Error(w, "OK", defaultStatusCodeOK)
}
