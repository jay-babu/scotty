package openapi

import (
	"os"

	oapi "github.com/sashabaranov/go-openai"
)

var GptClient *oapi.Client

func GptToken() (token string) {
	if stage, ok := os.LookupEnv("GPT_TOKEN"); ok {
		return stage
	}
	panic("GPT_TOKEN not found")
}

func init() {
	GptClient = oapi.NewClient(GptToken())
}
