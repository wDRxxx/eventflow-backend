// edited version of https://github.com/go-slog-handler/slog-handler
// to be compatible with https://github.com/cappuccinotm/slogx

package pretty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Handler struct {
	slog.Handler

	format string
	pretty bool
	w      io.Writer
	b      *bytes.Buffer
	m      *sync.Mutex
}

func (h Handler) Handle(ctx context.Context, r slog.Record) error {
	h.m.Lock()

	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()

	var (
		fields = make(map[string]interface{}, r.NumAttrs())
		out    []byte
	)

	if h.format == "json" {
		fields["level"] = strings.ToLower(r.Level.String())
		fields["msg"] = r.Message
		fields["time"] = r.Time.Format(time.DateTime)
	} else {
		out = []byte(fmt.Sprintf("%s %s %s ",
			r.Time.Format(time.DateTime),
			ParseColor(r.Level.String()),
			color.CyanString(r.Message),
		))
	}

	if err := h.Handler.Handle(ctx, r); err != nil {
		return err
	}

	attrs := map[string]any{}
	if err := json.Unmarshal(h.b.Bytes(), &attrs); err != nil {
		return err
	}

	for k, v := range attrs {
		fields[k] = v
	}

	if h.pretty {
		if b, err := json.MarshalIndent(fields, "", "  "); err != nil {
			return err
		} else {
			out = append(out, b...)
		}
	} else {
		if b, err := json.Marshal(fields); err != nil {
			return err
		} else {
			out = append(out, b...)
		}
	}

	h.w.Write(append(out, "\n"...))

	return nil
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) < 1 {
		return &h
	}

	h2 := h
	h2.Handler = h.Handler.WithAttrs(attrs)

	return &h2
}

func (h Handler) WithGroup(name string) slog.Handler {
	h2 := h
	h2.Handler = h.Handler.WithGroup(name)

	return &h2
}

type Options struct {
	*slog.HandlerOptions

	AddSource bool
	Attr      []slog.Attr
	Format    string
	Level     string
	Pretty    bool
}

func NewHandler(out io.Writer, opts *Options) Handler {
	opts.HandlerOptions = &slog.HandlerOptions{
		AddSource: opts.AddSource,
		Level:     ParseLevel(opts.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// skip standard fields
			if map[string]bool{
				slog.LevelKey:   true,
				slog.MessageKey: true,
				slog.TimeKey:    true,
			}[a.Key] {
				return slog.Attr{}
			}

			key := strings.Split(a.Key, ";")

			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)

				dir, file := filepath.Split(s.File)

				a.Value = slog.StringValue(fmt.Sprintf("%s:%d",
					filepath.Join(filepath.Base(dir), file),
					s.Line,
				))
			} else if key[0] == "raw" {
				a.Key = strings.Join(key[1:], ";")
				a.Value = slog.StringValue(fmt.Sprintf("%#v", a.Value.Any()))
			}

			//switch a.Value.Kind() {
			//// other cases
			//
			//case slog.KindAny:
			//	switch v := a.Value.Any().(type) {
			//	case error:
			//		a.Value = fmtErr(v)
			//	}
			//}

			return a
		},
	}

	b := new(bytes.Buffer)

	if !map[string]bool{
		"json": true,
		"text": true,
	}[opts.Format] {
		opts.Format = "json"
	}

	return Handler{
		Handler: slog.NewJSONHandler(b, opts.HandlerOptions),
		format:  opts.Format,
		pretty:  opts.Pretty,
		b:       b,
		m:       &sync.Mutex{},
		w:       out,
	}
}

func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func ParseColor(level string) string {
	switch strings.ToLower(level) {
	case "debug":
		return color.WhiteString(level)
	case "error":
		return color.RedString(level)
	case "warn":
		return color.YellowString(level)
	default:
		return color.GreenString(level)
	}
}
