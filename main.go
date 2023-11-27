package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {

	// Load Env variables from .dot file
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	// Create a new client to slack by giving token
	// Set debug to true while developing
	// Also add a ApplicationToken option to the client
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
	socket := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socket *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socket.Events:
				// We have a new Events, let's type switch the event
				// Add more use cases here if you want to listen to other events.
				switch event.Type {
				// handle EventAPI events
				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socket.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch

					//log.Println(eventsAPI) // commenting for event hanndling

					//------------------------------------
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := HandleEventMessage(eventsAPI, client)
					if err != nil {
						// Replace with actual err handeling
						log.Fatal(err)
					}
				}
			}
		}
	}(ctx, client, socket)

	socket.Run()
}

// HandleEventMessage will take an event and handle it properly based on the type of event
func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := HandleAppMentionEventToBot(ev, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

// HandleAppMentionEventToBot is used to take care of the AppMentionEvent when the bot is mentioned
// HandleAppMentionEventToBot is used to take care of the AppMentionEvent when the bot is mentioned
func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client) error {
	// Grab the user name based on the ID of the one who mentioned the bot
	_, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	// Check if the user provided a command and two numbers
	text := strings.ToLower(event.Text)
	words := strings.Fields(text)
	fmt.Println("longitud", len(words))
	fmt.Println("words", words)
	// eliminar la primera palabra
	words = words[1:]
	fmt.Println("words", words)
	fmt.Println("name of user", event.Text)

	if words[0] == "check_bipartite" {
		G := construirGrafo(words[1:])
		fmt.Println("Grafo", G)
		if G == nil {
			respuesta := slack.Attachment{
				Text:       "Ingrese un grafo válido",
				Color:      "#ff0000",
				MarkdownIn: []string{"text"},
			}
			_, _, err := client.PostMessage(event.Channel, slack.MsgOptionAttachments(respuesta))
			if err != nil {
				return err
			}
		}
		if isBipartite(G) {
			respuesta := slack.Attachment{
				Text:       "El grafo es bipartito maestro",
				Color:      "#00ff00",
				MarkdownIn: []string{"text"},
			}
			_, _, err := client.PostMessage(event.Channel, slack.MsgOptionAttachments(respuesta))
			if err != nil {
				return err
			}
		} else {
			respuesta := slack.Attachment{
				Text:       "El grafo no es bipartito maestro",
				Color:      "#ff0000",
				MarkdownIn: []string{"text"},
			}
			_, _, err := client.PostMessage(event.Channel, slack.MsgOptionAttachments(respuesta))
			if err != nil {
				return err
			}
		}

		// respuesta := slack.Attachment{
		// 	Text:       "El grafo es bipartito",
		// 	Color:      "#00ff00",
		// 	MarkdownIn: []string{"text"},
		// }
		// _, _, err := client.PostMessage(event.Channel, slack.MsgOptionAttachments(respuesta))
		// if err != nil {
		// 	return err
		// }

	}

	return nil
}

func aristar(G [][]int, a int, b int) {
	G[a] = append(G[a], b)
	G[b] = append(G[b], a)
}

func construirGrafo(words []string) [][]int {
	// cambiar el primer elemento a int

	len := len(words)
	if len%2 != 0 {
		return nil
	}

	n, err := strconv.Atoi(words[0])
	if err != nil {
		return nil
	}
	// crear matriz de adyacencia
	G := make([][]int, n+1)
	A := make([]int, 0)
	B := make([]int, 0)
	for i := 2; i < len; i += 2 {
		a, err := strconv.Atoi(words[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		b, err := strconv.Atoi(words[i+1])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		A = append(A, a)
		B = append(B, b)
		aristar(G, a, b)
	}
	return G
}

func isBipartite(G [][]int) bool {
	// arreglo de colores
	n := len(G)
	color := make([]int, n)
	for i := 0; i < n; i++ {
		color[i] = -1
	}
	fmt.Println("longitud de color", n)
	for i := 1; i < n; i++ {
		if color[i] == -1 {
			if !dfs(G, i, color, 0) {
				fmt.Println("no es bipartito", i)
				return false
			}
		}
	}
	return true

}
func dfs(G [][]int, v int, color []int, curr int) bool {
	fmt.Println("entro  ", v, " con color ", curr)
	color[v] = curr
	n := len(G[v])
	for i := 0; i < n; i++ {
		w := G[v][i]
		if color[w] == -1 {
			if !dfs(G, w, color, 1-curr) {
				return false
			}
		} else {
			if color[w] != 1-curr {
				fmt.Println("no es bipartito entre ", v, " y ", w)
				fmt.Println("color de ", v, " es ", color[v])
				fmt.Println("color de ", w, " es ", color[w])
				return false
			}
		}
	}

	return true
}

// go run main.go
