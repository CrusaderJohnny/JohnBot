package internal

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

/*
Slash commands handler function
Checks and compares the commandname to a match case then handles the logic for those commands
*/
func HandleCommands(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	c := e.SlashCommandInteractionData()
	switch c.CommandName() {
	case "say":
		return e.CreateMessage(discord.NewMessageCreate().WithContent(data.String("message")).WithEphemeral(data.Bool("ephemeral")))
	case "buttons":
		_, err := e.Client().Rest.CreateMessage(
			e.Channel().ID(),
			discord.MessageCreate{
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
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
