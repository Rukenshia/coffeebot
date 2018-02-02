# Coffeebot

This bot pairs all Rocket.Chat users of a server every Monday for a cup of coffee, lunch or whatever you feel like.

## Lambda

The lambda is written in Go, hence Go1.9 is required.

## Deployment

Coffeebot is deployed using [sceptre](https://sceptre.cloudreach.com/latest/). You can use `pip install -r requirements.txt` to install sceptre. It is recommended to use a virtualenv.

The first required step is to create the KMS stack:
```bash
$ sceptre launch-stack dev kms
```

Then, use the created key to encrypt your Rocket.Chat username and password:
```bash
$ aws kms encrypt --key-id e318d9e7-f1dc-4bfa-b904-c10346a2abd6 --plaintext 'myUsername'
$ aws kms encrypt --key-id e318d9e7-f1dc-4bfa-b904-c10346a2abd6 --plaintext 'myPassword'
```

To deploy Coffeebot, create a `config/creds.yml`:
```yml
rocket_chat:
  url: https://my.rocket.chat
  username: 'KMSEncryptedUsername'
  password: 'KMSEncryptedPassword'
```

Afterwards, the whole stack can be deployed:
```shell
$ sceptre --var-file config/creds.yml launch-env dev
```