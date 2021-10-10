package ranking

import (
	"fmt"
	"io"
	"os"
)

const DISPLAY_MAX_RANK = 10

type Judge struct {
	last   float32
	rank   int
	skip   int
	record *Recorder
	writer io.Writer
}

func NewJudge(r *Recorder) *Judge {
	judge := new(Judge)
	judge.last = 0.0
	judge.rank = 0
	judge.skip = 0
	judge.record = r
	return judge
}

func (j *Judge) SetWriter(w io.Writer) {
	j.writer = w
}

func (j *Judge) Output() {
	w := j.writer
	if j.writer == nil {
		w = os.Stdout
	}
	for _, playerMean := range j.record.meanBoard {

		// 同点を考慮して例えば 1, 1, 3　のような順位を表示するための処理
		if j.last != playerMean.mean {
			j.last = playerMean.mean
			j.rank = j.rank + j.skip + 1
			j.skip = 0
		} else {
			j.skip++
		}

		if j.rank > DISPLAY_MAX_RANK {
			break
		}
		fmt.Fprintf(w, "%d,%s,%.2f\n", j.rank, playerMean.player, playerMean.mean)
	}
}
