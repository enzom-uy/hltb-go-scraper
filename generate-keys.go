package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateAPIKey(prefix string) string {
	bytes := make([]byte, 16) // 32 caracteres hex
	rand.Read(bytes)
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(bytes))
}

func main() {
	fmt.Println("Generating API keys:")
	fmt.Println()

	botKey := generateAPIKey("bot")
	websiteKey := generateAPIKey("website")

	fmt.Printf("Discord Bot APi Key: %s\n", botKey)
	fmt.Printf("Website APi Key: %s\n", websiteKey)

}
