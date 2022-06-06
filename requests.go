package main

import (
	"context"
	"errors"
	"log"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	utopiago "github.com/Sagleft/utopialib-go"
)

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

func (app *solution) onNewAuth(event utopiago.WsEvent) {
	// get pubkey
	userPubkey, err := event.GetString("pk")
	if err != nil {
		app.onWsError(err)
		return
	}

	// approve auth
	_, err = app.Config.Utopia.AcceptAuthRequest(userPubkey, "")
	if err != nil {
		app.onWsError(errors.New("failed to accept auth: " + err.Error()))
		return
	}

	print("user " + userPubkey + " auth accepted")
	if welcomeMessage != "" {
		_, err = app.Config.Utopia.SendInstantMessage(userPubkey, welcomeMessage)
		if err != nil {
			app.onWsError(errors.New("failed to send PM: " + err.Error()))
		}
	}
}

func (app *solution) sendReply(pubkey, message string) error {
	_, err := app.Bot.SendInstantMessage(pubkey, message)
	if err != nil {
		return errors.New("failed to send reply: " + err.Error())
	}
	return nil
}

func (app *solution) onUserMessage(event utopiago.WsEvent) {
	// check message
	isMessageIncoming, err := event.GetBool("isIncoming")
	if err != nil {
		app.onWsError(err)
		return
	}
	if !isMessageIncoming {
		return
	}

	// get message text
	messageText, err := event.GetString("text")
	if err != nil {
		app.onWsError(err)
		return
	}
	if messageText == "" {
		return // ignore empty message
	}

	// get user pubkey from message
	userPubkey, err := event.GetString("pk")
	if err != nil {
		app.onWsError(err)
		return
	}
	if len(messageText) < minMessageLength {
		err := app.sendReply(userPubkey, "The message is too short. Formulate your request in as much detail as possible")
		if err != nil {
			app.onWsError(err)
			return
		}
	}

	// process request
	botResponse, err := app.handleUserRequest(messageText)
	if err != nil {
		app.onWsError(err)
		return
	}
	// send response to user
	err = app.sendReply(userPubkey, botResponse)
	if err != nil {
		app.onWsError(err)
	}
}

func (app *solution) onWsConnected() {
	printSuccess("ws connected")
}

func (app *solution) onWsEvent(event utopiago.WsEvent) {
	messageHandler, isHandlerExists := app.WsHandlers[event.Type]
	if !isHandlerExists {
		return // ignore unknown event
	}

	// process event
	messageHandler(event)
}

func (app *solution) onWsError(err error) {
	printError(err.Error())
}
