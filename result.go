package asynm

import (
	"time"
)

type MissionResult interface {

	// state of mission
	GetState() MissionState

	// count_all, count_current
	GetMax() int
	GetCur() int

	// expiration time of the mission
	GetExpiration() int64

	// array of result in JSON
	GetResults() []*ResultItem

	// get specific result by index, max as count_all
	GetResult(int) *ResultItem

	// missionId
	GetId() string

	// create time, finish time
	GetTime() (int64, int64)
}

type missionResult struct {
	missionId string

	state      MissionState
	all        int
	cur        int
	createTime int64
	finishTime int64
	expireTime int64
	arr        []*ResultItem
}

func newMissionResult(client *redisClient, missionId string) MissionResult {
	return &missionResult{
		missionId: missionId,
	}
}

// part of result
type ResultItem struct {
	Start  int64
	End    int64
	Data   string
	ErrMsg string
}

func newResultItem(start int64, data string, err error) *ResultItem {

	ErrMsg := ""
	if err != nil {
		ErrMsg = err.Error()
	}

	return &ResultItem{
		Start:  start,
		End:    time.Now().UnixMicro(),
		Data:   data,
		ErrMsg: ErrMsg,
	}
}

// state of mission
func (m *missionResult) GetState() MissionState {
	return MissionInit
}

// count_all
func (m *missionResult) GetMax() int {
	return m.all
}

// count_current
func (m *missionResult) GetCur() int {
	return m.cur
}

// expiration time of the mission
func (m *missionResult) GetExpiration() int64 {
	return m.expireTime
}

// array of result in JSON
func (m *missionResult) GetResults() []*ResultItem {
	return m.arr
}

// get specific result by index, max as count_all
func (m *missionResult) GetResult(idx int) *ResultItem {

	if idx > 0 && idx < len(m.arr) {
		return m.arr[idx]
	}

	return nil
}

// missionId
func (m *missionResult) GetId() string {
	return m.missionId
}

// create time, finish time
func (m *missionResult) GetTime() (int64, int64) {
	return m.createTime, m.finishTime
}
