package tools

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// ExtractTraceInfo extracts trace_id and span_id from context.
// Returns empty strings if no valid span is found.
func ExtractTraceInfo(ctx context.Context) (traceID, spanID string) {
	span := trace.SpanFromContext(ctx)
	if span == nil || !span.SpanContext().IsValid() {
		return "", ""
	}
	sc := span.SpanContext()
	return sc.TraceID().String(), sc.SpanID().String()
}
