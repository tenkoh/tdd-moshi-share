package ranking

import (
	"errors"
	"sort"
	"strconv"
)

type ParsedRecord struct {
	player string
	score  int
}

func ValidateHeader(header []string) bool {
	expected := []string{"create_timestamp", "player_id", "score"}
	if len(header) != 3 {
		return false
	}
	for i, h := range expected {
		if header[i] != h {
			return false
		}
	}
	return true
}

// input: []string{date, player, score}
func ParseRecord(r []string) (*ParsedRecord, error) {
	if len(r) != 3 {
		return nil, errors.New("invalid input")
	}
	score, err := strconv.Atoi(r[2])
	if err != nil {
		return nil, err
	}
	// scoreの最大値超え確認は必要？
	return &ParsedRecord{r[1], score}, nil
}

func sortMeanBoard(board []*PlayerMean) {
	// 平均点で並び替え
	sort.Slice(board, func(i, j int) bool {
		return board[i].mean > board[j].mean
	})
}
