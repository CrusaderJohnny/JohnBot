package internal

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "say",
		Description: "says what you say",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "message",
				Description: "What to say",
				Required:    true,
			},
			discord.ApplicationCommandOptionBool{
				Name:        "ephemeral",
				Description: "If the response should only be visible to you",
				Required:    true,
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "buttons",
		Description: "Show clickable buttons",
	},
	discord.SlashCommandCreate{
		Name:        "selectmenu",
		Description: "Show select menu",
	},
	discord.SlashCommandCreate{
		Name:        "selectrole",
		Description: "Show role selection menu",
	},
	discord.SlashCommandCreate{
		Name:        "selectchannel",
		Description: "Show channel selection menu",
	},
	discord.SlashCommandCreate{
		Name:        "selectuser",
		Description: "Show user selection menu",
	},
	discord.SlashCommandCreate{
		Name:        "mentionablemenu",
		Description: "Show mentionable menu for role selection",
	},
}
