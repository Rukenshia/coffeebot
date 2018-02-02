BINARY_NAME=main
OUTPUT_DIR=dist
DEPLOYMENT_BUCKET=coffeebot-artifacts
LAMBDA_PATH=coffeebot/lambda.zip
BASE_URL=https://coffeebot.in.fkn.space

all: clean build zip

artifact: zip
	aws s3 cp $(OUTPUT_DIR)/lambda.zip s3://$(DEPLOYMENT_BUCKET)/$(LAMBDA_PATH)

upload: artifact
	$(eval function_name := `aws lambda list-functions | jq -r '.Functions[] | select(.FunctionName | contains("Coffeebot")).FunctionName'`)

	aws lambda update-function-code \
		--function-name $(function_name) \
		--s3-bucket $(DEPLOYMENT_BUCKET) \
		--s3-key $(LAMBDA_PATH) \
		--publish

invoke:
	$(eval function_name := `aws lambda list-functions | jq -r '.Functions[] | select(.FunctionName | contains("Coffeebot")).FunctionName'`)
	$(eval payload := `cat resources/invoke.json`)

	aws lambda invoke \
		--function-name $(function_name) \
		--log-type Tail \
		--payload "$(payload)" \
		/dev/null | jq -r '.LogResult' | base64 --decode

upload_invoke: upload invoke


clean:
	rm -rf $(OUTPUT_DIR)

build:
	mkdir -p $(OUTPUT_DIR)
	GOOS=linux go build -o $(OUTPUT_DIR)/$(BINARY_NAME)

zip: build
	cd $(OUTPUT_DIR) && zip lambda.zip $(BINARY_NAME)
