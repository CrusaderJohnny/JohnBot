package main

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

/*
Handles click events from buttons
*/
func clickButton(event *events.ComponentInteractionCreate) {
	customID := event.ButtonInteractionData().CustomID()
	var response string
	switch strings.ToLower(customID) {
	case "option_1":
		response = "Selected Option 1"
	case "option_2":
		response = "Selected Option 2"
	case "option_3":
		response = "Selected Option 3"
	case "option_4":
		response = "Selected Option 4"
	case "github":
		response = "Selected GitHub"
	default:
		response = "Selected Unknown"
	}
	_ = event.CreateMessage(discord.MessageCreate{
		Content: response,
		Flags:   discord.MessageFlagEphemeral,
	})
}

/*
Buttons on message function with reactions to button clicks
*/
func buttonClicker(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}
	if event.Message.Content == "buttons" {
		_, _ = event.Client().Rest.CreateMessage(
			event.ChannelID,
			discord.MessageCreate{
				Content: "Choose an action:",
				Components: []discord.LayoutComponent{
					discord.NewActionRow(
						discord.NewSuccessButton("Option 1", "Option_1"),
						discord.NewSecondaryButton("Option 2", "Option_2"),
						discord.NewSuccessButton("Option 3", "Option_3"),
						discord.NewDangerButton("Option 4", "Option_4"),
						discord.NewLinkButton("github", "https://github.com/"),
					),
				},
			})
	}
}

/*
	Slash command function. saves to data variable and then if/else or match case conditions to select based off command
*/

func commandListener(appEvent *events.ApplicationCommandInteractionCreate, messageEvent *events.MessageCreate) {
	data := appEvent.SlashCommandInteractionData()
	if data.CommandName() == "say" {
		// How to create and structure the message replies. using discord.MessageCreate{} causes issues with bools
		err := appEvent.CreateMessage(discord.NewMessageCreate().WithContent(data.String("message")).WithEphemeral(data.Bool("ephemeral")))
		if err != nil {
			slog.Error("error on sending response", slog.Any("error", err))
		}
	}
	if data.CommandName() == "buttons" {
		_, _ = messageEvent.Client().Rest.CreateMessage(
			messageEvent.ChannelID,
			discord.MessageCreate{
				Content: "Choose an action:",
				Components: []discord.LayoutComponent{
					discord.NewActionRow(
						discord.NewSuccessButton("Option 1", "Option_1"),
						discord.NewSecondaryButton("Option 2", "Option_2"),
						discord.NewSuccessButton("Option 3", "Option_3"),
						discord.NewDangerButton("Option 4", "Option_4"),
						discord.NewLinkButton("github", "https://github.com/"),
					),
				},
			})
	}
}

/*
	Function to read channels and respond to ping or pong messages with the opposite
*/

func onMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}
	var message string
	if event.Message.Content == "ping" {
		message = "pong"
	} else if event.Message.Content == "pong" {
		message = "ping"
	}
	if message != "" {
		_, _ = event.Client().Rest.CreateMessage(event.ChannelID, discord.MessageCreate{
			Content: message,
		})
	}
}

/*
Handler function for /say command. Allows you to enter a message and has the bot say the message.
Ephemeral tag controls message visibility from the bot
True for only visible to person doing /say command
False if visible to others
*/
func HandleSay(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.NewMessageCreate().WithContent(data.String("message")).WithEphemeral(data.Bool("ephemeral")))
}
