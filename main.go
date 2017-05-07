package main
import (
	"flag"
	"net/http"
	"io/ioutil"
	"os/exec"
	"bytes"
	"log"
	"encoding/json"
	"strings"
	"net/url"
)

type TranslateApiResponse struct {
	Code int
	Lang string
	Text []string
}

type TranslateApiRequest struct {
	ApiUrl string
	ApiKey string
	Lang string
	Phrase string
}

type YandexConfig struct {
	ApiKey string
	ApiUrl string
	DefaultTranslateDirection string
}

func apiRequest(requestParams TranslateApiRequest) string {
	var apiResponse TranslateApiResponse

	phraseToTranslate := url.QueryEscape(requestParams.Phrase)

	resp, requestError := http.Get(requestParams.ApiUrl + "?key=" + requestParams.ApiKey + "&lang=" + requestParams.Lang + "&text=" + phraseToTranslate)

	if requestError != nil {
		notifySend(requestError.Error())
	}

	defer resp.Body.Close()

	body, ioError := ioutil.ReadAll(resp.Body)

	if ioError != nil {
		notifySend(ioError.Error())
	}

	unMarshalError := json.Unmarshal(body, &apiResponse)

	if unMarshalError != nil {
		notifySend(unMarshalError.Error())
	}

	return strings.Join(apiResponse.Text, "\n")
}

func getConfig() YandexConfig {
	var config YandexConfig

	configData, err := ioutil.ReadFile("config.json")

	if err != nil {
		notifySend(err.Error())
	}

	unMarshalError := json.Unmarshal(configData, &config)

	if unMarshalError != nil {
		notifySend(unMarshalError.Error())
	}

	return config
}

func getSelectedText() string {
	cmdXclip := exec.Command("xclip", "-o")
	var out bytes.Buffer

	cmdXclip.Stdout = &out
	err := cmdXclip.Run()

	if err != nil {
		notifySend(err.Error())
	}

	return out.String()
}

func notifySend(message string) {
	cmdNotify := exec.Command("notify-send", "Yandex Translate", message)
	notifyRunErr := cmdNotify.Run()

	if notifyRunErr != nil {
		log.Fatal(notifyRunErr)
	}
}

func main() {

	config := getConfig()
	selectedText := getSelectedText()
	lang := flag.String("lang",config.DefaultTranslateDirection,"Translate direction")
	flag.Parse()

	requestParams := TranslateApiRequest{
		config.ApiUrl,
		config.ApiKey,
		*lang,
		selectedText,
	}

	translatedText := apiRequest(requestParams)
	notifySend(translatedText)
}