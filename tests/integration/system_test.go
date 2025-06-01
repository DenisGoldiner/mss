package integration_test

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DenisGoldiner/merec"
	"github.com/DenisGoldiner/mss"
)

func Test_System_Run_FromSlice_LimitedQueueReject(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	call := func(ctx context.Context, in int) (string, error) {
		return strconv.Itoa(in), nil
	}
	givenTasks := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
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
	poolSize := 3

	source := mss.NewSliceSourcer(givenTasks)
	pre := mss.NewLimitedQueueRejectPreprocessor[int, string]()
	post := mss.NewLogPostprocessor[string](logger)
	system := mss.NewSystem[int, string](logger, source, pre, post, poolSize)
	outputs := system.Run(ctx, call)

	var actResults []merec.Result[string]
	for res := range outputs {
		actResults = append(actResults, res)
	}

	// then
	require.ElementsMatch(t, expResults, actResults)
}

func Test_System_Run_FromChan_LimitedQueueReject(t *testing.T) {
	t.Parallel()

	// given
	ctx := context.Background()
	call := func(ctx context.Context, in int) (string, error) {
		return strconv.Itoa(in), nil
	}
	givenTasks := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
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
	inputs := make(chan int, len(givenTasks))
	poolSize := 3

	source := mss.NewChanSourcer(inputs)
	pre := mss.NewLimitedQueueRejectPreprocessor[int, string]()
	post := mss.NewLogPostprocessor[string](logger)
	system := mss.NewSystem[int, string](logger, source, pre, post, poolSize)
	outputs := system.Run(ctx, call)

	var wgIn sync.WaitGroup
	var wgOut sync.WaitGroup

	// when
	wgIn.Add(1)
	go func() {
		for _, task := range givenTasks {
			inputs <- task
		}
		wgIn.Done()
	}()

	var actResults []merec.Result[string]
	wgOut.Add(1)
	go func() {
		for res := range outputs {
			actResults = append(actResults, res)
		}
		wgOut.Done()
	}()

	wgIn.Wait()
	close(inputs)

	// then
	wgOut.Wait()
	require.ElementsMatch(t, expResults, actResults)
}
