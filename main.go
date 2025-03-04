package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/drummer3333/go-smtpsrv"

	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
)

func main() {
	gotifyURL, _ := url.Parse(*flagGotifyURL)
	gotifyClient := gotify.NewClient(gotifyURL, &http.Client{})
	versionResponse, err := gotifyClient.Version.GetVersion(nil)

	if err != nil {
		log.Fatal("Could not request version ", err)
		return
	}
	version := versionResponse.Payload
	log.Println("Found Gotify version", *version)

	cfg := smtpsrv.ServerConfig{
		ReadTimeout:       time.Duration(*flagReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(*flagWriteTimeout) * time.Second,
		ListenAddr:        *flagListenAddr,
		MaxMessageBytes:   int(*flagMaxMessageSize),
		BannerDomain:      *flagServerName,
		AuthDisabled:      false,
		AllowInsecureAuth: true,
		Auther: smtpsrv.AuthFunc(func(username, password string) error {
			return nil
		}),
		Handler: smtpsrv.HandlerFunc(func(c *smtpsrv.Context) error {
			msg, err := c.Parse()
			if err != nil {
				log.Fatalf("Cannot read your message: " + err.Error())
				return nil
			}

			user, password, err := c.User()
			if err != nil {
				log.Fatalf("Could not get password and user: %v", err)
				return nil
			}

			params := message.NewCreateMessageParams()
			params.Body = &models.MessageExternal{
				Title:    msg.Subject,
				Message:  concatMsgText(msg, c),
				Priority: getPriority(user),
			}
			_, err = gotifyClient.Message.CreateMessage(params, auth.TokenAuth(password))

			if err != nil {
				log.Fatalf("Could not send message %v", err)
				return nil
			}
			log.Println("Message Sent!")

			return nil
		}),
	}

	fmt.Println(smtpsrv.ListenAndServe(&cfg))
}

func getPriority(user string) int {
	parts := strings.Split(user, "-")

	prio, err := strconv.Atoi(parts[len(parts)-1])
	if err == nil {
		return prio
	}

	return 0
}

func concatMsgText(msg *smtpsrv.Email, c *smtpsrv.Context) string {
	var text strings.Builder

	text.WriteString(string(msg.TextBody))
	text.WriteString("\n\n========== Adresses ==========\n")

	text.WriteString("from: ")
	text.WriteString(c.From().String())
	text.WriteString("\n")

	text.WriteString("to: ")
	text.WriteString(c.To().String())
	text.WriteString("\n\n\n")

	if len(msg.Attachments) > 0 {
		text.WriteString("\n\n========== Attatchments ==========\n")
	}

	for _, a := range msg.Attachments {
		data, _ := io.ReadAll(a.Data)
		if data != nil {
			text.WriteString("\n\n---------- ")
			text.WriteString(a.Filename)
			text.WriteString(" (")
			text.WriteString(a.ContentType)
			text.WriteString(")")
			text.WriteString(" ----------\n\n")
			text.WriteString(base64.StdEncoding.EncodeToString(data))
		}
	}

	msgText := text.String()
	return msgText
}
