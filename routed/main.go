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
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
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

	// Routing handler for slash commands. Will allow more complex commands without cluttering main
	r := handler.New()
	r.SlashCommand("/say", internal.HandleSay)
	r.SlashCommand("/buttons", internal.HandleCommands)

	// Client established, WithEventListeners for routing handler, WithEventListenerFunc for no router
	client, err := disgo.New(cfg.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		), bot.WithEventListeners(r),
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
