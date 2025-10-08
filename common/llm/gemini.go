package llm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"cookforyou.com/linebot-liff-template/common/models"

	"google.golang.org/genai"
)

type GoogleGemini interface {
	Chat(ctx context.Context, conversations []*models.Conversation) (string, error)
	io.Closer
}

type googleGemini struct {
	client  *genai.Client
	modelID string
}

func NewGoogleGemini(ctx context.Context, apiKey string, modelID string) (GoogleGemini, error) {
	retryClient := WrapClient(&http.Client{
		Timeout: 600 * time.Second,
	})
	retryClient.RetryMax = 5
	retryClient.RetryWaitMax = 10 * time.Second

	ai, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:     apiKey,
		Backend:    genai.BackendGeminiAPI,
		HTTPClient: retryClient.StandardClient(),
	})
	if err != nil {
		return nil, fmt.Errorf("googleGemini: create client: %w", err)
	}
	return &googleGemini{
		client:  ai,
		modelID: modelID,
	}, nil
}

func NewVertexAIGemini(ctx context.Context, projectID string, region string, modelID string) (GoogleGemini, error) {
	retryClient := WrapClient(&http.Client{
		Timeout: 30 * time.Second,
	})
	retryClient.RetryMax = 5
	retryClient.RetryWaitMax = 10 * time.Second

	ai, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:    projectID,
		Location:   region,
		Backend:    genai.BackendVertexAI,
		HTTPClient: retryClient.StandardClient(),
	})
	if err != nil {
		return nil, fmt.Errorf("googleGemini: create client: %w", err)
	}
	return &googleGemini{
		client:  ai,
		modelID: modelID,
	}, nil
}

func (g *googleGemini) Close() error {
	return nil
}

func (g *googleGemini) Chat(ctx context.Context, conversations []*models.Conversation) (string, error) {
	var history []*genai.Content

	for _, conv := range conversations {
		var role genai.Role
		if conv.Role == models.RoleUser {
			role = genai.RoleUser
		} else {
			role = genai.RoleModel
		}

		textPart := genai.NewPartFromText(conv.Content)
		content := genai.NewContentFromParts([]*genai.Part{textPart}, role)
		history = append(history, content)
	}

	config := &genai.GenerateContentConfig{
		SafetySettings: []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: genai.HarmBlockThresholdBlockOnlyHigh,
			},
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: genai.HarmBlockThresholdBlockOnlyHigh,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: genai.HarmBlockThresholdBlockOnlyHigh,
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: genai.HarmBlockThresholdBlockOnlyHigh,
			},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: genai.Ptr(int32(0)),
		},
	}

	resp, err := g.client.Models.GenerateContent(ctx, g.modelID, history, config)
	if err != nil {
		return "", fmt.Errorf("googleGemini: generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("googleGemini: no content generated")
	}

	var responseContent string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			responseContent += part.Text
		}
	}
	if responseContent == "" {
		return "", fmt.Errorf("googleGemini: no text content in response")
	}
	return responseContent, nil
}
