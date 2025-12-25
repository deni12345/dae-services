package prettylog

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
	otel "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

type PrettyHandler struct {
	slog.Handler
	internalLogger *log.Logger
	otelLogger     otel.Logger
}

func NewPrettyHandler(w io.Writer, otelLogger otel.Logger, opts *slog.HandlerOptions) *PrettyHandler {
	return &PrettyHandler{
		Handler:        slog.NewTextHandler(w, opts),
		internalLogger: log.New(w, "", 0),
		otelLogger:     otelLogger,
	}
}

func Fatal(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler:        h.Handler.WithAttrs(attrs),
		internalLogger: h.internalLogger,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler:        h.Handler.WithGroup(name),
		internalLogger: h.internalLogger,
	}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()
	var otelSeverity otel.Severity

	switch r.Level {
	case slog.LevelDebug:
		level = color.BlueString(level)
		otelSeverity = otel.SeverityDebug
	case slog.LevelInfo:
		level = color.GreenString(level)
		otelSeverity = otel.SeverityInfo
	case slog.LevelWarn:
		level = color.YellowString(level)
		otelSeverity = otel.SeverityWarn
	case slog.LevelError:
		level = color.RedString(level)
		otelSeverity = otel.SeverityError
	default:
		level = color.GreenString(level)
		otelSeverity = otel.SeverityInfo
	}

	fields := buildFieldsMap(r)

	b, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[2006-01-02 15:04:05]")
	msg := color.CyanString(r.Message)

	h.internalLogger.Println(color.YellowString(timeStr), level, msg, color.WhiteString(string(b)))

	// Emit to OTel logger provider if configured
	if h.otelLogger != nil {
		logRecord := otel.Record{}
		logRecord.SetTimestamp(r.Time)
		logRecord.SetBody(otel.StringValue(r.Message))
		logRecord.SetSeverity(otelSeverity)
		logRecord.SetSeverityText(r.Level.String())
		h.emitToOTel(ctx, r, logRecord, fields)
	}

	return nil
}

func buildFieldsMap(r slog.Record) map[string]interface{} {
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		appendAttr(fields, "", a)
		return true
	})
	return fields
}

func appendAttr(dst map[string]interface{}, prefix string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	key := a.Key

	if prefix != "" {
		if key != "" {
			key = prefix + "." + key
		} else {
			key = prefix // don't add trailing dot
		}
	}
	if a.Value.Kind() == slog.KindGroup {
		for _, ga := range a.Value.Group() {
			appendAttr(dst, key, ga)
		}
	} else {
		dst[key] = a.Value.Any()
	}
}

func (h *PrettyHandler) emitToOTel(ctx context.Context, r slog.Record, logRecord otel.Record, fields map[string]interface{}) {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		logRecord.AddAttributes(otel.String("trace_id", sc.TraceID().String()))
		logRecord.AddAttributes(otel.String("span_id", sc.SpanID().String()))
	}

	for k, v := range fields {
		logRecord.AddAttributes(otel.String(k, formatValue(v)))
	}

	h.otelLogger.Emit(ctx, logRecord)
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		b, _ := json.Marshal(val)
		return string(b)
	}
}
