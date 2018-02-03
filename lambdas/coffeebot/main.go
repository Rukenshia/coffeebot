package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Rukenshia/coffeebot/lambdas/coffeebot/rocketchat"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

// Handler handles the Amazon API Gateway Event
func Handler(event events.CloudWatchEvent) (interface{}, error) {
	rlog := log.WithField("EventId", event.ID).WithField("RocketChatUrl", os.Getenv("RocketChatUrl"))
	rlog.Infof("Handling CloudWatch event")

	username, err := DecryptEnvVar("", "RocketChatUsername")
	if err != nil {
		rlog.Errorf("Decrypting RocketChatUsername failed: %v", err)
		return nil, err
	}
	password, err := DecryptEnvVar("", "RocketChatPassword")
	if err != nil {
		rlog.Errorf("Decrypting RocketChatPassword failed: %v", err)
		return nil, err
	}

	rc := rocketchat.NewClient(os.Getenv("RocketChatUrl"))

	if err := rc.Login(username, password); err != nil {
		rlog.Error(err)
		return nil, err
	}

	users, err := rc.Users.List()
	if err != nil {
		rlog.Error(err)
		return nil, err
	}
	users = filter(users, func(u rocketchat.User) bool {
		botIdx := strings.Index(u.Username, "bot")
		return u.Active == true && botIdx != 0 && botIdx != len(u.Username)-3
	})

	if len(users) == 0 {
		rlog.Info("Not assigning any pairings, no users")
		return nil, nil
	}

	if len(users)%2 != 0 {
		rlog.Infof("%s is skipped this time due to user count", users[0].Username)
		users = users[1:]
	}

	indexes := []int{}
	for idx := range users {
		indexes = append(indexes, idx)
	}
	pairings := shuffle(indexes)
	if len(pairings)%2 != 0 {
		pairings = pairings[1:]
	}

	pa := pairings[0 : len(pairings)/2]
	pb := pairings[len(pairings)/2:]
	for idx, first := range pa {
		firstUser := users[first]
		secondUser := users[pb[idx]]

		if firstUser.Username == "jan" || secondUser.Username == "jan" {
			SendCoffeeInvitation(rc, firstUser, secondUser)
		}
	}

	return nil, nil
}

func main() {
	log.SetLevel(log.DebugLevel)
	lambda.Start(Handler)
}

func shuffle(vals []int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func filter(vs []rocketchat.User, f func(rocketchat.User) bool) []rocketchat.User {
	vsf := make([]rocketchat.User, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
