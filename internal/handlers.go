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
	case "selectmenu":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Choose a color",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewStringSelectMenu(
						"color_select",
						"Choose a color",
						discord.NewStringSelectMenuOption("Red", "red"),
						discord.NewStringSelectMenuOption("Green", "green"),
						discord.NewStringSelectMenuOption("Blue", "blue"),
					))},
		})
	}
	return nil
}

/*
Button click interactions logic
Handles the button click logic and reply for event data
*/
func ButtonClickRegister(event *events.ComponentInteractionCreate, response string) {
	// get the pressed button value from customID
	buttonID := event.ButtonInteractionData().CustomID()
	// Compares the value via switch case using strings.ToLower to ensure case matching
	switch strings.ToLower(buttonID) {
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
	// updates the message to now say the button option selected and disables the other buttons
	_ = event.UpdateMessage(discord.MessageUpdate{
		Content: &response,
		Components: &[]discord.LayoutComponent{
			discord.NewActionRow(
				discord.NewPrimaryButton("Primary Button", "Primary_1").WithDisabled(true),
				discord.NewSecondaryButton("Secondary Button", "Secondary_1").WithDisabled(true),
				discord.NewSuccessButton("Success Button", "Success_1").WithDisabled(true),
				discord.NewDangerButton("Danger Button", "Danger_1").WithDisabled(true),
			),
		},
	})
	/*
		Default message response. Sends the response back as a message only the user can see.
		Control the response via switch case assignments.
		Disabled currently for testing message updates functionality


		_ = event.CreateMessage(discord.MessageCreate{
			Content: response,
			Flags:   discord.MessageFlagEphemeral,
		})
	*/
}

/*
Select menu interactions logic
Handles select menu selection logic and reply for event data
*/
func SelectMenuClickRegister(event *events.ComponentInteractionCreate, response string) {
	// gets the menu type ID from customID
	selectID := event.StringSelectMenuInteractionData().CustomID()
	// Compares the value via switch case using strings.ToLower to ensure case matching
	switch strings.ToLower(selectID) {
	// Current color select menu. Can create multiple menu types and switch case for different types of replies.
	case "color_select":
		selectedColour := event.StringSelectMenuInteractionData().Values[0]
		response = "You selected " + selectedColour
	}
	// updates the message locking out the selection buttons. Can also have it updated to not lock out changing the response with each selection
	_ = event.UpdateMessage(discord.MessageUpdate{
		Content: &response,
		Components: &[]discord.LayoutComponent{
			discord.NewActionRow(
				discord.NewStringSelectMenu(
					"color_select",
					"Choose a color",
					discord.NewStringSelectMenuOption("Red", "red"),
					discord.NewStringSelectMenuOption("Green", "green"),
					discord.NewStringSelectMenuOption("Blue", "blue"),
				).WithDisabled(true),
			),
		},
	})
	/*
		Default message response. Sends the response back as a message only the user can see.
		Control the response via switch case assignments.
		Disabled currently for testing message updates functionality


		_ = event.CreateMessage(discord.MessageCreate{
			Content: response,
			Flags:   discord.MessageFlagEphemeral,
		})
	*/
}

/*
Improved Helper to handle a variety of interactions
Checks the type of interaction then sorts to switch cases for logic
*/
func HandlerInteractions(event *events.ComponentInteractionCreate) {
	// assign a string var to hold responses. Allows pointer to be used for updating messages
	var response string
	// switch case off of the type of data passed into
	switch event.Data.Type() {
	// button interactions
	case discord.ComponentTypeButton:
		ButtonClickRegister(event, response)
	case discord.ComponentTypeStringSelectMenu:
		SelectMenuClickRegister(event, response)
	}
}
