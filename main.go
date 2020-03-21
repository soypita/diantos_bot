package main

import (
	"encoding/json"
	"log"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
)

type DataAddRequest struct {
	DataList []string `json:"dataList"`
}

var dataProv *dataProvider

func addNewData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	dataList := DataAddRequest{}
	if err := json.NewDecoder(r.Body).Decode(&dataList); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	dataProv.insertNewPhrases(dataList.DataList)
	respondWithJson(w, http.StatusOK, `{"status": "success"}`)
}

func main() {
	dataProv = NewDataProvider()
	token := os.Getenv("TELEGRAM_TOKEN")
	webHookUrl := os.Getenv("WEBHOOK_URL")
	port := os.Getenv("PORT")
	log.Println(token)
	log.Println(webHookUrl)
	log.Println(port)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webHookUrl + "/" + bot.Token))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/addPhrase", addNewData)

	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServe(":"+port, nil)

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	for update := range updates {
		log.Println(update.Message.Text)
		resp := dataProv.getMatchPhrase(update.Message.Text)
		if resp != "" {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Как говорится " + strings.ToLower(resp),
			))
		}
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
