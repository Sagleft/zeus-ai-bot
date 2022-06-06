package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
)

type errorFunc func() error

func checkErrors(errChecks ...errorFunc) error {
	for _, errFunc := range errChecks {
		err := errFunc()
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *solution) parseConfig() error {
	print("parse config..")

	// parse config file
	if _, err := os.Stat(configJSONPath); os.IsNotExist(err) {
		return errors.New("failed to find config file")
	}

	jsonBytes, err := ioutil.ReadFile(configJSONPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, &app.Config)
}

func getFloat(val float32) *float32 {
	return &val
}

func getInt(val int) *int {
	return &val
}

func wrapPrintedMessage(info string) string {
	return "[ " + info + " ]"
}

func printSuccess(info string) {
	color.Green(wrapPrintedMessage(info))
}

func print(info string) {
	fmt.Println(info)
}

func printError(info string) {
	color.Red("ERROR: " + info)
}
