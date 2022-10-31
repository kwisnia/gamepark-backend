package awsconf

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
)

var AWSSession *session.Session

func ConnectAws() {
	AccessKeyID := config.GetEnv("AWS_ACCESS_KEY_ID")
	SecretAccessKey := config.GetEnv("AWS_SECRET_ACCESS_KEY")
	MyRegion := config.GetEnv("AWS_REGION")
	fmt.Println("region", MyRegion)
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	AWSSession = sess
}
