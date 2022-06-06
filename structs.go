package main

import (
	gpt3 "github.com/PullRequestInc/go-gpt3"
	utopiago "github.com/Sagleft/utopialib-go"
)

type solution struct {
	Config config
	Bot    *utopiago.UtopiaClient
	OpenAI gpt3.Client

	WsHandlers map[string]wsHandler
}

type wsHandler func(event utopiago.WsEvent)

type config struct {
	Utopia       utopiago.UtopiaClient `json:"utopia"`
	EnableWsSSL  bool                  `json:"enable_ws_ssl"`
	OpenAIToken  string                `json:"openai_token"`
	OpenAIEngine string                `json:"openai_engine"`
	MaxTokens    int                   `json:"max_tokens"`
}
