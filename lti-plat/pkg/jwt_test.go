package pkg

import (
	"fmt"
	"testing"
)

func TestIdToken(t *testing.T) {
	token := IdToken("c1", "abc", "123456", "r1")
	fmt.Println(string(decryptToken(token).Payload))
}
