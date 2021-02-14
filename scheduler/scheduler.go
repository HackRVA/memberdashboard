package scheduler

import (
	"memberserver/database"
	"memberserver/payments"
	"memberserver/resourcemanager"
	"time"

	log "github.com/sirupsen/logrus"
)

// checkPaymentsInterval - check the resources every 24 hours
const checkPaymentsInterval = 24

// evaluateMemberStatusInterval - check the resources every 25 hours
const evaluateMemberStatusInterval = 25

// resourceStatusCheckInterval - check the resources every hour
const resourceStatusCheckInterval = 1

// Setup Scheduler
//  We want certain tasks to happen on a regular basis
//  The scheduler will make sure that happens
func Setup() {
	checkPayments()
	updateMemberStatus()
	checkResourceACLs()
}

func checkPayments() {
	payments.GetPayments()

	log.Debugf("checking payments")

	// quietly check the resource status on an interval
	ticker := time.NewTicker(checkPaymentsInterval * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				payments.GetPayments()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func updateMemberStatus() {
	checkMemberStatus()
	// quietly check the resource status on an interval
	ticker := time.NewTicker(evaluateMemberStatusInterval * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				checkMemberStatus()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func checkMemberStatus() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	members := db.GetMembers()

	for _, m := range members {
		err = db.EvaluateMemberStatus(m.ID)
		if err != nil {
			log.Errorf("error evaluating member's status: %s", err.Error())
		}
	}

	db.Release()
}

func checkResourceACLs() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
		return
	}

	resources := db.GetResources()

	// on startup we will subscribe to resources and publish an initial status check
	for _, r := range resources {
		resourcemanager.Subscribe(r.Name+"/result", resourcemanager.HealthCheck)
		resourcemanager.CheckStatus(r)
	}

	// quietly check the resource status on an interval
	ticker := time.NewTicker(resourceStatusCheckInterval * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, r := range resources {
					resourcemanager.CheckStatus(r)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	db.Release()
}
