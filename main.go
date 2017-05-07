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

type YandexConfig struct {
	apiKey string
}

func apiRequest(apiKey, lang, phraseToTranslate string) string {
	var apiResponse TranslateApiResponse

	phraseToTranslate = url.QueryEscape(phraseToTranslate)

	resp, requestError := http.Get("https://translate.yandex.net/api/v1.5/tr.json/translate?key=" + apiKey + "&lang=" + lang + "&text=" + phraseToTranslate)

	if requestError != nil {
		log.Fatal(requestError)
	}

	defer resp.Body.Close()

	body, ioError := ioutil.ReadAll(resp.Body)

	if ioError != nil {
		log.Fatal(ioError)
	}

	unMarshalError := json.Unmarshal(body, &apiResponse)

	if unMarshalError != nil {
		log.Fatal(unMarshalError)
	}

	return strings.Join(apiResponse.Text, "\n")
}

func getConfig() YandexConfig {
	config := YandexConfig{
		apiKey: "trnsl.1.1.20161028T214536Z.e4ee95b3a5a57e20.6cf997cd1d39c117d1b90db5a6bbf90f22c893e4",
	}

	return config
}

func getSelectedText() string {
	cmdXclip := exec.Command("xclip", "-o")
	var out bytes.Buffer

	cmdXclip.Stdout = &out
	err := cmdXclip.Run()

	if err != nil {
		log.Fatal(err)
	}

	return out.String()
}

func notifySend(translatedText *string) {
	cmdNotify := exec.Command("notify-send", "Yandex Translate", *translatedText)
	notifyRunErr := cmdNotify.Run()

	if notifyRunErr != nil {
		log.Fatal(notifyRunErr)
	}
}

func main() {

	config := getConfig()

	selectedText := getSelectedText()

	lang := flag.String("lang","en-ru","Translate direction")

	flag.Parse()

	translatedText := apiRequest(config.apiKey, *lang, selectedText)

	notifySend(&translatedText)
}