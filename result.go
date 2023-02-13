package asynm

import (
	"time"
)

type MissionResult interface {

	// state of mission
	GetState() (MissionState, error)

	// count_all, count_current
	GetProgress() (int, int, error)

	// expiration time of the mission
	GetExpiration() (int64, error)

	// array of result in JSON
	GetResults() ([]*ResultItem, error)

	// get specific result by index, max as count_all
	GetResult(int) (*ResultItem, error)

	// missionId
	GetId() string

	// create time, finish time
	GetTime() (int64, int64, error)
}

type missionResult struct {
	missionId string
	r         *redisClient
}

func newMissionResult(client *redisClient, missionId string) MissionResult {
	return &missionResult{
		missionId: missionId,
		r:         client,
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
func (m *missionResult) GetState() (MissionState, error) {
	return MissionInit, nil
}

// count_all, count_current
func (m *missionResult) GetProgress() (int, int, error) {
	return 0, 0, nil
}

// expiration time of the mission
func (m *missionResult) GetExpiration() (int64, error) {
	return 0, nil
}

// array of result in JSON
func (m *missionResult) GetResults() ([]*ResultItem, error) {
	return nil, nil
}

// get specific result by index, max as count_all
func (m *missionResult) GetResult(idx int) (*ResultItem, error) {
	return nil, nil
}

// missionId
func (m *missionResult) GetId() string {
	return m.missionId
}

// create time, finish time
func (m *missionResult) GetTime() (int64, int64, error) {
	return 0, 0, nil
}
