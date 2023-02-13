package asynm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type MissionState string

const (
	// mission state
	MissionInit     MissionState = "MissionInit"
	MissionOpened   MissionState = "MissionOpened"
	MissionFinished MissionState = "MissionFinished"
	MissionStopped  MissionState = "MissionStopped"
	MissionError    MissionState = "MissionError"
	MissionNotExist MissionState = "MissionNotExist"
)

const (
	fieldCountAll     string = "count_all"
	fieldCountCur     string = "count_cur"
	fieldCreateTime   string = "creaete_time"
	fieldFinishTime   string = "finish_time"
	fieldExpireTime   string = "expire_time"
	fieldMissionState string = "mission_state"
	fieldMissionData  string = "mission_data"
)

type asynm struct {
	client *redisClient

	// options
	opt *options
}

type Asynm interface {

	// return mission id
	OpenAsyncMission(data string, count int, expiration int64) (string, error)

	// submit submission result
	SubmitMissionResult(missionId string, idx int, start int64, data string, err error) error

	// return bool success
	CloseAsyncMission(missionId string) error

	// return mission result
	GetMissionResult(missionId string) (MissionResult, error)
}

func NewAsynm(conf *redis.UniversalOptions, opts ...Option) (Asynm, error) {

	// redis client
	client, err := newRedisClient(conf)
	if err != nil {
		return nil, fmt.Errorf("newRedisClient error: %w", err)
	}

	// apply option
	opt := newOption()
	for _, o := range opts {
		o.apply(opt)
	}

	return &asynm{
		client: client,
		opt:    opt,
	}, nil
}

// return mission id
func (a *asynm) OpenAsyncMission(data string, count int, expiration int64) (string, error) {

	uuid, err := GenerateUuid()
	if err != nil {
		return "", fmt.Errorf("GenerateUuid error: %w", err)
	}

	key := fmt.Sprintf("asynm-%s", uuid)

	// init mission in hash
	if err := a.client.HSet(key, fieldMissionState, MissionInit); err != nil {
		return "", fmt.Errorf("Init mission error: %w", err)
	}

	// set key expire
	if expiration == DefaultExpiration {
		expiration = a.opt.expiration
	}

	if expiration != NoExpiration {
		duration := time.Microsecond * time.Duration(expiration)
		if err := a.client.Expire(key, duration); err != nil {
			return "", fmt.Errorf("Init expiration error: %w", err)
		}
	}

	// open mission
	ts := time.Now().UnixMicro()
	if err := a.client.HSet(key, fieldMissionState, MissionOpened, fieldCountAll, count, fieldCountCur, 0, fieldCreateTime, ts, fieldFinishTime, 0, fieldExpireTime, expiration, fieldMissionData, data); err != nil {
		return "", fmt.Errorf("Open mission error: %w", err)
	}

	return uuid, nil
}

// submit submission result
func (a *asynm) SubmitMissionResult(missionId string, idx int, start int64, data string, missionError error) error {

	key := fmt.Sprintf("asynm-%s", missionId)
	if state := a.client.HGet(key, fieldMissionState); state != string(MissionOpened) {
		if state == "" {
			return fmt.Errorf("HGet %s %s error: empty value [%s]", key, fieldMissionState, state)
		}
		return fmt.Errorf("SubmitMissionResult error: invalid state %s", state)
	}

	// check all-cur
	all, err := strconv.Atoi(a.client.HGet(key, fieldCountAll))
	if err != nil {
		return fmt.Errorf("HGet %s %s error: not a number", key, fieldCountAll)
	}

	cur, err := strconv.Atoi(a.client.HGet(key, fieldCountCur))
	if err != nil {
		return fmt.Errorf("HGet %s %s error: not a number", key, fieldCountCur)
	}

	if cur >= all {
		return fmt.Errorf("Mission [%s] state incorrect: count_cur >= count_all", missionId)
	}

	// submit mission result item
	item := newResultItem(start, data, missionError)
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("json.Marshal(item) error: %w", err)
	}

	field := fmt.Sprintf("item_%d", idx)
	if err := a.client.HSet(key, field, string(itemBytes)); err != nil {
		return fmt.Errorf("HSet %s %s error: submit item", key, field)
	}

	// update finish time
	// though it may be covered by other submission, but it is very near the real finish time
	if err := a.client.HSet(key, fieldFinishTime, item.End); err != nil {
		return fmt.Errorf("HSet %s %s error: %w", key, fieldFinishTime, err)
	}

	// update cur
	// return error is equal to rollback
	if err := a.client.HIncrBy(key, fieldCountCur, 1); err != nil {
		return fmt.Errorf("HIncrBy %s %s 1 error: %w", key, fieldCountCur, err)
	}

	return nil
}

// return bool success
func (a *asynm) CloseAsyncMission(missionId string) error {

	key := fmt.Sprintf("asynm-%s", missionId)
	state := a.client.HGet(key, fieldMissionState)
	if state == "" {
		return fmt.Errorf("HGet %s %s error: empty value [%s]", key, fieldMissionState, state)
	}

	return nil
}

// return mission result
func (a *asynm) GetMissionResult(missionId string) (MissionResult, error) {
	return nil, nil
}
