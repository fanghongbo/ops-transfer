package proc

// trace
var (
	RecvDataTrace = NewDataTrace("RecvDataTrace", 3)
)

// filter
var (
	RecvDataFilter = NewDataFilter("RecvDataFilter", 5)
)

// 统计指标的整体数据
var (
	// 计数统计,正确计数,错误计数, ...
	RecvCnt     = NewSCounterQps("RecvCnt")
	RpcRecvCnt  = NewSCounterQps("RpcRecvCnt")
	HttpRecvCnt = NewSCounterQps("HttpRecvCnt")

	SendToJudgeCnt = NewSCounterQps("SendToJudgeCnt")
	SendToTsDBCnt  = NewSCounterQps("SendToTsDBCnt")
	SendToGraphCnt = NewSCounterQps("SendToGraphCnt")

	SendToJudgeDropCnt = NewSCounterQps("SendToJudgeDropCnt")
	SendToTsDBDropCnt  = NewSCounterQps("SendToTsDBDropCnt")
	SendToGraphDropCnt = NewSCounterQps("SendToGraphDropCnt")

	SendToJudgeFailCnt = NewSCounterQps("SendToJudgeFailCnt")
	SendToTsDBFailCnt  = NewSCounterQps("SendToTsDBFailCnt")
	SendToGraphFailCnt = NewSCounterQps("SendToGraphFailCnt")

	// 发送缓存大小
	JudgeQueuesCnt = NewSCounterBase("JudgeSendCacheCnt")
	TsDBQueuesCnt  = NewSCounterBase("TsDBSendCacheCnt")
	GraphQueuesCnt = NewSCounterBase("GraphSendCacheCnt")
)

func GetAll() []interface{} {
	var ret []interface{}

	ret = make([]interface{}, 0)

	// recv cnt
	ret = append(ret, RecvCnt.Get())
	ret = append(ret, RpcRecvCnt.Get())
	ret = append(ret, HttpRecvCnt.Get())

	// send cnt
	ret = append(ret, SendToJudgeCnt.Get())
	ret = append(ret, SendToTsDBCnt.Get())
	ret = append(ret, SendToGraphCnt.Get())

	// drop cnt
	ret = append(ret, SendToJudgeDropCnt.Get())
	ret = append(ret, SendToTsDBDropCnt.Get())
	ret = append(ret, SendToGraphDropCnt.Get())

	// send fail cnt
	ret = append(ret, SendToJudgeFailCnt.Get())
	ret = append(ret, SendToTsDBFailCnt.Get())
	ret = append(ret, SendToGraphFailCnt.Get())

	// cache cnt
	ret = append(ret, JudgeQueuesCnt.Get())
	ret = append(ret, TsDBQueuesCnt.Get())
	ret = append(ret, GraphQueuesCnt.Get())

	return ret
}
