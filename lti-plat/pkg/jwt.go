package pkg

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/jwt"
)

type LTIContext struct {
	Id    string   `json:"id"`
	Label string   `json:"label"`
	Title string   `json:"title"`
	Type  []string `json:"type"`
}

type LTIResourceLink struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type LTILaunchPresentation struct {
	DocumentTarget string `json:"document_target"`
	Height         int8   `json:"height"`
	Width          int8   `json:"width"`
	ReturnUrl      string `json:"return_url"`
}

type LTIClaims struct {
	jwt.Claims
	Nonce              string                `json:"nonce"`
	DeploymentId       string                `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id"`
	MessageType        string                `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
	Roles              []string              `json:"https://purl.imsglobal.org/spec/lti/claim/roles"`
	Context            LTIContext            `json:"https://purl.imsglobal.org/spec/lti/claim/context"`
	ResourceLink       LTIResourceLink       `json:"https://purl.imsglobal.org/spec/lti/claim/resource_link"`
	TargetLink         string                `json:"https://purl.imsglobal.org/spec/lti/claim/target_link_uri"`
	LaunchPresentation LTILaunchPresentation `json:"https://purl.imsglobal.org/spec/lti/claim/launch_presentation""`
	CustomClaim        map[string]string     `json:"https://purl.imsglobal.org/spec/lti/claim/custom"`
	Version            string                `json:"https://purl.imsglobal.org/spec/lti/claim/version"`
}

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	var err error
	privateKey, err = jwt.ParsePrivateKeyRSA([]byte(PrivateKey))
	if err != nil {
		fmt.Println(fmt.Sprintf("err when load private key:%v", err))
	}
	publicKey, err = jwt.ParsePublicKeyRSA([]byte(PublicKey))
	if err != nil {
		fmt.Println(fmt.Sprintf("err when load public key:%v", err))
	}
}

func IdToken(clientId, userId, nonce, resId string) string {
	claims := LTIClaims{
		Claims: jwt.Claims{
			IssuedAt: time.Now().Unix(),
			Expiry:   time.Now().Add(24 * time.Hour).Unix(),
			ID:       uuid.New().String(),
			Issuer:   ISSUER,
			Audience: jwt.Audience{clientId},
			Subject:  userId,
		},
		Nonce:        nonce,
		TargetLink:   "http://localhost:9000/launch",
		DeploymentId: "1",
		MessageType:  "LtiResourceLinkRequest",
		Version:      "1.3.0",
		CustomClaim: map[string]string{
			"Foo": "bar",
		},
		ResourceLink: LTIResourceLink{
			Id: resId,
		},
		LaunchPresentation: LTILaunchPresentation{
			DocumentTarget: "iframe",
		},
	}
	t, _ := jwt.Sign(jwt.RS256, privateKey, claims)
	return string(t)
}

func decryptToken(token string) *jwt.VerifiedToken {
	vt, err := jwt.Verify(jwt.RS256, publicKey, []byte(token))
	if err != nil {
		fmt.Printf("veriry token failed:%v", err)
		return nil
	}
	return vt
}
