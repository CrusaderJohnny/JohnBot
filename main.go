package main

import (
	"JohnBot/internal"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

/*
Config struct for holding general information to be passed around by the bot and used in various ways.
Currently, has fields for bot token and for guild/server ID
*/
type Config struct {
	Token   string
	GuildID snowflake.ID
}

/*
Function to load config struct.
Assigns values from .env file
*/
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

	/*
		Handler instantiation and assignment for slash commands.
		All commands require registration via SlashCommand as well as the pattern, including the '/' followed by the handler
		Can route all commands through a single handler or design seperate handlers for all commands.
	*/
	r := handler.New()
	r.SlashCommand("/say", internal.CommandsHandler)
	r.SlashCommand("/buttons", internal.CommandsHandler)
	r.SlashCommand("/selectmenu", internal.CommandsHandler)
	r.SlashCommand("/selectuser", internal.CommandsHandler)
	r.SlashCommand("/selectrole", internal.CommandsHandler)
	r.SlashCommand("/selectchannel", internal.CommandsHandler)
	r.SlashCommand("/mentionablemenu", internal.CommandsHandler)
	r.SlashCommand("/modal", internal.CommandsHandler)

	/*
		Instantiate the client with bot info.
		Can use default gateway or with config options to specify functionality and/or restrict features.
		Uses WithEventListeners for routing handler, WithEventListenerFunc for no router.
		Can also use events.ListenerAdapter{} struct to route all events the bot registers based on type of event.
	*/
	client, err := disgo.New(cfg.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		), bot.WithEventListeners(r),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnComponentInteraction: internal.HandlerComponentInteractions,
			OnModalSubmit:          internal.HandlerModalInteractions,
		}),
	)

	if err != nil {
		slog.Error("error while building disgo instance", slog.Any("error", err))
		return
	}
	// Defer closing with to do context for now. Can select a specific context later
	defer client.Close(context.TODO())

	// registers commands with the serverID allowing instant command updates
	if _, err = client.Rest.SetGuildCommands(client.ApplicationID, cfg.GuildID, internal.Commands); err != nil {
		slog.Error("error while registering guild commands", slog.Any("error", err))
	}

	/* registers commands for the bot itself, takes up to an hour for these to register
	if _, err = client.Rest.SetGlobalCommands(client.ApplicationID, []discord.ApplicationCommandCreate{}); err != nil {
		slog.Error("error while registering user commands", slog.Any("error", err))
	}
	*/

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while opening disgo gateway", slog.Any("error", err))
	}

	slog.Info("bot is now running.  Press CTRL-C to exit.")

	// Create channel to listen for sigint
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
