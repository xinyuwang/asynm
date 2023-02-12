package asynm

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type MissionState int

const (
	// mission state
	MissionOpened   MissionState = 1
	MissionFinished MissionState = 2
	MissionStopped  MissionState = 3
	MissionNotExist MissionState = 4
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
	SubmitMissionResult(missionId string, start int64, data string, err error) error

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
	return "", nil
}

// submit submission result
func (a *asynm) SubmitMissionResult(missionId string, start int64, data string, err error) error {
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
