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

type TranslateResponse struct {
	Code int
	Lang string
	Text []string
}

func handleWebRequests(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Дратути")
}

func main() {
	go SetupServer()

	tg_token := os.Getenv("TG_TOKEN") // TODO check if it exists
	bot, err := tgbotapi.NewBotAPI(tg_token)
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

func SetupServer(){
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

func translateWithYandex(text string) string {
	ya_token := os.Getenv("YA_TOKEN") // TODO check if it exists

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
