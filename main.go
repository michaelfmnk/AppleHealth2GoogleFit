package main

import (
	"fmt"
	"github.com/alecthomas/kong"
)

const (
	Scope = "https://www.googleapis.com/auth/fitness.body.write"
)

type CLI struct {
	ParseAppleHealthXML ParseAppleHealthXML `cmd:"" help:"Parse Apple Health XML file" name:"parse"`
	ImportGoogleFitness ImportGoogleFitness `cmd:"" help:"Import Google Fitness data" name:"import"`
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli)
	var err error
	switch ctx.Command() {
	case "parse-apple-health-xml":
		err = cli.ParseAppleHealthXML.Run()
	case "import-google-fitness":
		err = cli.ImportGoogleFitness.Run()
	default:
		_ = kong.DefaultHelpPrinter(kong.HelpOptions{}, ctx)
	}
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
