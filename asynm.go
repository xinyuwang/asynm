package asynm

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type MissionState int

const (
	// mission state
	MissionInit     MissionState = 1
	MissionOpened   MissionState = 2
	MissionFinished MissionState = 3
	MissionStopped  MissionState = 4
	MissionError    MissionState = 5
	MissionNotExist MissionState = 6
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
	CloseAsyncMission(missionId string) (bool, error)

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
	if err := a.client.HSet(key, "mission_state", MissionInit); err != nil {
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
	if err := a.client.HSet(key, "mission_state", MissionOpened, "count_all", count, "count_cur", 0, "create_time", ts, "finish_time", 0, "expire_time", expiration, "mission_data", data); err != nil {
		return "", fmt.Errorf("Open mission error: %w", err)
	}

	return "", nil
}

// submit submission result
func (a *asynm) SubmitMissionResult(missionId string, idx int, start int64, data string, err error) error {

	return nil
}

// return bool success
func (a *asynm) CloseAsyncMission(missionId string) (bool, error) {
	return false, nil
}

// return mission result
func (a *asynm) GetMissionResult(missionId string) (MissionResult, error) {
	return nil, nil
}
