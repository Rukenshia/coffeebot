STAGE=dev

kms:
	sceptre launch-stack $(STAGE) kms

deploy:
	sceptre --var-file config/creds.yaml launch-env $(STAGE)

.PHONY: kms deploy