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
	"os/user"
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
}

func notifyIfErr(err error)  {
	if err != nil {
		notifySend(err.Error())
	}
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

func getConfig(fileName string) YandexConfig {
	var config YandexConfig

	configData, ioError := ioutil.ReadFile(fileName)
	notifyIfErr(ioError)

	unMarshalError := json.Unmarshal(configData, &config)
	notifyIfErr(unMarshalError)

	return config
}

func getSelectedText() string {
	cmdXclip := exec.Command("xclip", "-o")
	var out bytes.Buffer

	cmdXclip.Stdout = &out
	err := cmdXclip.Run()
	notifyIfErr(err)

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

	usr, getUserError := user.Current()

	notifyIfErr(getUserError)

	lang := flag.String("lang", "en-ru","Translate direction")
	configFileName := flag.String("config", usr.HomeDir + "/.yagotranslate/config.json", "Config file name location")
	flag.Parse()

	config := getConfig(*configFileName)

	requestParams := TranslateApiRequest{
		config.ApiUrl,
		config.ApiKey,
		*lang,
		getSelectedText(),
	}

	translatedText := apiRequest(requestParams)
	notifySend(translatedText)
}