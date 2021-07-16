package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Movie struct {
	ID   string
	Name string
}

func findAll(request events.APIGatewayProxyResponse) (events.APIGatewayProxyResponse, error) {

	size, _ := strconv.Atoi(request.Headers["size"])

	config := &aws.Config{
		Region: aws.String("ap-southeast-1"),
	}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)
	req, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Limit:     aws.Int64(int64(size)),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while scanning DynamoDB",
		}, nil

	}

	var movies []Movie

	err = dynamodbattribute.UnmarshalListOfMaps(req.Items, &movies)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while unmarshal result",
		}, nil
	}

	response, err := json.Marshal(movies)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding to string value",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil

}

func main() {
	lambda.Start(findAll)
}
