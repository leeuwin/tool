package ticker

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type tickerTest struct {
	name   string
	count  int
	errMap map[int]any
}

func NewTickerTest(name string, em map[int]any) *tickerTest {

	return &tickerTest{
		name:   name,
		errMap: em,
	}
}

func (t *tickerTest) handler() error {

	defer func() {
		t.count++
	}()

	now := time.Now()
	if _, ok := t.errMap[t.count]; ok {
		fmt.Printf("name:%s ts:%d count:%d return error\n", t.name, now.Local().Second(), t.count)
		return errors.New("error on purpers")
	}

	fmt.Printf("name:%s ts:%d count:%d return ok\n", t.name, now.Local().Second(), t.count)
	return nil
}

func TestTicker(t *testing.T) {
	type args struct {
		handler      TickerHandler
		okInterval   time.Duration
		failInterval time.Duration
		retryTimes   []int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				handler:      NewTickerTest("test1", map[int]any{1: nil, 3: nil}).handler,
				okInterval:   5 * time.Second,
				failInterval: 1 * time.Second,
				retryTimes:   []int{1},
			},
		},
		{
			name: "test2",
			args: args{
				handler:      NewTickerTest("test2", map[int]any{1: nil, 5: nil}).handler,
				okInterval:   5 * time.Second,
				failInterval: 1 * time.Second,
				retryTimes:   []int{2},
			},
		},
		{
			name: "test3",
			args: args{
				handler:      NewTickerTest("test3", map[int]any{0: nil, 1: nil, 2: nil, 3: nil, 4: nil, 5: nil, 6: nil, 7: nil}).handler,
				okInterval:   3 * time.Second,
				failInterval: 1 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Ticker(tt.args.handler, tt.args.okInterval, tt.args.failInterval, tt.args.retryTimes...)
		})
	}

	select {}
}
