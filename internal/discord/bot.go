package discord

import (
	auth "github.com/Floor-Gang/authclient"
	"github.com/Floor-Gang/init-discord-bot/internal"
	util "github.com/Floor-Gang/utilpkg"
	dg "github.com/bwmarrin/discordgo"
	"log"
)

// Bot structure
type Bot struct {
	Auth   *auth.AuthClient
	Client *dg.Session
	Config *internal.Config
}

// Start starts discord client, configuration and database
func Start() {
	var err error

	// get Config.yml
	config := internal.GetConfig()

	// setup authentication server
	// you can use this to get the bot's access token
	// and authenticate each user using a command.
	authClient, err := auth.GetClient(config.Auth)

	if err != nil {
		log.Fatalln("Failed to connect to authentication server", err)
	}

	register, err := authClient.Register(
		auth.Feature{
			Name:        "", // Give this bot / feature a name
			Description: "", // Describe what this bot is doing
			Commands: []auth.SubCommand{ // list all the commands this bot / feature has
				{
					Name:        "",           // Command name like "add"
					Description: "",           // Describe what the command does
					Example:     []string{""}, // [command name, argument 1, argument 2] like [add, #channel, #channel]
				},
			},
			CommandPrefix: config.Prefix,
		},
	)

	if err != nil {
		log.Fatalln("Failed to register with authentication server", err)
	}

	client, err := dg.New(register.Token)

	if err != nil {
		panic(err)
	}

	bot := Bot{
		Auth:   &authClient,
		Client: client,
		Config: &config,
	}

	client.AddHandler(bot.onReady)
	client.AddHandler(bot.onMessage)

	if err = client.Open(); err != nil {
		util.Report("Was an authentication token provided?", err)
	}
}
