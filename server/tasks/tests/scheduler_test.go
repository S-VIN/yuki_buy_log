package tasks_test

import (
	"sync/atomic"
	"testing"
	"time"

	"yuki_buy_log/tasks"
)

func TestScheduler_AddTask(t *testing.T) {
	s := tasks.NewScheduler()
	task := tasks.Task{
		Name:     "test_task",
		Interval: 1 * time.Second,
		Run:      func() {},
	}

	s.AddTask(task)

	if s.TaskCount() != 1 {
		t.Errorf("Expected 1 task, got %d", s.TaskCount())
	}

	if s.GetTask(0).Name != "test_task" {
		t.Errorf("Expected task name 'test_task', got '%s'", s.GetTask(0).Name)
	}
}

func TestScheduler_StartStop(t *testing.T) {
	s := tasks.NewScheduler()
	var counter atomic.Int32

	task := tasks.Task{
		Name:     "test_task",
		Interval: 100 * time.Millisecond,
		Run: func() {
			counter.Add(1)
		},
	}

	s.AddTask(task)
	s.Start()

	// Wait for task to execute a few times
	time.Sleep(350 * time.Millisecond)

	s.Stop()

	count := counter.Load()
	if count < 2 || count > 4 {
		t.Errorf("Expected task to run 2-4 times, ran %d times", count)
	}

	// Verify that task stops running after Stop()
	countAfterStop := counter.Load()
	time.Sleep(200 * time.Millisecond)
	countFinal := counter.Load()

	if countFinal != countAfterStop {
		t.Errorf("Task continued running after Stop(). Count before: %d, after: %d", countAfterStop, countFinal)
	}
}

func TestScheduler_MultipleTasks(t *testing.T) {
	s := tasks.NewScheduler()
	var counter1 atomic.Int32
	var counter2 atomic.Int32

	task1 := tasks.Task{
		Name:     "task1",
		Interval: 50 * time.Millisecond,
		Run: func() {
			counter1.Add(1)
		},
	}

	task2 := tasks.Task{
		Name:     "task2",
		Interval: 100 * time.Millisecond,
		Run: func() {
			counter2.Add(1)
		},
	}

	s.AddTask(task1)
	s.AddTask(task2)
	s.Start()

	time.Sleep(250 * time.Millisecond)
	s.Stop()

	count1 := counter1.Load()
	count2 := counter2.Load()

	// task1 should run more times than task2 (twice as fast)
	if count1 <= count2 {
		t.Errorf("Expected task1 to run more times than task2. task1: %d, task2: %d", count1, count2)
	}

	// task1 should run approximately 4-5 times
	if count1 < 3 || count1 > 6 {
		t.Errorf("Expected task1 to run 3-6 times, ran %d times", count1)
	}

	// task2 should run approximately 2 times
	if count2 < 1 || count2 > 3 {
		t.Errorf("Expected task2 to run 1-3 times, ran %d times", count2)
	}
}

func TestScheduler_EmptyScheduler(t *testing.T) {
	s := tasks.NewScheduler()
	s.Start()
	time.Sleep(100 * time.Millisecond)
	s.Stop()
	// Should not panic or error with no tasks
}
