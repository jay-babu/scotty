package main

import (
	"net/http"
	"time"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jay-babu/scotty/log"
	"github.com/jay-babu/scotty/openapi"
	oapi "github.com/sashabaranov/go-openai"
)

type ServerImpl struct{}

// ScottyChat implements openapi.ServerInterface.
func (ServerImpl) ScottyChat(c *gin.Context) {
	input := openapi.ScottyChatInput{}
	err := c.ShouldBind(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	resp, err := openapi.GptClient.CreateChatCompletion(c, oapi.ChatCompletionRequest{
		Model: oapi.GPT3Dot5Turbo,
		Messages: []oapi.ChatCompletionMessage{
			{
				Role:    oapi.ChatMessageRoleUser,
				Content: input.Message,
			},
		},
		User: c.ClientIP(),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	messages := make([]openapi.Message, 0)
	for _, message := range resp.Choices {
		messages = append(messages, openapi.Message{Message: message.Message.Content})
	}

	log.SLogger.Info(resp.Choices[0].Message.Content)
	output := openapi.ScottyChatOutput{
		Messages: messages,
	}
	c.JSON(http.StatusOK, output)
}

var _ openapi.ServerInterface = (*ServerImpl)(nil)

func main() {
	r := gin.New()
	r.Use(requestid.New())

	r.Use(ginzap.Ginzap(log.Logger, time.RFC3339, true))
	r.Use(cors.Default())
	r.Use(ginzap.RecoveryWithZap(log.Logger, true))

	swagger, err := openapi.GetSwagger()
	if err != nil {
		// This should never error
		panic("there was an error getting the swagger")
	}

	// Clear out the servers array in the swagger spec. It is recommended to do this so that it skips validating
	// that server names match.
	swagger.Servers = nil

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Use(middleware.OapiRequestValidator(swagger))

	var myAPI ServerImpl

	openapi.RegisterHandlers(r, myAPI)
	r.Run()
}
