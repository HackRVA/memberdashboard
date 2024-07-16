package scheduler

import (
	"time"

	"github.com/HackRVA/memberserver/pkg/membermgr/services"
)

const (
	// checkPaymentsInterval - check the resources every 24 hours
	checkPaymentsInterval = 24

	// updateMemberCountInterval
	updateMemberCountInterval = 24

	// evaluateMemberStatusInterval - check the resources every 25 hours
	evaluateMemberStatusInterval = 25

	// resourceStatusCheckInterval - check the resources every hour
	resourceStatusCheckInterval = 1

	resourceUpdateInterval = 4

	// checkIPInterval - check the IP Address daily
	checkIPInterval = 24
)

type Scheduler struct{}

type Task struct {
	interval time.Duration
	initFunc func()
	tickFunc func()
}

// Setup Scheduler
//
//	We want certain tasks to happen on a regular basis
//	The scheduler will make sure that happens
func (s *Scheduler) Setup(j services.Job) {
	tasks := []Task{
		{interval: checkPaymentsInterval * time.Hour, initFunc: j.CheckActiveMembersWithoutSubscription, tickFunc: j.CheckActiveMembersWithoutSubscription},
		{interval: checkPaymentsInterval * time.Hour, initFunc: j.CheckMemberSubscriptions, tickFunc: j.CheckMemberSubscriptions},
		{interval: evaluateMemberStatusInterval * time.Hour, initFunc: j.RemovedInvalidUIDs, tickFunc: j.RemovedInvalidUIDs},
		{interval: evaluateMemberStatusInterval * time.Hour, initFunc: j.EnableValidUIDs, tickFunc: j.EnableValidUIDs},
		{interval: resourceStatusCheckInterval * time.Hour, initFunc: j.CheckResourceInit, tickFunc: j.CheckResourceInterval},
		{interval: resourceUpdateInterval * time.Hour, initFunc: j.UpdateResources, tickFunc: j.UpdateResources},
		{interval: checkIPInterval * time.Hour, initFunc: j.CheckIPAddressInterval, tickFunc: j.CheckIPAddressInterval},
		{interval: updateMemberCountInterval * time.Hour, initFunc: j.UpdateMemberCounts, tickFunc: j.UpdateMemberCounts},
	}

	for _, task := range tasks {
		s.scheduleTask(task.interval, task.initFunc, task.tickFunc)
	}
}

func (s *Scheduler) scheduleTask(interval time.Duration, initFunc func(), tickFunc func()) {
	go initFunc()

	// quietly check the resource status on an interval
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go tickFunc()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
