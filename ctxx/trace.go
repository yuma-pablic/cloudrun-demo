package ctxx

import "context"

type contextKey string

const TraceIDKey contextKey = "trace_id"

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	if v := ctx.Value(TraceIDKey); v != nil {
		if traceID, ok := v.(string); ok {
			return traceID
		}
	}
	return ""
}
