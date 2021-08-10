// instagram @edermxf 	twitter @ederhmaia

package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gookit/color"
	"github.com/robfig/cron/v3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//estruturação da response
//Cidade is....
type Cidade struct {
	Name         string `json:"city"`
	Confirmed    int    `json:"last_available_confirmed"`
	Deaths       int    `json:"last_available_deaths"`
	NewConfirmed int    `json:"new_confirmed"`
	NewDeaths    int    `json:"new_deaths"`
}

//APIResultados is....
type APIResultados struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Cidade `json:"results"`
}

// vc pode usar o token como uma var num .env
var TOKEN string = ""

func getCity() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	requestTimeout := 120

	reqAb, _ := http.NewRequest(http.MethodGet, "https://api.brasil.io/v1/dataset/covid19/caso_full/data/?is_last=True&state=SC&city=Abelardo Luz&is_last=True", nil)

	reqAb.Header.Set("Authorization", "Token "+TOKEN)

	clientAb := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(requestTimeout) * time.Second,
	}

	resAb, err := clientAb.Do(reqAb)
	if err != nil {
		panic(err.Error())
	}

	defer resAb.Body.Close()

	bodyAb, readErr := ioutil.ReadAll(resAb.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	apiResultsAb := APIResultados{}
	jsonAb := json.Unmarshal(bodyAb, &apiResultsAb)
	if jsonAb != nil {
		log.Fatal(jsonAb)
	}

	Nome := apiResultsAb.Results[0].Name
	Confirmados := apiResultsAb.Results[0].Confirmed
	NewConfirmados := apiResultsAb.Results[0].NewConfirmed
	Mortes := apiResultsAb.Results[0].Deaths
	NewMortes := apiResultsAb.Results[0].NewDeaths

	var consumerKey string = "CONSUMER KEY DO TWITTER AQ"
	var consumerSecret string = "BLABLABLA"
	var accessToken string = "TOKEN"
	var accessSecret string = "nome da var"

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	p := message.NewPrinter(language.BrazilianPortuguese)

	tweet := p.Sprintf("\n\n🏙️ %s\n😷 %d casos / 🏥%d novos\n✝️ %d óbitos / 🏴 %d novos", Nome, Confirmados, NewConfirmados, Mortes, NewMortes)

	api := twitter.NewClient(httpClient)

	_, responseTwitter, err := api.Statuses.Update("🦠 Dados COVID19 🦠"+tweet, nil)
	ResponseFMT := string(responseTwitter.Status)
	color.Warn.Println("|  " + time.Now().String() + "Status: " + ResponseFMT)
}

func main() {

	color.White.Println(`
                                                 *******
                                 ~             *---*******
                                ~             *-----*******
                         ~                   *-------*******
                        __      _   _!__     *-------*******
                   _   /  \_  _/ \  |::| ___ **-----********   ~
                 _/ \_/^    \/   ^\/|::|\|:|  **---*****/^\_
              /\/  ^ /  ^    / ^ ___|::|_|:|_/\_******/  ^  \
             /  \  _/ ^ ^   /    |::|--|:|---|  \__/  ^     ^\___
           _/_^  \/  ^    _/ ^   |::|::|:|-::| ^ /_  ^    ^  ^   \_
          /   \^ /    /\ /       |::|--|:|:--|  /  \        ^      \
         /     \/    /  /        |::|::|:|:-:| / ^  \  ^      ^     \
   _Q   / _Q  _Q_Q  / _Q    _Q   |::|::|:|:::|/    ^ \   _Q      ^
  /_\)   /_\)/_/\\)  /_\)  /_\)  |::|::|:|:::|          /_\)
_O|/O___O|/O_OO|/O__O|/O__O|/O__________________________O|/O__________
//////////////////////////////////////////////////////////////////////
`)
	color.White.Println("|----------------------------------------------------------------------|")
	color.Magenta.Println("| + ]                   COVID TWITTER BOT!                       [ + |")
	color.White.Println("|----------------------------------------------------------------------|")
	color.White.Println("|                                                                      |")
	color.Cyan.Println("|                    INSTAGRAM           TWITTER                       |")
	color.Cyan.Println("|                     @edermxf          @ederhmaia                     |")
	color.White.Println("|                                                                      |")
	color.White.Println("|  - Using Twitter API for GoLang                                      |")
	color.White.Println("|  - Using Datasets by Brasil.io                                       |")
	color.White.Println("|----------------------------------------------------------------------|")
	color.White.Println("\n\n|----------------------------------------------------------------------|")
	color.Magenta.Println("| + ]                           BOT LOG                            [ + |")
	color.White.Println("|----------------------------------------------------------------------|")
	startCronn()
	listen()
}

//cronjob 12h
func startCronn() {
	color.Info.Println("|  " + time.Now().String() + " - BOT STARTED")
	color.White.Println("|")
	c := cron.New()

	id, _ := c.AddFunc("@every 12h", getCity)
	c.Entry(id).Job.Run()

	go c.Start()
}

//listener 
func listen() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	c := cron.New()
	c.Stop()
	color.Red.Println("\n|  " + time.Now().String() + " BOT CLOSED!!!")
}
