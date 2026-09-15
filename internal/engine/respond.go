package engine

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ju4n97/esquema/internal/config"
)

// executeRespond serializes the final HTTP response payload and performs egress masking.
func (e *Engine) executeRespond(ctx *Context, w http.ResponseWriter, step *config.RespondStep) {
	status := step.Status
	if status == 0 {
		status = http.StatusOK
	}

	body, _ := ctx.EvalAny(step.BodyExpr)

	if step.SchemaRef != "" && body != nil {
		body = e.maskResponse(body, step.SchemaRef)
	}

	for k, v := range step.Headers {
		w.Header().Set(k, v)
	}

	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}

	// Respect user-specified Content-Type; fallback to application/json
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
		w.Header().Set("Content-Type", contentType)
	}

	w.WriteHeader(status)
	if body == nil {
		return
	}

	// Stream directly if Content-Type is JSON
	if strings.HasPrefix(contentType, "application/json") {
		_ = json.NewEncoder(w).Encode(body)
		return
	}

	// Write raw bytes or string directly for non-JSON payloads
	switch v := body.(type) {
	case []byte:
		_, _ = w.Write(v)
	case string:
		_, _ = w.Write([]byte(v))
	default:
		_ = json.NewEncoder(w).Encode(body)
	}
}

// maskResponse filters out fields not declared in the target response schema.
func (e *Engine) maskResponse(data any, schemaRef string) any {
	cleanRef := strings.TrimPrefix(schemaRef, "[]")
	targetSchema, exists := e.cfg.Schemas[cleanRef]
	if !exists {
		return data
	}

	maskMap := func(m map[string]any) map[string]any {
		out := make(map[string]any, len(targetSchema.Fields))
		for fName := range targetSchema.Fields {
			if val, ok := m[fName]; ok {
				out[fName] = val
			}
		}
		return out
	}

	switch v := data.(type) {
	case []map[string]any:
		outList := make([]any, len(v))
		for i, m := range v {
			outList[i] = maskMap(m)
		}
		return outList
	case []any:
		outList := make([]any, len(v))
		for i, item := range v {
			if m, ok := item.(map[string]any); ok {
				outList[i] = maskMap(m)
			} else {
				outList[i] = item
			}
		}
		return outList
	case map[string]any:
		return maskMap(v)
	default:
		return data
	}
}
