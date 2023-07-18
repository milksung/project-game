package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

func CreateSign(data string, timeNow time.Time) string {

	timeStamp := fmt.Sprintf("%d", timeNow.Unix())
	apiKey := os.Getenv("AGENT_KEY")

	word := strings.ToLower(data) + timeStamp + strings.ToLower(apiKey)
	fmt.Println("word", word)
	hasher := sha256.New()

	if _, err := hasher.Write([]byte(word)); err != nil {
		return ""
	}

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}
