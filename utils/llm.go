package utils

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func openaiRequester(base_url, apiKey, model string, messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam) (*openai.ChatCompletion, error) {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(base_url),
	)

	req := openai.ChatCompletionNewParams{
		Model:    model,
		Messages: messages,
	}

	if tools != nil {
		req.Tools = tools
	}

	ctx := context.Background()
	completion, err := client.Chat.Completions.New(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	return completion, nil
}

func DeepSeekChat(messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam) (*openai.ChatCompletion, error) {
	return openaiRequester("https://api.deepseek.com/v1", "sk-775b4e8eebb4435f88ffc1fdcb37c7f6", "deepseek-chat", messages, tools)
}
