# Cloud Interactions SDK for Go

An unofficial Go GenAI SDK port for the Vertex AI **Interactions & Managed Agents API**.

This package provides a high-performance, lightweight, and highly concurrent Go client tailored explicitly for the real-time streaming SSE (Server-Sent Events) architectures of the Gemini Interactions platform.

---

## Features

* **High-Performance SSE Stream Scanner**: Out-of-the-box support for massive streaming buffers (up to 128MB) optimized specifically to handle raw interleaved multimodal payloads (large video, audio, or image base64 bytes) without memory fragmentation.
* **Managed Agents & Tool Integration**: Complete support for provisioning, executing, and managing enterprise workflows on the [Google Cloud Managed Agents Platform](https://docs.cloud.google.com/gemini-enterprise-agent-platform/build/managed-agents).
* **Dynamic MCP Injection**: Inline override structures allowing developers to attach or re-bind custom Model Context Protocol (MCP) servers at runtime.
* **Sanitized Multi-Tenant Environments**: Provision network access policies and custom source allowlists dynamically per-interaction.

---

## Installation

Include the package in your Go application:

```bash
go get github.com/sourcerepo-genai-sa/cloud-interactions-go
```

---

## Quickstart

Here is how you initialize the client and stream a simple conversational agent turn:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sourcerepo-genai-sa/cloud-interactions-go"
)

func main() {
	ctx := context.Background()
	
	// Initialize pointing to the Vertex Interactions regional gateway
	baseURL := "https://us-central1-aiplatform.googleapis.com/v1beta1/projects/my-project/locations/us-central1/interactions"
	
	client := interactions.NewClient(baseURL)
	client.WithBearerToken("YOUR_ACCESS_TOKEN")

	req := &interactions.InteractionRequest{
		Model:  "gemini-2.5-flash",
		Stream: true,
		Input: []interactions.Content{
			{
				Type: "user_input",
				Content: []interactions.Part{
					{
						Type: "text",
						Text: "Explain Quantum Computing in one short sentence.",
					},
				},
			},
		},
	}

	err := client.StreamCreate(ctx, req, func(event, data string) error {
		fmt.Printf("Event: %s | Chunk Size: %d bytes\n", event, len(data))
		return nil
	})
	if err != nil {
		log.Fatalf("Stream failed: %v", err)
	}
}
```

---

## Standard Types & Future Compatibility

This client has been designed to strictly mimic the casing, parameter structures, and snake_case tags of the official Google USDK releases:
* Uses `omitzero` tag annotations natively.
* Structures mirror the underlying `Interaction`, `Step`, and `Usage` specifications, ensuring that migrating to the official SDK in the future will require zero rewrite of downstream logic.
