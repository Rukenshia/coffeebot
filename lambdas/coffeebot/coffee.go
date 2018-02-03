package main

import (
	"bytes"
	"log"
	"math/rand"
	"text/template"

	"github.com/Rukenshia/coffeebot/lambdas/coffeebot/rocketchat"
)

var messages = []string{
	"Good morning, @{{.First.Username}}! Don't you think this is an awesome week to get to know @{{.Second.Username}} a little bit better? I think you'd be great friends. Why don't you contact @{{.Second.Username}} and go grab a coffee together this week?",
	"Hey @{{.First.Username}}! Even if you don't like coffee, why don't you meet up with @{{.Second.Username}} this week? @{{.Second.Username}} will surely appreciate it!",
	"Coofffeeee time! For this week, your coffee date is @{{.Second.Username}}! If you're more of a cake person, why don't you go grab a piece of cake from Balzac down the road together with @{{.Second.Username}}?",
	"It's this time of the week again! @{{.Second.Username}} is dieing to learn more about you. Coffee/tea/water with @{{.Second.Username}} sounds great, doesn't it?",
	"I bet you haven't had enough of @{{.Second.Username}} yet. Why don't you two have some 1on1 time this week, go for a coffee or even lunch?",
}

// SendCoffeeInvitation sends a randomised, personalised message to the pair
func SendCoffeeInvitation(rc *rocketchat.Client, first, second rocketchat.User) error {
	m1 := generateMessage(first, second)
	m2 := generateMessage(second, first)

	if err := rc.Chat.PostMessage(rocketchat.Message{
		Channel: "@jan",
		Text:    m1,
		Emoji:   ":coffee:",
	}); err != nil {
		return err
	}

	return rc.Chat.PostMessage(rocketchat.Message{
		Channel: "@jan",
		Text:    m2,
		Emoji:   ":coffee:",
	})
}

func generateMessage(first, second rocketchat.User) string {
	msg := messages[rand.Intn(len(messages))]
	tpl, err := template.New("msg").Parse(msg)

	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, map[string]interface{}{
		"First":  first,
		"Second": second,
	}); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}
