package main

import (
	"context"
	"errors"
	"log"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	utopiago "github.com/Sagleft/utopialib-go"
)

type solution struct {
	Config config
	Bot    *utopiago.UtopiaClient
	OpenAI gpt3.Client
}

func newSolution() solution {
	return solution{}
}

type config struct {
	Utopia       utopiago.UtopiaClient `json:"utopia"`
	EnableWsSSL  bool                  `json:"enable_ws_ssl"`
	OpenAIToken  string                `json:"openai_token"`
	OpenAIEngine string                `json:"openai_engine"`
	MaxTokens    int                   `json:"max_tokens"`
}

func (app *solution) setupOpenAIClient() error {
	print("setup OpenAI..")
	app.OpenAI = gpt3.NewClient(app.Config.OpenAIToken)
	return nil
}

func (app *solution) utopiaConnect() error {
	app.Bot = &app.Config.Utopia

	if !app.Bot.CheckClientConnection() {
		return errors.New("failed to connect to Utopia client")
	}
	return nil
}

func (app *solution) onWsConnected() {
	printSuccess("ws connected")
}

func (app *solution) onWsEvent(event utopiago.WsEvent) {
	if event.Type == "" {
		// TODO
	}
}

func (app *solution) onWsError(err error) {
	printError(err.Error())
}

func (app *solution) runBot() error {
	print("enable ws connection..")
	err := app.Bot.SetWebSocketState(utopiago.SetWsStateTask{
		Enabled:       true,
		Port:          app.Bot.WsPort,
		EnableSSL:     false,
		Notifications: "contact",
	})
	if err != nil {
		return errors.New("failed to setup Utopia websocket state: " + err.Error())
	}

	print("subscribe to events..")
	err = app.Bot.WsSubscribe(utopiago.WsSubscribeTask{
		OnConnected: app.onWsConnected,
		Callback:    app.onWsEvent,
		ErrCallback: app.onWsError,
		Port:        app.Bot.WsPort,
	})
	if err != nil {
		return errors.New("failed to subscribe to Utopia ws events: " + err.Error())
	}

	return nil
}

func main() {
	app := newSolution()

	err := checkErrors(
		app.parseConfig,
		app.setupOpenAIClient,
		app.utopiaConnect,
		app.runBot,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func (app *solution) handleUserRequest(request string) (string, error) {
	response, err := app.OpenAI.CompletionWithEngine(context.TODO(), app.Config.OpenAIEngine, gpt3.CompletionRequest{
		Prompt: []string{
			request,
		},
		Temperature:      getFloat(0.6),
		MaxTokens:        getInt(app.Config.MaxTokens),
		TopP:             getFloat(1),
		N:                getInt(1),
		FrequencyPenalty: 1,
		PresencePenalty:  1,
	})
	if err != nil {
		log.Fatalln(err)
	}

	dataArray := []string{}
	for _, data := range response.Choices {
		dataArray = append(dataArray, data.Text)
	}

	return strings.Join(dataArray, "\n"), nil
}