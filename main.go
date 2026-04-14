package main

import (
	"JohnBot/internal"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

type Config struct {
	Token   string
	GuildID snowflake.ID
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		Token:   os.Getenv("DISGO_TOKEN"),
		GuildID: snowflake.GetEnv("DISGO_GUILD_ID"),
	}
}

func main() {

	// Loads config file with tokens and serverID
	cfg := LoadConfig()

	// Client established, WithEventListeners for routing handler, WithEventListenerFunc for no router
	client, err := disgo.New(cfg.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		), bot.WithEventListenerFunc(buttonClicker),
		bot.WithEventListenerFunc(clickButton),
	)

	if err != nil {
		slog.Error("error while building disgo instance", slog.Any("error", err))
		return
	}
	// Defer closing with to do context for now. Can select a specific context later
	defer client.Close(context.TODO())

	// registers commands with the serverID allowing instant command updates
	if _, err = client.Rest.SetGuildCommands(client.ApplicationID, cfg.GuildID, internal.Commands); err != nil {
		slog.Error("error while registering commands", slog.Any("error", err))
	}

	// registers commands for the bot itself, takes up to an hour for these to register
	if _, err = client.Rest.SetGlobalCommands(client.ApplicationID, internal.Commands); err != nil {
		slog.Error("error while registering commands", slog.Any("error", err))
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while opening disgo gateway", slog.Any("error", err))
	}

	slog.Info("bot is now running.  Press CTRL-C to exit.")

	// Create channel to listen for sigint
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

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
