package outfmt

import "context"

type contextKey string

const (
	formatKey contextKey = "output_format"
	queryKey  contextKey = "output_query"
	colorKey  contextKey = "output_color"
)

func WithFormat(ctx context.Context, format string) context.Context {
	return context.WithValue(ctx, formatKey, format)
}

func GetFormat(ctx context.Context) string {
	if v, ok := ctx.Value(formatKey).(string); ok {
		return v
	}
	return "text"
}

func WithQuery(ctx context.Context, query string) context.Context {
	return context.WithValue(ctx, queryKey, query)
}

func GetQuery(ctx context.Context) string {
	if v, ok := ctx.Value(queryKey).(string); ok {
		return v
	}
	return ""
}

func WithColor(ctx context.Context, mode string) context.Context {
	return context.WithValue(ctx, colorKey, mode)
}

func GetColor(ctx context.Context) string {
	if v, ok := ctx.Value(colorKey).(string); ok {
		return v
	}
	return "auto"
}
