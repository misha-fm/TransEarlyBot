package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type TranslateResponse struct {
	Code int
	Lang string
	Text []string
}

func main2() {
	textToTranslate := "ich versuche weitere assets zu erstellen und dann auch noch ein paar supplies"

	translation := translateWithYandex(textToTranslate)

	fmt.Printf("Translation: %s", translation)

	//testJson()
}

func testJson() {

	rawJson := []byte(`{
    "code": 200,
    "lang": "de-ru",
    "text": [
        "я пытаюсь создать дополнительные ресурсы и потом еще пару supplies"
    ]
}`)

	var bodySerialized TranslateResponse

	err := json.Unmarshal(rawJson, &bodySerialized)
	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		code := bodySerialized.Code
		lang := bodySerialized.Lang
		text := bodySerialized.Text

		fmt.Printf("code: %d\nlang: %s\ntext: %s\n", code, lang, text[0])
	}
}

func main() {
	tg_token := os.Getenv("TG_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tg_token)
	if err != nil {
		fmt.Println("Panic!!! ")
		fmt.Print(err)
	}

	bot.Debug = false

	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		fmt.Printf("[%s] %s\n\n", update.Message.From.UserName, update.Message.Text)

		translation := translateWithYandex(update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, translation)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func translateWithNaturalIntelligence(text string) string {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Перевод: ")

	input.Scan()
	translation := input.Text()
	return translation
}

func translateWithYandex(text string) string {
	ya_token := os.Getenv("YA_TOKEN")

	var Url *url.URL

	Url, err := url.Parse("https://translate.yandex.net/api/v1.5/tr.json/translate")
	if err != nil {
		panic("AAAAA")
	}

	params := url.Values{}
	params.Add("key", ya_token)
	params.Add("lang", "de-ru")
	params.Add("text", text)
	Url.RawQuery = params.Encode()

	resp, err := http.Get(Url.String())

	//fmt.Printf("Encoded URL is %q\n", Url.String())

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
