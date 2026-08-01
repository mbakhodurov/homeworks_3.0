package assembly

import "time"

type service struct {
	producer     ShipAssembledProducer
	minBuildTime time.Duration
	maxBuildTime time.Duration
}

func New(producer ShipAssembledProducer, minBuildTimeSec, maxBuildTimeSec int) *service {
	return &service{
		producer:     producer,
		minBuildTime: time.Duration(minBuildTimeSec) * time.Second,
		maxBuildTime: time.Duration(maxBuildTimeSec) * time.Second,
	}
}
