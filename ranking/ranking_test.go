package ranking

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadCsv(t *testing.T) {
	path := "../testdata/input.csv"
	parsed, err := filepath.Abs(path)
	if err != nil {
		t.Error("could not parse filepath")
	}
	f, err := os.Open(parsed)
	if err != nil {
		t.Error("could not open file")
	}
	defer f.Close()
	reader := csv.NewReader(f)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			t.Error("could not read csv line")
		}
		fmt.Println(record)
	}
}

func TestValidateHeader(t *testing.T) {
	validHeader := []string{"create_timestamp", "player_id", "score"}
	if !ValidateHeader(validHeader) {
		t.Error("validation fails")
	}

	invalidHeader := []string{"create_timestamp ", "player_id ", "score "}
	if ValidateHeader(invalidHeader) {
		t.Error("validation fails")
	}

	invalidHeader = []string{"create_timestamp", "player_id"}
	if ValidateHeader(invalidHeader) {
		t.Error("validation fails")
	}
}

func TestParseRecord(t *testing.T) {
	r := []string{"2021/01/01 12:00", "player0001", "12345"}
	parsed, err := ParseRecord(r)
	if err != nil {
		t.Error("could not parse a valid input")
	}
	if parsed.player != r[1] {
		t.Error("player: not expected value")
	}
	if parsed.score != 12345 {
		t.Error("score: not expected value")
	}
	// scoreは型変換が入るのでチェックする
	if fmt.Sprintf("%T", parsed.score) != "int" {
		t.Error("score type: not expected type")
	}
}

func TestAggregate(t *testing.T) {
	path, _ := filepath.Abs("../testdata/input.csv")
	f, err := os.Open(path)
	if err != nil {
		t.Error("bad test: could not load test input")
	}
	defer f.Close()

	data := csv.NewReader(f)
	recorder := NewRecorder()
	recorder.Aggregate(data)

	sumBoard := recorder.sumBoard
	tests := map[string]int{
		"player0001": 22345,
		"player0002": 10000,
	}
	for player, sum := range tests {
		if sumCounter := sumBoard[player]; sumCounter.sum != sum {
			t.Errorf("expected %d: got %d\n", sum, sumCounter.sum)
		}
	}

	meanBoard := recorder.meanBoard
	meanTests := []PlayerMean{
		{"player0001", 11173},
		{"player0002", 10000},
	}
	hits := 0
	for _, expected := range meanTests {
		for _, calculated := range meanBoard {
			if expected.player == calculated.player {
				hits++
				if expected.mean != calculated.mean {
					t.Errorf("expected %d, got %d\n", expected.mean, calculated.mean)
				}
			}
		}
	}
	if hits != len(meanTests) {
		t.Errorf("player num is insufficient")
	}
}

func TestSortMeanBoard(t *testing.T) {
	test := []*PlayerMean{
		{"player2", 100},
		{"player1", 100},
		{"player3", 200},
		{"player4", 50},
	}
	// 要求にないので同点の時にプレイヤー名の五十音で並び替えるなどは実装していない
	expected := []*PlayerMean{
		{"player3", 200},
		{"player2", 100},
		{"player1", 100},
		{"player4", 50},
	}
	// this method override order
	sortMeanBoard(test)
	for i, player := range expected {
		if test[i].player != player.player {
			t.Errorf("unexpected order at index of %d", i)
		}
	}
}

func TestJudgeOutput(t *testing.T) {
	path, _ := filepath.Abs("../testdata/input.csv")
	f, err := os.Open(path)
	if err != nil {
		t.Error("bad test: could not load test input")
	}
	defer f.Close()

	data := csv.NewReader(f)
	recorder := NewRecorder()
	recorder.Aggregate(data)

	judge := NewJudge(recorder)
	writer := new(bytes.Buffer)

	judge.SetWriter(writer)
	judge.Output()
	expectedSlice := []string{
		"1,player0001,11173",
		"2,player0002,10000",
		"3,player0031,300",
		"3,player0041,300",
		"5,player0021,100",
	}

	s := writer.String()
	s = strings.TrimRight(s, "\n") // trim last newline
	ss := strings.Split(s, "\n")

	for _, expected := range expectedSlice {
		if !isin(expected, ss) {
			t.Errorf("could not find %s\n", expected)
		}
	}
}

func TestPrintCompare(t *testing.T) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%d,%d\n", 1, 2)
	fmt.Fprintf(buf, "%d,%d\n", 3, 4)
	// warning: below case fails
	// if s := buf.String(); s != "1,2\n3, 4\n" {
	// 	t.Errorf("got %s", s)
	// }
	expected := new(bytes.Buffer)
	fmt.Fprintf(expected, "%d,%d\n", 1, 2)
	fmt.Fprintf(expected, "%d,%d\n", 3, 4)
	if buf.String() != expected.String() {
		t.Error("Fprint to beffer test fails")
	}
}

func TestIntegration(t *testing.T) {
	es := []string{
		"1,player0001,11563",
		"2,player0002,10000",
		"3,player001,1000",
		"3,player002,1000",
		"5,player0031,300",
		"5,player0041,300",
		"7,player0021,100",
		"7,player01,100",
		"7,player02,100",
		"10,player031,30",
		"10,player041,30",
	}

	filename := "../testdata/large.csv"
	writer := new(bytes.Buffer)
	err := GetRank(filename, writer)
	if err != nil {
		t.Errorf("Got error: %s\n", err.Error())
	}

	// note: return values' order could not be guaranteed
	s := writer.String()
	s = strings.TrimRight(s, "\n") // trim last newline
	ss := strings.Split(s, "\n")

	for _, expected := range es {
		if !isin(expected, ss) {
			t.Errorf("could not find %s\n", expected)
		}
	}
}

func isin(target string, words []string) bool {
	for _, w := range words {
		if target == w {
			return true
		}
	}
	return false
}
