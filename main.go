package main

import (
	"errors"
	"log"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	utopiago "github.com/Sagleft/utopialib-go"
)

func newSolution() *solution {
	app := &solution{}
	app.WsHandlers = map[string]wsHandler{
		"newAuthorization":  app.onNewAuth,
		"newInstantMessage": app.onUserMessage,
	}
	return app
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

func (app *solution) runInBackground() error {
	printSuccess("bot started")

	ch := make(chan struct{})
	// background
	<-ch
	return nil
}

func main() {
	app := newSolution()

	err := checkErrors(
		app.parseConfig,
		app.setupOpenAIClient,
		app.utopiaConnect,
		app.runBot,
		app.runInBackground,
	)
	if err != nil {
		log.Fatalln(err)
	}
}
