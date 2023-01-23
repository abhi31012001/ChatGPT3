package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		log.Fatalln("API KEY is not present")
	}
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	const inputfile = "./input.txt"
	file, err := os.ReadFile(inputfile)
	if err != nil {
		log.Fatalf("failed to read: %v", err)
	}
	msgPrefix := "give me a short list of libraries that are used in the code \n ```java\n"
	msgSuffix := "\n```"
	msg := msgPrefix + string(file) + msgSuffix
	outputBilder := strings.Builder{}
	errr := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			msg,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(cr *gpt3.CompletionResponse) {
		outputBilder.WriteString(cr.Choices[0].Text)
	},
	)
	if errr != nil {
		log.Fatalln(errr)
	}
	output := strings.TrimSpace(outputBilder.String())
	const outputFile = "./output.txt"
	errs := os.WriteFile(outputFile, []byte(output), os.ModePerm)
	if errs != nil {
		log.Panic(errs)
	}
}
