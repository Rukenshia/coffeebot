# Coffeebot

This bot helps users to find someone to have a cup of coffee with via RocketChat. It encourages social activity within a group of people, e.g. coworkers. All users of a server are paired every Monday and get another person proposed. 

## Lambda

The lambda is written in Go, hence Go1.9 is required. The infrastructure needs to be deployed on AWS.

## Deployment

Coffeebot is deployed using [sceptre](https://sceptre.cloudreach.com/latest/). You can use `pip install -r requirements.txt` to install sceptre. It is recommended to use a virtualenv.

Use `./configure` to configure coffeebot. It will ask you for the Rocket.Chat URL, Username and Password.

Afterwards, the whole stack can be deployed:
```shell
$ make deploy
```
