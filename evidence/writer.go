package evidence

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Writer struct{ encoder *json.Encoder }

func NewWriter(w io.Writer) *Writer {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return &Writer{encoder: encoder}
}

func (w *Writer) Write(run Run) error {
	if err := run.Validate(); err != nil {
		return err
	}
	return w.encoder.Encode(run)
}

func WriteFile(path string, run Run) error {
	if err := run.Validate(); err != nil {
		return err
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := NewWriter(file).Write(run); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
