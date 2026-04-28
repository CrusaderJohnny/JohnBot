package internal

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

const defaultURL = "https://discord.com/"

var defaultErrorMessage = "Error: Invalid selection"

/*
Button click interactions logic
Handles the button click logic and reply for event data
*/
func buttonClickRegister(event *events.ComponentInteractionCreate, response string) {
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
	return
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
func selectMenuClickRegister(event *events.ComponentInteractionCreate, response string) {
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
	return
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
User select menu interactions logic
Handles user selection menu logic. Currently, updates the message to have a reply using @ to give a clickable link to the person
*/
func selectUserClickRegister(event *events.ComponentInteractionCreate, response string) {
	data := event.UserSelectMenuInteractionData()
	selectedUsers := data.Values
	if len(selectedUsers) > 0 {
		selectedUser := selectedUsers[0]
		response = "You picked: <@" + selectedUser.String() + ">"
		_ = event.UpdateMessage(discord.MessageUpdate{
			Content: &response,
		})
		return
	}
}

/*
Channel select menu interactions logic
Handles channel selection logic. Currently, updates the message to have a clickable link leading to the selected channel.
Will add and use this feature for channel selection on posts for bot
*/
func selectChannelClickRegister(event *events.ComponentInteractionCreate, response string) {
	data := event.ChannelSelectMenuInteractionData()
	if len(data.Channels()) <= 0 {
		_ = event.UpdateMessage(discord.MessageUpdate{
			Content: &defaultErrorMessage,
		})
		return
	}
	channelName := data.Channels()[0].Name
	channels := data.Values
	if len(channels) > 0 {
		selectedChannel := channels[0]
		completeURL := defaultURL + "channels/" + event.GuildID().String() + "/" + selectedChannel.String()
		response = "Click [HERE](" + completeURL + ") to go to `#" + channelName + "`"
		_ = event.UpdateMessage(discord.MessageUpdate{
			Content: &response,
		})
		return
	}
}

/*
Role select menu interactions logic
Allows a user to select a role via a menu and assigns it to them.
Will add more validation logic for restricted roles and filtering
*/
func selectRoleClickRegister(event *events.ComponentInteractionCreate, response string) {
	data := event.RoleSelectMenuInteractionData()
	if len(data.Roles()) <= 0 {
		_ = event.UpdateMessage(discord.MessageUpdate{
			Content: &defaultErrorMessage,
		})
		return
	}
	selectedRoleIDs := event.Member().RoleIDs
	selectedRoleIDs = append(selectedRoleIDs, data.Roles()[0].ID)
	_, err := event.Client().Rest.UpdateMember(*event.GuildID(), event.User().ID, discord.MemberUpdate{
		Roles: &selectedRoleIDs,
	})
	if err != nil {
		defaultErrorMessage = err.Error()
		_ = event.UpdateMessage(discord.MessageUpdate{
			Content: &defaultErrorMessage,
		})
	}
	response = "Succesfully updated " + event.Member().User.Username + " with the role " + data.Roles()[0].Name
	_ = event.UpdateMessage(discord.MessageUpdate{
		Content: &response,
	})
}

/*
Role mentionable menu interactions logic
Handles the mentionable role menu logic
*/
func selectMentionableMenuClickRegister(event *events.ComponentInteractionCreate, response string) {
	data := event.MentionableSelectMenuInteractionData()
	role := data.Roles()
	if len(role) > 0 {
		response = "Hey You! <@&" + role[0].ID.String() + ">"
		_ = event.CreateMessage(discord.MessageCreate{
			Content: response,
		})
	}
}

/*
Modal Interaction logic
Handles Modal interaction logic. Currently very minimal logic involved
*/
func selectModalClickRegister(event *events.ModalSubmitInteractionCreate, response string) {
	data := event.Data
	nameValue := data.Text("Name")
	descriptionValue := data.Text("Description")
	response = "You have entered '" + nameValue + "' as your name and '" + descriptionValue + "' as your description."
	_ = event.CreateMessage(discord.MessageCreate{
		Content: response,
		Flags:   discord.MessageFlagEphemeral,
	})
}
