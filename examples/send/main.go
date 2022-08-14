package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nmeilick/go-telenote"
	"github.com/urfave/cli/v2"
)

const (
	VERSION = "1.0.1"
)

func main() {
	app := &cli.App{
		Name:        "tgsend",
		Version:     VERSION,
		Usage:       "Send a telegram message",
		UsageText:   "tgsend [OPTIONS] TEXT ID...",
		Description: "Send the given text to one or more chats identified by their id.\nTo find the id, forward a message from the chat to @getidsbot.",
		HideHelp:    true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				EnvVars: []string{"TOKEN"},
				Usage:   "Telegram bot token",
			},
			&cli.BoolFlag{
				Name:    "no-preview",
				Aliases: []string{"n"},
				EnvVars: []string{"NO_PREVIEW"},
				Usage:   "Disable automatic link preview",
			},
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"m"},
				EnvVars: []string{"MODE"},
				Value:   "MarkdownV2",
				Usage:   "Parse mode (Markdown, MarkdownV2, HTML)",
			},
		},
		Action: sendMessage,
	}

	if err := app.Run(os.Args); err != nil {
		fail(err.Error())
	}
}

func fail(format string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, format)
	} else {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
	os.Exit(1)
}

func sendMessage(ctx *cli.Context) error {
	token := ctx.String("token")
	if token == "" {
		fail("Please specify the telegram bot token (-t, see -h for help)!")
	}

	var opts []telenote.Option
	if ctx.Bool("no-preview") {
		opts = append(opts, telenote.NoPreview())
	}
	if s := ctx.String("mode"); s != "" {
		opts = append(opts, telenote.ParseMode(s))
	}

	var text string
	var ids []int64

	args := ctx.Args().Slice()
	if len(args) > 0 {
		text = strings.TrimSpace(args[0])
	}

	if len(args) > 1 {
		for _, s := range args[1:] {
			id, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fail("Not a numeric chat id: %s", s)
			}
			ids = append(ids, id)
		}
	}

	if text == "" {
		fail("Please specify the message to send.")
	} else if len(ids) == 0 {
		fail("Please specify one or more chat ids to send to")
	}

	rc := 0

	notifier := telenote.NewNotifier(token)
	for _, id := range ids {
		if err := notifier.Notify(id, text, opts...); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send message to %d: %s\n", id, err)
			rc = 1
		}
	}

	os.Exit(rc)

	return nil
}
