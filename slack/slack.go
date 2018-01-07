package messages

import (
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

// Listen receives Slack events and triggers actions.
func Listen() chan []string {
	// Retrieve a Slack API token from our environment variables
	s := make(chan []string)
	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	go func(chan []string) {
	Loop:
		for {
			select {
			case msg := <-rtm.IncomingEvents:
				// fmt.Println("Event Received: ")
				switch ev := msg.Data.(type) {
				case *slack.ConnectedEvent:
					// fmt.Println("Connection counter:", ev.ConnectionCount)

				case *slack.MessageEvent:
					// fmt.Printf("Message: %v\n", ev)
					info := rtm.GetInfo()
					botuser := info.User.ID
					// Set a prefix that should be met in order to warrant a response from us
					prefix := fmt.Sprintf("<@%s> ", botuser)
					// If the original message wasn’t posted by our bot AND it contains our
					// prefix @botuser, then we’ll respond to the channel.
					// if ev.User != botuser && strings.HasPrefix(ev.Text, prefix) {
					//	rtm.SendMessage(rtm.NewOutgoingMessage("What's up buddy!?!?", ev.Channel))
					// }
					if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
						respond(rtm, ev, prefix, s)
					}

				case *slack.RTMError:
					fmt.Printf("Error: %s\n", ev.Error())

				case *slack.InvalidAuthEvent:
					fmt.Println("Invalid credentials")
					close(s)
					break Loop

				default:
					//Take no action
				}
			}
		}
	}(s)
	return s
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string, ch chan []string) {
	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	acceptedGreetings := map[string]bool{
		"what's up?": true,
		"hey!":       true,
		"yo":         true,
	}
	acceptedHowAreYou := map[string]bool{
		"how's it going?": true,
		"how are ya?":     true,
		"feeling okay?":   true,
	}

	if acceptedGreetings[text] {
		response = "What's up buddy!?"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else if acceptedHowAreYou[text] {
		response = "Good. How are you?"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else if strings.HasPrefix(text, "route ") {
		text = strings.TrimPrefix(text, "route ")
		sl := strings.Split(text, " ")
		if len(sl) < 2 {
			response = "Incorrect route"
			rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
			return
		}
		response = fmt.Sprintf("Configuring route to %s via %s", sl[0], sl[1])
		ch <- sl
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else {
		response = "What!?... "
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}
