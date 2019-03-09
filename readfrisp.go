package main
/*
*	"crypto/rsa"
*	"github.com/aws/aws-sdk-go/aws"
*/
import (
	"fmt"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)
/**
* env is path in the s3 bucket or what frisp calls a vault
*/
type Event struct {
	bucket string `json:"bucket"`
	env string `json:"environment"`
	region string `json:"region"`
	filename string `json:"filename"`
}

type KeyResponse struct {
	ciper string `json:"cipher:"`
	plain string `json:"plain:"`
	//key,cert,status?
}

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}
/*
 encrypt data with a public key (defaults to RSAES PKCS#1 v1.5)
var encrypted = publicKey.encrypt(bytes);
*/
func HandleRequest(ctx context.Context, config Event) (string, error) {
	/*
	* get frisp privatekey from s3
	* attempt to read environment encrypted data (SSL cert and key)
	* use frisp private key to decrypt environment SSL cert and key
	* print/dump result to caller as json
	*/
	sess := session.Must(session.NewSession())

	svc := s3.New(sess)

	i := 0
	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: &config.region,
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		fmt.Println("Page,", i)
		i++

		for _, obj := range p.Contents {
			fmt.Println("Object:", *obj.Key)
		}
		return true
	})
	if err != nil {
		return "failed to list objects", err
	}
	return fmt.Sprintf("Hello %s!", config.region ), nil
}

func main() {
	lambda.Start(HandleRequest)
}