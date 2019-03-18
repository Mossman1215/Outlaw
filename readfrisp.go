package main

/*
*	"crypto/rsa"
*	"github.com/aws/aws-sdk-go/aws"
 */
import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/**
* Event contains input data for the lambda
* env is folder in the s3 bucket or what frisp calls a vault
 */
type event struct {
	Bucket   string `json:"bucket"`
	Env      string `json:"environment"`
	Region   string `json:"region"`
	Filename string `json:"filename"`
}

type keyResponse struct {
	Ciper string `json:"cipher:"`
	Plain string `json:"plain:"`
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
func HandleRequest(ctx context.Context, config event) (string, error) {
	/*
	* get frisp privatekey from s3
	* attempt to read environment encrypted data (SSL cert and key)
	* use frisp private key to decrypt environment SSL cert and key
	* print/dump result to caller as json
	 */
	bucket := config.Bucket
	env := config.Env
	//filename := config.Filename
	item := "_private/" + env + ".enc"
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(config.Region)},
	)

	downloader := s3manager.NewDownloader(sess)
	encValue := aws.NewWriteAtBuffer([]byte{})
	requestInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	}
	numBytes, err := downloader.Download(encValue,
		requestInput)
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}

	fmt.Println("Downloaded", len(encValue.Bytes()), numBytes, "bytes")
	encryptedStr := fmt.Sprintf("%q", encValue.Bytes())

	//read private key into memory
	//privKey := aws.NewWriteAtBuffer([]byte{})
	//privKeyRequestInput := &s3.GetObjectInput{
	//	Bucket: aws.String(bucket),
	//	Key:    aws.String("_private/" + env + ".inc"),
	//}
	//numBytespriv, err := downloader.Download(privKey,
	//	privKeyRequestInput)
	//if err != nil {
	//	exitErrorf("Unable to download item %q, %v", item, err)
	//}
	//create response struct and assign existing encrypted key to it

	response := &keyResponse{encryptedStr, ""}
	responsestr, err := json.Marshal(response)
	if err != nil {
		exitErrorf("failure to encode json", err)
	}
	return string(responsestr), nil
}

func main() {
	lambda.Start(HandleRequest)
}
