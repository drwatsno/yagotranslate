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
	"time"
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
	phraseToTranslate := url.QueryEscape(requestParams.Phrase)
	response := httpRequestGet(requestParams.ApiUrl + "?key=" + requestParams.ApiKey + "&lang=" + requestParams.Lang + "&text=" + phraseToTranslate)
	apiResponse := getParsedBodyResponse(&response)
	defer response.Body.Close()

	return strings.Join(apiResponse.Text, "\n")
}

func httpRequestGet(url string) http.Response {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	client.Get(url)
	response, requestError := client.Get(url)
	notifyIfErr(requestError)

	return *response
}

func getParsedBodyResponse(response *http.Response) TranslateApiResponse  {
	var apiResponse TranslateApiResponse

	body, ioError := ioutil.ReadAll(response.Body)
	notifyIfErr(ioError)

	unMarshalError := json.Unmarshal(body, &apiResponse)
	notifyIfErr(unMarshalError)

	return apiResponse
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