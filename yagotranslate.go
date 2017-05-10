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

type InputArguments struct {
	Lang string
	ConfigFileName string
}

func main() {
	arguments := getInputArguments()
	config := getConfig(arguments.ConfigFileName)

	requestParams := TranslateApiRequest{
		config.ApiUrl,
		config.ApiKey,
		arguments.Lang,
		getSelectedText(),
	}

	translatedText := apiRequest(requestParams)
	notifySend(translatedText)
}

func getInputArguments() InputArguments {
	var arguments InputArguments

	usr, getUserError := user.Current()

	notifyIfErr(getUserError)

	arguments.Lang = *flag.String("lang", "en-ru","Translate direction")
	arguments.ConfigFileName = *flag.String("config", usr.HomeDir + "/.yagotranslate/config.json", "Config file name location")
	flag.Parse()

	return arguments
}

func apiRequest(requestParams TranslateApiRequest) string {
	var apiResponse TranslateApiResponse

	phraseToTranslate := url.QueryEscape(requestParams.Phrase)

	resp, requestError := http.Get(requestParams.ApiUrl + "?key=" + requestParams.ApiKey + "&lang=" + requestParams.Lang + "&text=" + phraseToTranslate)
	notifyIfErr(requestError)
	defer resp.Body.Close()

	body, ioError := ioutil.ReadAll(resp.Body)
	notifyIfErr(ioError)

	unMarshalError := json.Unmarshal(body, &apiResponse)
	notifyIfErr(unMarshalError)

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

func notifyIfErr(err error)  {
	if err != nil {
		notifySend(err.Error())
	}
}

func notifySend(message string) {
	cmdNotify := exec.Command("notify-send", "Yandex Translate", message)
	notifyRunErr := cmdNotify.Run()

	if notifyRunErr != nil {
		log.Fatal(notifyRunErr)
	}
}