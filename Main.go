package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"log"
	"strconv"
	"gopkg.in/telegram-bot-api.v4"
)

func handleWebRequests(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Дратути")
}

var yaToken string

func main() {
	go SetupServer()

	yaToken = os.Getenv("YA_TOKEN")  // TODO check if it exists
	tgToken := os.Getenv("TG_TOKEN") // TODO check if it exists
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		fmt.Println("Panic!!! ")
		fmt.Println(err)
	}

	bot.Debug = false

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			cbq := update.CallbackQuery

			newLang := fmt.Sprintf("ru-%s", cbq.Data)
			textToTranslate := cbq.Message.ReplyToMessage.Text
			newTranslation := translateWithYandex(textToTranslate, newLang)
			mssg := tgbotapi.NewEditMessageText(cbq.Message.Chat.ID, cbq.Message.MessageID, newTranslation)
			bot.Send(mssg)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Это было команда. Меня не проведешь!")
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
			continue
		}

		fmt.Printf("[%s] %s\n\n", update.Message.From.UserName, update.Message.Text)

		inputText := update.Message.Text;
		language := detectLanguage(inputText)
		direction := getTranslateDirection(language)

		translation := translateWithYandex(inputText, direction)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, translation)

		if translation != "ru" {
			inlineBtns := []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("english", "en"),
				tgbotapi.NewInlineKeyboardButtonData("french", "fr"),
				//tgbotapi.NewInlineKeyboardButtonData("deutsch", "de"),
			}
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineBtns)
		}

		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func getTranslateDirection(code string) string {
	switch code {
	case "ru":
		return "ru-de"
	}
	return fmt.Sprintf("%s-ru", code)
}

func SetupServer() {
	portStr := os.Getenv("PORT")

	if portStr == "" {
		// TODO handle error
	}

	port, _ := strconv.Atoi(portStr)

	httpBinding := fmt.Sprintf(":%d", port)
	http.HandleFunc("/", handleWebRequests)
	log.Fatal(http.ListenAndServe(httpBinding, nil))
}

func translateWithNaturalIntelligence(text string) string {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Перевод: ")

	input.Scan()
	translation := input.Text()
	return translation
}

func detectLanguage(text string) string {
	type DetectResponse struct {
		Code int
		Lang string
	}

	var Url *url.URL

	Url, err := url.Parse("https://translate.yandex.net/api/v1.5/tr.json/detect")
	if err != nil {
		panic("AAAAA")
	}

	params := url.Values{}
	params.Add("key", yaToken)
	params.Add("text", text)
	Url.RawQuery = params.Encode()

	resp, err := http.Get(Url.String())

	if err != nil {

	}
	defer resp.Body.Close()

	rawJson, _ := ioutil.ReadAll(resp.Body)

	//fmt.Printf("%s", rawJson)

	var bodySerialized DetectResponse
	json.Unmarshal(rawJson, &bodySerialized)

	languageCode := bodySerialized.Lang

	return languageCode
}

func translateWithYandex(text string, lang string) string {
	type TranslateResponse struct {
		Code int
		Lang string
		Text []string
	}

	var Url *url.URL

	Url, err := url.Parse("https://translate.yandex.net/api/v1.5/tr.json/translate")
	if err != nil {
		panic("AAAAA")
	}

	params := url.Values{}
	params.Add("key", yaToken)
	params.Add("lang", lang)
	params.Add("text", text)
	Url.RawQuery = params.Encode()

	resp, err := http.Get(Url.String())

	if err != nil {

	}
	defer resp.Body.Close()

	rawJson, _ := ioutil.ReadAll(resp.Body)

	//fmt.Printf("%s", rawJson)

	var bodySerialized TranslateResponse
	json.Unmarshal(rawJson, &bodySerialized)

	translation := bodySerialized.Text[0]

	return translation
}
