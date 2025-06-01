package integration_test

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DenisGoldiner/merec"
	"github.com/DenisGoldiner/mss"
)

type intSized struct {
	size int
	val  int
}

func (is intSized) Size() int {
	return is.size
}

func Test_System_Run_FromSlice_LimitedQueue_Priorities_Reject(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	call := func(ctx context.Context, in intSized) (string, error) {
		return strconv.Itoa(in.val), nil
	}
	givenTasks := []intSized{
		{size: 1, val: 0},
		{size: 1, val: 1},
		{size: 1, val: 2},
		{size: 1, val: 3},
		{size: 1, val: 4},
		{size: 1, val: 5},
		{size: 1, val: 6},
		{size: 1, val: 7},
		{size: 1, val: 8},
		{size: 1, val: 9},
		{size: 1, val: 10},
		{size: 1, val: 11},
		{size: 1, val: 12},
		{size: 1, val: 13},
		{size: 1, val: 14},
		{size: 1, val: 15},
		{size: 1, val: 16},
		{size: 1, val: 17},
		{size: 1, val: 18},
		{size: 1, val: 19},
	}
	expResults := []merec.Result[string]{
		merec.ValueResult("0"), merec.ValueResult("1"),
		merec.ValueResult("2"), merec.ValueResult("3"),
		merec.ValueResult("4"), merec.ValueResult("5"),
		merec.ValueResult("6"), merec.ValueResult("7"),
		merec.ValueResult("8"), merec.ValueResult("9"),
		merec.ValueResult("10"), merec.ValueResult("11"),
		merec.ValueResult("12"), merec.ValueResult("13"),
		merec.ValueResult("14"), merec.ValueResult("15"),
		merec.ValueResult("16"), merec.ValueResult("17"),
		merec.ValueResult("18"), merec.ValueResult("19"),
	}

	options := slog.HandlerOptions{Level: slog.LevelDebug}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &options))
	incomeIntense := 3

	source := mss.NewSliceSourcer(givenTasks)
	pre := mss.NewLimitedSizedQueueRejectPreprocessor[intSized, string](logger, incomeIntense)
	post := mss.NewLogPostprocessor[string](logger)
	system := mss.NewSystem[intSized, string](logger, source, pre, post, incomeIntense)
	outputs := system.Run(ctx, call)

	var actResults []merec.Result[string]
	for res := range outputs {
		actResults = append(actResults, res)
	}

	// then
	require.ElementsMatch(t, expResults, actResults)
}
