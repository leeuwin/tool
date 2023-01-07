package ticker

import (
	"time"
)

const (
	defalutInterval = 60 * time.Second
)

type (
	TickerHandler func() error
)

// @Description 定时驱动调用函数，支持根据函数是否返回err执行ok和fail不通的请求间隔，且支持设置fail的最大连续重试次数, 其中，首次执行函数是同步的;
// @Param handler: 函数执行句柄
// @Param okInterval: 函数执行返回err为nil即成功时的执行间隔; 若<=0则退化为defatuleInterval即60s
// @Param failInterval: 函数执行返回err非nil即失败时的执行间隔; 若<=0则退化为okInterval
// @Param retry: 可选变量，若不设置或者设置值<0则将默认无限次数；若=0则不尝试错误; 若0<则为限制handler执行返回错误将最大进行连续的尝试次数, 超过此重试次数后，退化为okInterval
func Ticker(handler TickerHandler, okInterval, failInterval time.Duration, retry ...int) {
	tickerDo(handler, false, okInterval, failInterval, retry...)
}

// @Description 同Ticker(),区别仅为：首次执行函数是异步;
func TickerAsync(handler TickerHandler, okInterval, failInterval time.Duration, retry ...int) {
	tickerDo(handler, true, okInterval, failInterval, retry...)
}

func tickerDo(handler TickerHandler, async bool, okInterval, failInterval time.Duration, retryTimes ...int) {

	//参数校验
	if okInterval <= 0 {
		okInterval = defalutInterval
	}
	if failInterval <= 0 {
		failInterval = okInterval
	}

	retryMax := -1 //默认无限制次数
	if 0 < len(retryTimes) {
		if 0 <= retryTimes[0] {
			retryMax = retryTimes[0]
		}
	}

	// 首次执行(同步模式)
	var firstErr error
	if !async {
		firstErr = handler()
	}

	go func() {
		retryLeft := retryMax
		firstInterval := okInterval
		if !async {
			// 考察同步执行的结果,推算下一次执行间隔
			if firstErr != nil {
				if retryLeft != 0 {
					firstInterval = failInterval
					if 0 < retryLeft {
						retryLeft--
					}
				}
			}
		} else {
			// 首次执行(异步模式)
			firstInterval = 0
		}

		timer := time.NewTimer(firstInterval)
		for range timer.C {
			//execute handler
			err := handler()
			if err != nil {
				if retryLeft != 0 {
					timer.Reset(failInterval)
					if 0 < retryLeft {
						retryLeft--
					}
					continue
				}
			}
			//ok
			if 0 <= retryLeft {
				retryLeft = retryMax
			}
			timer.Reset(okInterval)
		}
	}()
}
