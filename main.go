package main

import (
	"encoding/json"
	"log"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
)

type DataAddRequest struct {
	DataList []string `json:"dataList"`
}

var dataProv *dataProvider

func addNewData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webHookUrl))
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.PUT("/addPhrase", addNewData)

	updates := bot.ListenForWebhook("/")

	go log.Fatal(http.ListenAndServe(":"+port, router))

	for update := range updates {
		log.Println(update.Message.Text)
		resp := dataProv.getMatchPhrase(update.Message.Text)
		bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			resp,
		))
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
