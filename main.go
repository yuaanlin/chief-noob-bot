package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {

	file, err := os.ReadFile("/etc/resolv.conf")
	if err == nil {
		println(string(file))
	}

	setRoleChannelID := "999922492422504488"
	setRoleMsgID := "999924590820196482"

	roles := map[string]string{
		"meow_b":         "999923150861115402",
		"4459_ComfyBlob": "999923603938218025",
	}

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Printf("Error while starting bot: %s", err)
		return
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessageReactions

	content := "大家好，我是 Chief Noob 機器小雞，請給這個訊息加入 Reaction 來自動設定你在這個 Server 的身分組：\n"
	content += "\n點擊 <:meow_b:968378576162422784> 可以把自己設定為「前端工程師」"
	content += "\n點擊 <:4459_ComfyBlob:968822210263404584> 可以把自己設定為「後端工程師」"
	content += "\n\n未來會繼續追加更多身分組的選項，請大家踴躍提供 emoji 和身分組的創意哦！"

	_, err = dg.ChannelMessageEdit(setRoleChannelID, setRoleMsgID, content)
	if err != nil {
		println(err.Error())
		return
	}

	dg.AddHandler(
		func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
			if r.UserID == s.State.User.ID {
				return
			}
			if r.ChannelID == "999922492422504488" && r.MessageID == "999924590820196482" {
				for emoji, role := range roles {
					if r.Emoji.Name == emoji {
						err := s.GuildMemberRoleAdd(r.GuildID, r.UserID, role)
						if err != nil {
							println(err.Error())
							return
						}
					}
				}
			}
		},
	)

	dg.AddHandler(
		func(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
			if r.UserID == s.State.User.ID {
				return
			}

			if r.ChannelID == "999922492422504488" && r.MessageID == "999924590820196482" {
				for emoji, role := range roles {
					if r.Emoji.Name == emoji {
						err := s.GuildMemberRoleRemove(
							r.GuildID, r.UserID, role,
						)
						if err != nil {
							println(err.Error())
							return
						}
					}
				}
			}
		},
	)

	// Connect to the gateway
	err = dg.Open()
	if err != nil {
		fmt.Printf("Error while connecting to gateway: %s", err)
		return
	}

	// Wait until Ctrl+C or another signal is received
	fmt.Println("The bot is now running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close the Discord session
	dg.Close()
}
