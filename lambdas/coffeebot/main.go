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

func init() {
	level := os.Getenv("LogLevel")
	switch level {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}
}

// Handler handles the Amazon API Gateway Event
func Handler(event events.CloudWatchEvent) (interface{}, error) {
	silentMode := os.Getenv("SilentMode") == "true"

	rlog := log.WithField("EventId", event.ID).WithField("RocketChatUrl", os.Getenv("RocketChatUrl"))
	rlog.WithField("silent", silentMode).Infof("Handling CloudWatch event")

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

	blacklist := strings.Split(os.Getenv("RocketChatBlacklist"), ",")

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
		// Remove users that are blacklisted
		for _, b := range blacklist {
			if u.Username == b {
				return false
			}
		}
		return true
	})
	users = filter(users, func(u rocketchat.User) bool {
		botIdx := strings.Index(u.Username, "bot")
		return u.Active == true && botIdx != 0 && botIdx != len(u.Username)-3
	})
	users = filter(users, func(u rocketchat.User) bool {
		if u.Username == "" {
			rlog.WithField("user", u).Warn("Found empty username")
			return false
		}
		return true
	})

	var usernames []string
	for _, user := range users {
		usernames = append(usernames, user.Username)
	}
	rlog.WithField("users", usernames).Debugf("Determined eligible users")

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

		if silentMode {
			rlog.Debugf("Would send message to %s and %s now", firstUser.Username, secondUser.Username)
		} else {
			rlog.Debugf("Sending messages to %s and %s", firstUser.Username, secondUser.Username)
			SendCoffeeInvitation(rc, firstUser, secondUser)
		}
	}

	return nil, nil
}

func main() {
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
