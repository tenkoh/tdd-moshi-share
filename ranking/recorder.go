package ranking

import (
	"encoding/csv"
	"io"
	"math"
)

type SumCounter struct {
	sum   int
	count int
}

type PlayerMean struct {
	player string
	mean   int
}

type Recorder struct {
	sumBoard  map[string]SumCounter
	meanBoard []*PlayerMean
}

func NewRecorder() *Recorder {
	cap_player := 1000
	recorder := new(Recorder)
	recorder.sumBoard = make(map[string]SumCounter)
	recorder.meanBoard = make([]*PlayerMean, 0, cap_player)
	return recorder
}

func (r *Recorder) Aggregate(records *csv.Reader) {
	if records == nil {
		return
	}

	header, err := records.Read()
	if err != nil {
		return
	}
	if !ValidateHeader(header) {
		return
	}

	for {
		record, err := records.Read()
		if err == io.EOF {
			break
		}
		// ある行が無効だったら飛ばして処理を継続する
		if err != nil {
			continue
		}
		r.updateSum(record)
	}
	r.updateMean()
}

func (r *Recorder) updateSum(record []string) {
	parsed, err := ParseRecord(record)
	if err != nil {
		return
	}
	current, exist := r.sumBoard[parsed.player]
	if !exist {
		r.sumBoard[parsed.player] = SumCounter{parsed.score, 1}
	} else {
		r.sumBoard[parsed.player] = SumCounter{
			current.sum + parsed.score,
			current.count + 1,
		}
	}
}

func (r *Recorder) updateMean() {
	for player, sumCounter := range r.sumBoard {
		cnt := sumCounter.count
		if cnt == 0 {
			continue
		}
		fmean := float64(sumCounter.sum) / float64(cnt)
		fmean = math.Round(fmean)
		mean := int(fmean)
		r.meanBoard = append(r.meanBoard, &PlayerMean{player, mean})
	}
	sortMeanBoard(r.meanBoard)
}
