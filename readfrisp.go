package readfrisp

import (
	"fmt"
	"crypto/rsa"
)
/*
 encrypt data with a public key (defaults to RSAES PKCS#1 v1.5)
var encrypted = publicKey.encrypt(bytes);
*/
func main() {
	fmt.Println("Hello, 世界")
	/*
	* get frisp privatekey from s3
	* attempt to read environment encrypted data (SSL cert and key)
	* use frisp private key to decrypt environment SSL cert and key
	* print/dump result to caller as json
	*/
}
