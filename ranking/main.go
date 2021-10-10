package ranking

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
)

func GetRank(filename string, w io.Writer) error {
	path, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	data := csv.NewReader(f)
	recorder := NewRecorder()
	recorder.Aggregate(data)

	judge := NewJudge(recorder)
	judge.SetWriter(w)
	judge.Output()
	return nil
}
