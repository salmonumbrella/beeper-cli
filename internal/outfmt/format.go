package outfmt

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"text/tabwriter"

	"github.com/itchyny/gojq"
)

type TableWriter struct {
	w       *tabwriter.Writer
	headers []string
	rows    [][]string
}

func NewTableWriter(w io.Writer) *TableWriter {
	return &TableWriter{
		w:    tabwriter.NewWriter(w, 0, 0, 2, ' ', 0),
		rows: make([][]string, 0),
	}
}

func (t *TableWriter) SetHeader(headers []string) {
	t.headers = headers
}

func (t *TableWriter) Append(row []string) {
	t.rows = append(t.rows, row)
}

func (t *TableWriter) Render() {
	// Write header first if set
	if len(t.headers) > 0 {
		for i, h := range t.headers {
			if i > 0 {
				_, _ = t.w.Write([]byte("\t"))
			}
			_, _ = t.w.Write([]byte(h))
		}
		_, _ = t.w.Write([]byte("\n"))
	}
	// Write all rows
	for _, row := range t.rows {
		for i, cell := range row {
			if i > 0 {
				_, _ = t.w.Write([]byte("\t"))
			}
			_, _ = t.w.Write([]byte(cell))
		}
		_, _ = t.w.Write([]byte("\n"))
	}
	_ = t.w.Flush()
}

func WriteJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

func WriteJSONPretty(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func WriteJSONWithQuery(w io.Writer, data any, queryStr string) error {
	if queryStr == "" {
		return WriteJSON(w, data)
	}

	query, err := gojq.Parse(queryStr)
	if err != nil {
		return err
	}

	iter := query.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return err
		}
		if err := WriteJSON(w, v); err != nil {
			return err
		}
	}
	return nil
}

func Output(ctx context.Context, data any, textFn func(io.Writer)) error {
	format := GetFormat(ctx)
	query := GetQuery(ctx)
	return OutputWithQuery(format, query, data, textFn)
}

func OutputWithQuery(format, query string, data any, textFn func(io.Writer)) error {
	switch format {
	case "json":
		return WriteJSONWithQuery(os.Stdout, data, query)
	default:
		textFn(os.Stdout)
		return nil
	}
}
