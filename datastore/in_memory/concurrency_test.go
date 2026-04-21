package in_memory

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/HackRVA/memberserver/models"
)

// TestConcurrentMemberOps hammers the store from many goroutines. Run with
// `go test -race` to detect any unsynchronized map access.
func TestConcurrentMemberOps(t *testing.T) {
	ctx := context.Background()
	store := &In_memory{}

	const writers = 16
	const reads = 100

	var wg sync.WaitGroup

	// writers: AddNewMember
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				_, _ = store.AddNewMember(ctx, models.Member{
					Name:  "m" + strconv.Itoa(id) + "-" + strconv.Itoa(j),
					Email: "m" + strconv.Itoa(id) + "-" + strconv.Itoa(j) + "@ex.com",
					Level: 4,
				})
			}
		}(w)
	}

	// readers: GetMembers + GetMemberCount in parallel
	for r := 0; r < reads; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.GetMembers(ctx)
			_, _ = store.GetMemberCount(ctx, true)
		}()
	}

	// mixed: update + assign
	for u := 0; u < writers; u++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			email := "update-" + strconv.Itoa(id) + "@ex.com"
			_, _ = store.AddNewMember(ctx, models.Member{Name: "u" + strconv.Itoa(id), Email: email})
			_ = store.UpdateMember(ctx, models.Member{Name: "updated", Email: email})
			_, _ = store.AssignRFID(ctx, email, "rfid"+strconv.Itoa(id))
		}(u)
	}

	wg.Wait()
}

func TestConcurrentResourceOps(t *testing.T) {
	ctx := context.Background()
	store := &In_memory{}

	var wg sync.WaitGroup
	for w := 0; w < 32; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, _ = store.RegisterResource(ctx, "res-"+strconv.Itoa(id), "1.1.1."+strconv.Itoa(id%255), false)
			_ = store.GetResources(ctx)
		}(w)
	}
	wg.Wait()
}
