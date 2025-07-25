#!/bin/bash
set -e

REGION="us-east-1"
ENDPOINT_URL="http://localhost:4566"

GATEWAY_ID=$(aws apigateway get-rest-apis --endpoint-url $ENDPOINT_URL --region $REGION | jq -r '.items[0].id')
echo "Using API Gateway ID: $GATEWAY_ID"

curl --request POST \
  --url $ENDPOINT_URL/restapis/$GATEWAY_ID/dev/_user_request_/callbacks-trackingagrovarejo/v1/eventos_operacao \
  --header 'content-type: application/json' \
  --data '{
  "foo": "bar"
}'

awslocal sqs receive-message --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/my-queue --region $REGION
