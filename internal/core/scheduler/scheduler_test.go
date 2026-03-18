package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
	if s.jobs == nil {
		t.Error("jobs map not initialized")
	}
}

func TestScheduler_StartStop(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	// Scheduler should be running
	time.Sleep(100 * time.Millisecond)
}

func TestScheduler_AddJob(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	job := &Job{
		ID:       "test-job",
		Name:     "Test Job",
		Schedule: "*/5 * * * * *",
		Handler: func(ctx context.Context) error {
			return nil
		},
		Enabled: true,
	}

	err := s.AddJob(job)
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	retrieved, ok := s.GetJob("test-job")
	if !ok {
		t.Fatal("job not found after adding")
	}
	if retrieved.ID != job.ID {
		t.Errorf("expected ID %s, got %s", job.ID, retrieved.ID)
	}
}

func TestScheduler_AddJob_Duplicate(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	job := &Job{
		ID:       "test-job",
		Schedule: "*/5 * * * * *",
		Handler:  func(ctx context.Context) error { return nil },
		Enabled:  true,
	}

	_ = s.AddJob(job)
	err := s.AddJob(job)
	if err != ErrJobExists {
		t.Errorf("expected ErrJobExists, got %v", err)
	}
}

func TestScheduler_RemoveJob(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	job := &Job{
		ID:       "test-job",
		Schedule: "*/5 * * * * *",
		Handler:  func(ctx context.Context) error { return nil },
		Enabled:  true,
	}

	_ = s.AddJob(job)
	s.RemoveJob("test-job")

	_, ok := s.GetJob("test-job")
	if ok {
		t.Error("job still exists after removal")
	}
}

func TestScheduler_ListJobs(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	for i := 0; i < 3; i++ {
		_ = s.AddJob(&Job{
			ID:       string(rune('a' + i)),
			Schedule: "*/5 * * * * *",
			Handler:  func(ctx context.Context) error { return nil },
			Enabled:  true,
		})
	}

	jobs := s.ListJobs()
	if len(jobs) != 3 {
		t.Errorf("expected 3 jobs, got %d", len(jobs))
	}
}

func TestScheduler_EnableDisableJob(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	job := &Job{
		ID:       "test-job",
		Schedule: "*/5 * * * * *",
		Handler:  func(ctx context.Context) error { return nil },
		Enabled:  true,
	}
	_ = s.AddJob(job)

	s.DisableJob("test-job")
	retrieved, _ := s.GetJob("test-job")
	if retrieved.Enabled {
		t.Error("job should be disabled")
	}

	s.EnableJob("test-job")
	retrieved, _ = s.GetJob("test-job")
	if !retrieved.Enabled {
		t.Error("job should be enabled")
	}
}

func TestScheduler_RunJob(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	var called int32
	job := &Job{
		ID:       "test-job",
		Schedule: "*/5 * * * * *",
		Handler: func(ctx context.Context) error {
			atomic.StoreInt32(&called, 1)
			return nil
		},
		Enabled: true,
	}
	_ = s.AddJob(job)

	err := s.RunJob("test-job")
	if err != nil {
		t.Fatalf("RunJob failed: %v", err)
	}

	if atomic.LoadInt32(&called) != 1 {
		t.Error("job handler not called")
	}
}

func TestScheduler_RunJob_NotFound(t *testing.T) {
	s := New()

	err := s.RunJob("nonexistent")
	if err != ErrJobNotFound {
		t.Errorf("expected ErrJobNotFound, got %v", err)
	}
}

func TestScheduler_InvalidSchedule(t *testing.T) {
	s := New()
	s.Start()
	defer s.Stop()

	job := &Job{
		ID:       "test-job",
		Schedule: "invalid-cron",
		Handler:  func(ctx context.Context) error { return nil },
		Enabled:  true,
	}

	err := s.AddJob(job)
	if err == nil {
		t.Error("expected error for invalid cron expression")
	}
}

func TestJobError(t *testing.T) {
	err := ErrJobExists
	if err.Error() != "job already exists" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	err = ErrJobNotFound
	if err.Error() != "job not found" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}