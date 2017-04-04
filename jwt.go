package box

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
)

// Box app information, like ClientID, sub id, etc.
var (
	privateKey []byte

	ClientID     = os.Getenv("BOX_CLIENT_ID")
	ClientSecret = os.Getenv("BOX_CLIENT_SECRET")
	SubID        = os.Getenv("BOX_SUB_ID")
	SubType      = os.Getenv("BOX_SUB_TYPE")

	JWTExpiration = 30 * time.Second
)

func init() {
	_ = SetRSAPrivateKeyFile(os.Getenv("BOX_PRIVATE_KEY_FILE"))
}

// SetRSAPrivateKeyFile reads the given private key file into memory
func SetRSAPrivateKeyFile(fn string) (err error) {
	privateKey, err = ioutil.ReadFile(fn)
	return
}

// SetRSAPrivateKey sets the private key to the contents of the supplied byte slice
func SetRSAPrivateKey(key []byte) {
	privateKey = key
}

// JWTToken returns a JWT token string suitable for use in as an auth header.
// The private key file is read from  BOX_PRIVATE_KEY_FILE or using SetRSAPrivateKeyFile
func JWTToken() (tokenString string) {

	if privateKey == nil || len(privateKey) < 1 {
		return
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"iss":          ClientID,
		"sub":          SubID,
		"box_sub_type": SubType,
		"aud":          "https://api.box.com/oauth2/token",
		"jti":          uuid.New(),
		"exp":          now.Add(JWTExpiration).Unix(),
	})

	fmt.Println(now.Add(JWTExpiration).Unix())

	if tokenString, err = token.SignedString(key); err != nil {
		panic(err)
	}

	return
}
