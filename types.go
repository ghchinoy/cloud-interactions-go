// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interactions

import (
	"encoding/json"
	"time"
)

// Role defines the sender of a message part.
type Role string

const (
	RoleUser  Role = "user"
	RoleModel Role = "model"
)

// Part represents a single segment of an interaction input or output.
type Part struct {
	Type                string                      `json:"type,omitzero"`
	Text                string                      `json:"text,omitzero"`
	MimeType            string                      `json:"mime_type,omitzero"`
	URI                 string                      `json:"uri,omitzero"`
	Data                string                      `json:"data,omitzero"`
	Signature           string                      `json:"signature,omitzero"`
	Thought             *ThoughtContent             `json:"thought,omitzero"`
	Call                *ToolCall                   `json:"tool_call,omitzero"`
	Response            *ToolResult                 `json:"tool_response,omitzero"`
	CodeExecutionCall   *CodeExecutionCallContent   `json:"code_execution_call,omitzero"`
	CodeExecutionResult *CodeExecutionResultContent `json:"code_execution_result,omitzero"`
	GoogleSearchCall    *GoogleSearchCallContent    `json:"google_search_call,omitzero"`
	GoogleSearchResult  *GoogleSearchResultContent  `json:"google_search_result,omitzero"`
	URLContextCall      *URLContextCallContent      `json:"url_context_call,omitzero"`
	URLContextResult    *URLContextResultContent    `json:"url_context_result,omitzero"`
	MCPServerToolCall   *MCPServerToolCallContent   `json:"mcp_server_tool_call,omitzero"`
	MCPServerToolResult *MCPServerToolResultContent `json:"mcp_server_tool_result,omitzero"`
	FileSearchCall      *FileSearchCallContent      `json:"file_search_call,omitzero"`
	FileSearchResult    *FileSearchResultContent    `json:"file_search_result,omitzero"`
	GoogleMapsCall      *GoogleMapsCallContent      `json:"google_maps_call,omitzero"`
	GoogleMapsResult    *GoogleMapsResultContent    `json:"google_maps_result,omitzero"`
}

// ThoughtContent represents the reasoning trace from the model.
type ThoughtContent struct {
	Text      string `json:"text,omitempty"`
	Signature string `json:"signature,omitempty"`
	Summary   string `json:"summary,omitempty"`
}

type CodeExecutionCallContent struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Arguments any    `json:"arguments,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type CodeExecutionResultContent struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Output  string `json:"output,omitempty"`
	Outcome string `json:"outcome,omitempty"`
}

type GoogleSearchCallContent struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Arguments any    `json:"arguments,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type GoogleSearchResultContent struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Result any    `json:"result,omitempty"`
}

type URLContextCallContent struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Arguments any    `json:"arguments,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type URLContextResultContent struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Result any    `json:"result,omitempty"`
}

type MCPServerToolCallContent struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Arguments any    `json:"arguments,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type MCPServerToolResultContent struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Result any    `json:"result,omitempty"`
}

type FileSearchCallContent struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Arguments any    `json:"arguments,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type FileSearchResultContent struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Result any    `json:"result,omitempty"`
}

type GoogleMapsCallContent struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Arguments any    `json:"arguments,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type GoogleMapsResultContent struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Result any    `json:"result,omitempty"`
}

// Blob represents inline binary data.
type Blob struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"` // Base64 encoded
}

// File represents a reference to a stored file.
type File struct {
	MimeType string `json:"mime_type"`
	FileURI  string `json:"file_uri"`
}

// ToolCall represents a request from the model to call a function.
type ToolCall struct {
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

// FunctionCall represents the actual function name and arguments.
type FunctionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}

// ToolResult represents the result of a function call.
type ToolResult struct {
	FunctionResponse *FunctionResponse `json:"function_response,omitempty"`
}

// FunctionResponse represents the output data from a function.
type FunctionResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

// Content represents a single turn in an interaction, containing the role and the actual data (parts or text).
type Content struct {
	Type    string `json:"type,omitempty"` // e.g. "user_input" or "model_output"
	Content []Part `json:"content,omitempty"`
	Text    string `json:"text,omitempty"` // Interactions API often flattens text here
}

// InteractionRequest defines the payload for creating a new interaction.
// It supports both standard model input and specialized agent execution.
type InteractionRequest struct {
	Model                 string            `json:"model,omitempty"`
	Agent                 string            `json:"agent,omitempty"`
	AgentConfig           any               `json:"agent_config,omitempty"`
	Input                 any               `json:"input,omitempty"` // Can be string or []Content
	PreviousInteractionID string            `json:"previous_interaction_id,omitempty"`
	Store                 *bool             `json:"store,omitempty"`
	Background            bool              `json:"background,omitempty"`
	Stream                bool              `json:"stream,omitempty"`
	Environment           *Environment      `json:"environment,omitempty"`
	SystemInstruction     any               `json:"system_instruction,omitempty"` // Can be string or Content
	ResponseModalities    []string          `json:"response_modalities,omitempty"`
	ResponseFormat        any               `json:"response_format,omitempty"` // JSON Schema
	ServiceTier           string            `json:"service_tier,omitempty"`    // "flex", "standard", "priority"
	WebhookConfig         *WebhookConfig    `json:"webhook_config,omitempty"`
	GenerationConfig      *GenerationConfig `json:"generation_config,omitempty"`
	Tools                 []Tool            `json:"tools,omitempty"`
}

// WebhookConfig defines optional callback configurations for long-running interactions.
type WebhookConfig struct {
	Url string `json:"url,omitempty"`
}

// Environment specifies an environment to run the agent in.
type Environment struct {
	EnvID   string       `json:"env_id,omitzero"`
	Type    string       `json:"type,omitzero"`
	Sources []Source     `json:"sources,omitzero"`
	Network *NetworkConf `json:"network,omitzero"`
}

type Source struct {
	Type   string `json:"type"`             // "gcs" or "skill_registry"
	Source string `json:"source"`           // "gs://..." or "projects/.../skills/..."
	Target string `json:"target,omitzero"` // "./agent"
}

type NetworkConf struct {
	Allowlist []AllowlistEntry `json:"allowlist,omitzero"`
}

type AllowlistEntry struct {
	Domain string `json:"domain"`
}

// GenerationConfig defines model sampling and output parameters.
type GenerationConfig struct {
	Temperature      *float32     `json:"temperature,omitempty"`
	TopP             *float32     `json:"top_p,omitempty"`
	TopK             *int         `json:"top_k,omitempty"`
	MaxOutputTokens  *int         `json:"max_output_tokens,omitempty"`
	StopSequences    []string     `json:"stop_sequences,omitempty"`
	ResponseMimeType string       `json:"response_mime_type,omitempty"`
	ImageConfig      *ImageConfig `json:"image_config,omitempty"`
}

// ImageConfig defines parameters for image generation.
type ImageConfig struct {
	AspectRatio string `json:"aspect_ratio,omitempty"`
	ImageSize   string `json:"image_size,omitempty"`
}

// Tool represents an external capability the model can use.
type Tool struct {
	Type                 string                `json:"type,omitzero"`
	URL                  string                `json:"url,omitzero"`
	Name                 string                `json:"name,omitzero"`
	Headers              map[string]string     `json:"headers,omitzero"`
	FunctionDeclarations []FunctionDeclaration `json:"function_declarations,omitzero"`
	GoogleSearch         *GoogleSearchTool     `json:"google_search,omitzero"`
	CodeExecution        *CodeExecutionTool    `json:"code_execution,omitzero"`
}

// GoogleSearchTool enables the Google Search tool.
type GoogleSearchTool struct{}

// CodeExecutionTool enables the built-in Code Execution tool.
type CodeExecutionTool struct{}

// FunctionDeclaration defines a tool that the model can call.
type FunctionDeclaration struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters,omitempty"` // JSON Schema
}

// InteractionResponse defines the result of an interaction.
type InteractionResponse struct {
	ID                    string     `json:"id"`
	Name                  string     `json:"name,omitempty"`
	Status                string     `json:"status"` // e.g., "COMPLETED", "WORKING"
	Object                string     `json:"object,omitempty"`
	EnvironmentID         string     `json:"environment_id,omitempty"`
	Outputs               []Content  `json:"outputs,omitempty"`
	Error                 *Error     `json:"error,omitempty"`
	Usage                 *Usage     `json:"usage,omitempty"`
	PreviousInteractionID string     `json:"previous_interaction_id,omitempty"`
	CreateTime            *time.Time `json:"create_time,omitempty"`
	UpdateTime            *time.Time `json:"update_time,omitempty"`
}

// Usage represents token metrics for the interaction.
type Usage struct {
	TotalTokens        int `json:"total_tokens,omitzero"`
	TotalInputTokens   int `json:"total_input_tokens,omitzero"`
	TotalOutputTokens  int `json:"total_output_tokens,omitzero"`
	TotalThoughtTokens int `json:"total_thought_tokens,omitzero"`
}

// Error represents an API error.
type Error struct {
	Code    int    `json:"code,omitzero"`
	Message string `json:"message,omitzero"`
	Status  string `json:"status,omitzero"`
}

// UnmarshalJSON handles the dynamic 'input' field which can be string or array.
func (r *InteractionRequest) MarshalJSON() ([]byte, error) {
	type Alias InteractionRequest
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}
