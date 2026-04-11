package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
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

	cfg := LoadConfig()

	client, err := disgo.New(cfg.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(commandListener), bot.WithEventListenerFunc(onMessageCreate),
	)

	if err != nil {
		slog.Error("error while building disgo instance", slog.Any("error", err))
		return
	}
	defer client.Close(context.TODO())

	if _, err = client.Rest.SetGuildCommands(client.ApplicationID, cfg.GuildID, commands); err != nil {
		slog.Error("error while registering commands", slog.Any("error", err))
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while opening disgo gateway", slog.Any("error", err))
	}

	slog.Info("bot is now running.  Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == "say" {
		err := event.CreateMessage(discord.NewMessageCreate().WithContent(data.String("message")).WithEphemeral(data.Bool("ephemeral")))
		if err != nil {
			slog.Error("error on sending response", slog.Any("error", err))
		}
	}
}

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
