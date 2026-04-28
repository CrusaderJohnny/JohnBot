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
			Flags: discord.MessageFlagEphemeral,
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
			Flags: discord.MessageFlagEphemeral,
		})
	case "selectuser":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Select a user",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewUserSelectMenu("user_select", "Select a user")),
			},
			Flags: discord.MessageFlagEphemeral,
		})
	case "selectrole":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Select a role",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewRoleSelectMenu("role_select", "Select a role")),
			},
			Flags: discord.MessageFlagEphemeral,
		})
	case "selectchannel":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Select a channel",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewChannelSelectMenu("channel_select", "Select a channel")),
			},
			Flags: discord.MessageFlagEphemeral,
		})
	case "mentionablemenu":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Choose a mentionable menu",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewMentionableSelectMenu("mention_select", "Select a user or role")),
			},
			Flags: discord.MessageFlagEphemeral,
		})
	case "modal":
		return e.Modal(discord.ModalCreate{
			CustomID: "test_modal",
			Title:    "Test modal",
			Components: []discord.LayoutComponent{
				discord.LabelComponent{
					Label:     "Name",
					Component: discord.NewShortTextInput("Name").WithPlaceholder("Enter Name").WithMinLength(3).WithMaxLength(30),
				},
				discord.LabelComponent{
					Label:     "Description",
					Component: discord.NewParagraphTextInput("Description").WithPlaceholder("Enter a description here").WithMinLength(3).WithMaxLength(30),
				},
			},
		})
	}
	return nil
}

/*
Improved Helper to handle a variety of interactions
Checks the type of interaction then sorts to switch cases for logic
*/
func HandlerComponentInteractions(event *events.ComponentInteractionCreate) {
	// assign a string var to hold responses. Allows pointer to be used for updating messages
	var response string
	// switch case off of the type of data passed into
	switch event.Data.Type() {
	// button interactions
	case discord.ComponentTypeButton:
		buttonClickRegister(event, response)
	case discord.ComponentTypeStringSelectMenu:
		selectMenuClickRegister(event, response)
	case discord.ComponentTypeUserSelectMenu:
		selectUserClickRegister(event, response)
	case discord.ComponentTypeRoleSelectMenu:
		selectRoleClickRegister(event, response)
	case discord.ComponentTypeChannelSelectMenu:
		selectChannelClickRegister(event, response)
	case discord.ComponentTypeMentionableSelectMenu:
		selectMentionableMenuClickRegister(event, response)
	default:
		_ = event.CreateMessage(discord.MessageCreate{
			Content: "Unknown Interaction",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
}

func HandlerModalInteractions(event *events.ModalSubmitInteractionCreate) {
	var response string
	switch strings.ToLower(event.Data.CustomID) {
	case "test_modal":
		selectModalClickRegister(event, response)
	default:
		_ = event.CreateMessage(discord.MessageCreate{
			Content: "Unknown Interaction",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
}
