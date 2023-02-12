package asynm

import "fmt"

type MissionResult interface {

	// state of mission
	GetState() MissionState

	// count_all, count_current
	GetProgress() (int, int)

	// expiration time of the mission
	GetExpiration() int64

	// array of result in JSON
	GetResults() []string

	// get specific result by index, max as count_all
	GetResult(int) (string, error)

	// missionId
	GetId() string

	// create time, finish time
	GetTime() (int64, int64)
}

type missionResult struct {
	missionId string
	countAll  int
	countCur  int

	createTime int64
	finishTime int64
	expireTime int64

	state MissionState
	data  []string
}

func newMissionResult(r *redisClient, missionId string) MissionResult {
	return nil
}

// state of mission
func (m *missionResult) GetState() MissionState {
	return m.state
}

// count_all, count_current
func (m *missionResult) GetProgress() (int, int) {
	return m.countAll, m.countCur
}

// expiration time of the mission
func (m *missionResult) GetExpiration() int64 {
	return m.expireTime
}

// array of result in JSON
func (m *missionResult) GetResults() []string {
	return m.data
}

// get specific result by index, max as count_all
func (m *missionResult) GetResult(idx int) (string, error) {
	if idx < 0 || idx >= m.countAll {
		return "", fmt.Errorf("Result index error")
	}

	return m.data[idx], nil
}

// missionId
func (m *missionResult) GetId() string {
	return m.missionId
}

// create time, finish time
func (m *missionResult) GetTime() (int64, int64) {
	return m.createTime, m.finishTime
}
