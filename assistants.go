package telnyx

import (
	"context"
	"fmt"
)

// TelephonySettings controls call handling behaviour for an assistant.
type TelephonySettings struct {
	NoiseSuppression bool `json:"noise_suppression,omitempty"`
	SilenceTimeoutMS int  `json:"silence_timeout_ms,omitempty"`
	MaxCallDuration  int  `json:"max_call_duration_secs,omitempty"`
}

// MessagingSettings controls SMS/MMS behaviour for an assistant.
type MessagingSettings struct {
	MessagingProfileID  string `json:"messaging_profile_id,omitempty"`
	InactivityTimeoutMS int    `json:"inactivity_timeout_ms,omitempty"`
}

// ExternalLLM points to a custom LLM endpoint.
type ExternalLLM struct {
	URL   string `json:"url"`
	Model string `json:"model,omitempty"`
}

// FallbackConfig configures a backup model when the primary fails.
type FallbackConfig struct {
	Model string `json:"model"`
}

// MCPServer references an MCP server integration.
type MCPServer struct {
	ID           string   `json:"id"`
	AllowedTools []string `json:"allowed_tools,omitempty"`
}

// AssistantTranscriptionConfig controls STT for the assistant.
type AssistantTranscriptionConfig struct {
	Model    string `json:"model,omitempty"`
	Language string `json:"language,omitempty"`
}

// CreateAssistantRequest is the payload for creating a new AI assistant.
type CreateAssistantRequest struct {
	Name         string `json:"name"`
	Instructions string `json:"instructions"`

	Model           string          `json:"model,omitempty"`
	Description     string          `json:"description,omitempty"`
	Greeting        string          `json:"greeting,omitempty"`
	ToolIDs         []string        `json:"tool_ids,omitempty"`
	MCPServers      []MCPServer     `json:"mcp_servers,omitempty"`
	LLMAPIKeyRef    string          `json:"llm_api_key_ref,omitempty"`
	ExternalLLM     *ExternalLLM    `json:"external_llm,omitempty"`
	FallbackConfig  *FallbackConfig `json:"fallback_config,omitempty"`
	EnabledFeatures []string        `json:"enabled_features,omitempty"` // telephony | messaging

	VoiceSettings     *VoiceSettings                `json:"voice_settings,omitempty"`
	Transcription     *AssistantTranscriptionConfig `json:"transcription,omitempty"`
	TelephonySettings *TelephonySettings            `json:"telephony_settings,omitempty"`
	MessagingSettings *MessagingSettings            `json:"messaging_settings,omitempty"`

	InterruptionSettings *InterruptionSettings `json:"interruption_settings,omitempty"`

	Tags                              []string       `json:"tags,omitempty"`
	DynamicVariables                  map[string]any `json:"dynamic_variables,omitempty"`
	DynamicVariablesWebhookURL        string         `json:"dynamic_variables_webhook_url,omitempty"`
	DynamicVariablesWebhookTimeoutMS  int            `json:"dynamic_variables_webhook_timeout_ms,omitempty"`
}

// Assistant is the assistant object returned by the API.
type Assistant struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Model        string `json:"model"`
	Instructions string `json:"instructions"`
	Description  string `json:"description"`
	Greeting     string `json:"greeting"`
	CreatedAt    string `json:"created_at"`
	VersionID    string `json:"version_id"`

	ToolIDs         []string        `json:"tool_ids"`
	MCPServers      []MCPServer     `json:"mcp_servers"`
	EnabledFeatures []string        `json:"enabled_features"`
	Tags            []string        `json:"tags"`
	DynamicVariables map[string]any `json:"dynamic_variables"`

	VoiceSettings     *VoiceSettings                `json:"voice_settings"`
	Transcription     *AssistantTranscriptionConfig `json:"transcription"`
	TelephonySettings *TelephonySettings            `json:"telephony_settings"`
	MessagingSettings *MessagingSettings            `json:"messaging_settings"`
	ExternalLLM       *ExternalLLM                  `json:"external_llm"`
	FallbackConfig    *FallbackConfig               `json:"fallback_config"`
}

type assistantResponse struct {
	Data Assistant `json:"data"`
}

type assistantListResponse struct {
	Data []Assistant `json:"data"`
	Meta struct {
		TotalResults int    `json:"total_results"`
		TotalPages   int    `json:"total_pages"`
		PageSize     int    `json:"page_size"`
		PageNumber   int    `json:"page_number"`
	} `json:"meta"`
}

// CreateAssistant creates a new AI assistant.
func (c *Client) CreateAssistant(ctx context.Context, req CreateAssistantRequest) (*Assistant, error) {
	var resp assistantResponse
	if err := c.post(ctx, "/ai/assistants", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetAssistant retrieves an assistant by ID.
func (c *Client) GetAssistant(ctx context.Context, id string) (*Assistant, error) {
	var resp assistantResponse
	if err := c.get(ctx, fmt.Sprintf("/ai/assistants/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpdateAssistantRequest carries partial updates for an assistant.
type UpdateAssistantRequest struct {
	Name         string `json:"name,omitempty"`
	Instructions string `json:"instructions,omitempty"`
	Model        string `json:"model,omitempty"`
	Description  string `json:"description,omitempty"`
	Greeting     string `json:"greeting,omitempty"`

	ToolIDs         []string        `json:"tool_ids,omitempty"`
	MCPServers      []MCPServer     `json:"mcp_servers,omitempty"`
	EnabledFeatures []string        `json:"enabled_features,omitempty"`

	VoiceSettings     *VoiceSettings                `json:"voice_settings,omitempty"`
	Transcription     *AssistantTranscriptionConfig `json:"transcription,omitempty"`
	TelephonySettings *TelephonySettings            `json:"telephony_settings,omitempty"`
	MessagingSettings *MessagingSettings            `json:"messaging_settings,omitempty"`
	FallbackConfig    *FallbackConfig               `json:"fallback_config,omitempty"`
	ExternalLLM       *ExternalLLM                  `json:"external_llm,omitempty"`

	InterruptionSettings *InterruptionSettings `json:"interruption_settings,omitempty"`

	Tags                             []string       `json:"tags,omitempty"`
	DynamicVariables                 map[string]any `json:"dynamic_variables,omitempty"`
	DynamicVariablesWebhookURL       string         `json:"dynamic_variables_webhook_url,omitempty"`
	DynamicVariablesWebhookTimeoutMS int            `json:"dynamic_variables_webhook_timeout_ms,omitempty"`
}

// UpdateAssistant applies partial updates to an existing assistant.
func (c *Client) UpdateAssistant(ctx context.Context, id string, req UpdateAssistantRequest) (*Assistant, error) {
	var resp assistantResponse
	if err := c.patch(ctx, fmt.Sprintf("/ai/assistants/%s", id), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteAssistant permanently removes an assistant.
func (c *Client) DeleteAssistant(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/ai/assistants/%s", id))
}

// ListAssistantsParams filters the assistant list.
type ListAssistantsParams struct {
	PageNumber int
	PageSize   int
}

// ListAssistants returns all assistants for the account.
func (c *Client) ListAssistants(ctx context.Context, params ListAssistantsParams) ([]Assistant, error) {
	path := "/ai/assistants"
	if params.PageNumber > 0 || params.PageSize > 0 {
		path = fmt.Sprintf("%s?page[number]=%d&page[size]=%d", path, params.PageNumber, params.PageSize)
	}
	var resp assistantListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ChatMessage is a single turn in an assistant chat.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest sends messages to an assistant for a response.
type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse holds the assistant's reply.
type ChatResponse struct {
	Data struct {
		Message ChatMessage `json:"message"`
	} `json:"data"`
}

// Chat sends a message to an assistant and returns the reply (beta).
func (c *Client) Chat(ctx context.Context, assistantID string, req ChatRequest) (*ChatResponse, error) {
	var resp ChatResponse
	if err := c.post(ctx, fmt.Sprintf("/ai/assistants/%s/chat", assistantID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
