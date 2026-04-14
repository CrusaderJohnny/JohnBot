package internal

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

/*
Slash commands handler function
Checks and compares the commandname to a match case then handles the logic for those commands
*/
func CommandsHandler(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	c := e.SlashCommandInteractionData()
	switch c.CommandName() {
	case "say":
		return e.CreateMessage(discord.NewMessageCreate().WithContent(data.String("message")).WithEphemeral(data.Bool("ephemeral")))
	case "buttons":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Choose an action:",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewPrimaryButton("Primary Button", "Primary_1"),
					discord.NewSecondaryButton("Secondary Button", "Secondary_1"),
					discord.NewSuccessButton("Success Button", "Success_1"),
					discord.NewDangerButton("Danger Button", "Danger_1"),
					discord.NewLinkButton("Link Button", "https://discord.com/channels/1492233366433562838/1493356833128054896"),
				),
			},
		})
	}
	return nil
}

/*
Helper to handle click events from buttons
*/
func ButtonClickRegister(event *events.ComponentInteractionCreate) {
	customID := event.ButtonInteractionData().CustomID()
	var response string
	switch strings.ToLower(customID) {
	case "primary_1":
		response = "Selected Option 1"
	case "secondary_1":
		response = "Selected Option 2"
	case "success_1":
		response = "Selected Option 3"
	case "danger_1":
		response = "Selected Option 4"
	default:
		response = "Selected Unknown"
	}
	_ = event.CreateMessage(discord.MessageCreate{
		Content: response,
		Flags:   discord.MessageFlagEphemeral,
	})
}
