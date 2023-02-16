package asynm

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func newMissionResult(client *redisClient, missionId string) (MissionResult, error) {

	key := fmt.Sprintf("asynm-%s", missionId)
	m := client.HGetAll(key)

	result := &missionResult{
		missionId: missionId,
		arr:       []*ResultItem{},
	}

	// count_all
	var err error
	if result.all, err = getMapInt(m, fieldCountAll); err != nil {
		return nil, err
	}

	// count_cur
	if result.cur, err = getMapInt(m, fieldCountCur); err != nil {
		return nil, err
	}

	// state
	if str, ok := m[fieldMissionState]; ok {
		result.state = MissionState(str)
	} else {
		return nil, fmt.Errorf("Field %s not exist", fieldMissionState)
	}

	if result.all == result.cur {
		result.state = MissionFinished
	}

	// time
	if result.createTime, err = getMapInt64(m, fieldCreateTime); err != nil {
		return nil, err
	}

	if result.finishTime, err = getMapInt64(m, fieldFinishTime); err != nil {
		return nil, err
	}

	if result.expireTime, err = getMapInt64(m, fieldExpireTime); err != nil {
		return nil, err
	}

	// data
	for i := 0; i < result.cur; i++ {
		if item, err := getMapResultItem(m, i); err != nil {
			return nil, err
		} else {

			result.arr = append(result.arr, item)
		}
	}

	return result, nil
}

func getMapInt(m map[string]string, k string) (int, error) {
	if str, ok := m[k]; ok {
		if ret, err := strconv.Atoi(str); err == nil {
			return ret, nil
		}
		return 0, fmt.Errorf("Field %s is not a number", k)

	} else {
		return 0, fmt.Errorf("Field %s not exist", k)
	}
}

func getMapInt64(m map[string]string, k string) (int64, error) {
	if str, ok := m[k]; ok {
		if ret, err := strconv.ParseInt(str, 10, 64); err == nil {
			return ret, nil
		}
		return 0, fmt.Errorf("Field %s is not a number", k)

	} else {
		return 0, fmt.Errorf("Field %s not exist", k)
	}
}

func getMapResultItem(m map[string]string, idx int) (*ResultItem, error) {
	k := fmt.Sprintf("item_{%d}", idx)
	if str, ok := m[k]; ok {
		if item, err := parseResultItem(str); err != nil {
			return nil, err
		} else {
			return item, nil
		}

	} else {
		return nil, fmt.Errorf("Field item_%d not exist", idx)
	}
}

// part of result
type ResultItem struct {
	Start  int64
	End    int64
	Data   string
	ErrMsg string
}

func parseResultItem(str string) (*ResultItem, error) {

	item := &ResultItem{}
	if err := json.Unmarshal([]byte(str), item); err != nil {
		return nil, fmt.Errorf("json.Unmarshal ResultItem error: %w", err)
	}

	return item, nil
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
