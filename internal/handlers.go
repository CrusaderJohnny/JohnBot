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
	// assign the slash command data to variable
	c := e.SlashCommandInteractionData()
	// switch logic based on custom command name sent to bot
	switch c.CommandName() {
	/*
		Say command
		Usage: "/say" [message] 'string' [ephemeral] 'bool'
		Description: Type /say in discord channel. Bot responds with two fields, [message] for the message you wish the bot to say, and [ephemeral] for if only you wish to see it.
	*/
	case "say":
		return e.CreateMessage(discord.NewMessageCreate().WithContent(data.String("message")).WithEphemeral(data.Bool("ephemeral")))
	/*
		Buttons command
		Usage: "/buttons"
		Description: Type /buttons in discord channel. Buttons will appear for the user that called the command.
		Ephemeral Flag hides the message for others so only user who calls command may see and click on the buttons
	*/
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
	/*
		Select Menu Command
		Usage: "/selectmenu"
		Description: Type /selectmenu to generate a selection menu. User may click on different colours to select them.
	*/
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
	/*
		Select User Command
		Usage: "/selectuser"
		Description: Type /selectuser to generate a drop down menu of users within the server.
	*/
	case "selectuser":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Select a user",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewUserSelectMenu("user_select", "Select a user")),
			},
			Flags: discord.MessageFlagEphemeral,
		})
	/*
		Select Role Command
		Usage: "/selectrole"
		Description: Type /selectrole to generate a drop down menu of roles within the server.
	*/
	case "selectrole":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Select a role",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewRoleSelectMenu("role_select", "Select a role")),
			},
			Flags: discord.MessageFlagEphemeral,
		})
	/*
		Select Channel Command
		Usage: "/selectchannel"
		Description: Type /selectchannel to generate a drop down menu of channels within the server.
	*/
	case "selectchannel":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Select a channel",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewChannelSelectMenu("channel_select", "Select a channel")),
			},
			Flags: discord.MessageFlagEphemeral,
		})
	/*
		Mentionable Menu Command
		Usage: "/mentionablemenu"
		Description: Type /mentionablemenu to genereate a drop down menu of users or roles within the server.
	*/
	case "mentionablemenu":
		return e.CreateMessage(discord.MessageCreate{
			Content: "Choose a mentionable menu",
			Components: []discord.LayoutComponent{
				discord.NewActionRow(
					discord.NewMentionableSelectMenu("mention_select", "Select a user or role")),
			},
			Flags: discord.MessageFlagEphemeral,
		})
	/*
		Generate Modal Command
		Usage: "/modal"
		Description: Type /modal to generate a pop out modal in discord with text fields.
	*/
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

/*
Modal Interactions Handler
Receives all modal interaction events and routes them according to customIDs
*/
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
