#!/bin/bash
set -e

REGION="us-east-1"

pushd ../app/cmd
GOOS=linux GOARCH=amd64 go build -o ../main
zip ../function.zip ../main
popd

SNS_TOPIC_ARN=$(aws sns create-topic --name my-topic --endpoint-url http://localhost:4566 --region $REGION | jq -r .TopicArn)
echo "SNS topic ARN: $SNS_TOPIC_ARN"

SQS_QUEUE_URL=$(aws sqs create-queue --queue-name my-queue --endpoint-url http://localhost:4566 --region $REGION | jq -r .QueueUrl)
echo "SQS queue URL: $SQS_QUEUE_URL"

SQS_QUEUE_ARN=$(aws sqs get-queue-attributes --queue-url $SQS_QUEUE_URL --attribute-names QueueArn --endpoint-url http://localhost:4566 --region $REGION | jq -r .Attributes.QueueArn)
echo "SQS queue ARN: $SQS_QUEUE_ARN"

aws sns subscribe \
  --topic-arn $SNS_TOPIC_ARN \
  --protocol sqs \
  --notification-endpoint $SQS_QUEUE_ARN \
  --endpoint-url http://localhost:4566 \
  --region $REGION

aws lambda create-function \
  --function-name my-lambda-function \
  --runtime go1.x \
  --handler main \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --zip-file fileb://../app/function.zip \
  --endpoint-url http://localhost:4566 \
  --region $REGION

API_ID=$(aws apigateway create-rest-api --name "MyAPI" --endpoint-url http://localhost:4566 --region $REGION | jq -r .id)
PARENT_RESOURCE_ID=$(aws apigateway get-resources --rest-api-id $API_ID --endpoint-url http://localhost:4566 --region $REGION | jq -r .items[0].id)
RESOURCE_ID=$(aws apigateway create-resource --rest-api-id $API_ID --parent-id $PARENT_RESOURCE_ID --path-part callbacks-trackingagrovarejo --endpoint-url http://localhost:4566 --region $REGION | jq -r .id)
RESOURCE_ID2=$(aws apigateway create-resource --rest-api-id $API_ID --parent-id $RESOURCE_ID --path-part v1 --endpoint-url http://localhost:4566 --region $REGION | jq -r .id)
RESOURCE_ID3=$(aws apigateway create-resource --rest-api-id $API_ID --parent-id $RESOURCE_ID2 --path-part eventos_operacao --endpoint-url http://localhost:4566 --region $REGION | jq -r .id)

aws apigateway put-method --rest-api-id $API_ID --resource-id $RESOURCE_ID3 --http-method POST --authorization-type "NONE" --endpoint-url http://localhost:4566 --region $REGION
aws apigateway put-integration --rest-api-id $API_ID --resource-id $RESOURCE_ID3 --http-method POST --type AWS_PROXY --integration-http-method POST --uri arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$REGION:000000000000:function:my-lambda-function/invocations --endpoint-url http://localhost:4566 --region $REGION
aws lambda add-permission --function-name my-lambda-function --statement-id apigateway-test-2 --action lambda:InvokeFunction --principal apigateway.amazonaws.com --source-arn arn:aws:execute-api:$REGION:000000000000:$API_ID/*/POST/callbacks-trackingagrovarejo/v1/eventos_operacao --endpoint-url http://localhost:4566 --region $REGION
aws apigateway create-deployment --rest-api-id $API_ID --stage-name dev --endpoint-url http://localhost:4566 --region $REGION

echo "API endpoint: http://localhost:4566/restapis/$API_ID/dev/_user_request_/callbacks-trackingagrovarejo/v1/eventos_operacao"
