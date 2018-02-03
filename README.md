# Coffeebot

This bot pairs all Rocket.Chat users of a server every Monday for a cup of coffee, lunch or whatever you feel like.

## Lambda

The lambda is written in Go, hence Go1.9 is required.

## Deployment

Coffeebot is deployed using [sceptre](https://sceptre.cloudreach.com/latest/). You can use `pip install -r requirements.txt` to install sceptre. It is recommended to use a virtualenv.

Use `./configure` to configure coffeebot. It will ask you for the Rocket.Chat URL, Username and Password.

Afterwards, the whole stack can be deployed:
```shell
$ make deploy
```