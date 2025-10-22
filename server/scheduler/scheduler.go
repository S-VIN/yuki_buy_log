package scheduler

import (
	"log"
	"sync"
	"time"
)

// Task represents a scheduled task with a name, interval, and function to execute
type Task struct {
	Name     string
	Interval time.Duration
	Run      func()
}

// Scheduler manages and executes periodic tasks
type Scheduler struct {
	tasks  []Task
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// New creates a new Scheduler instance
func New() *Scheduler {
	return &Scheduler{
		tasks:  make([]Task, 0),
		stopCh: make(chan struct{}),
	}
}

// AddTask adds a new task to the scheduler
func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

// Start begins executing all scheduled tasks in separate goroutines
func (s *Scheduler) Start() {
	log.Println("Starting scheduler...")
	for _, task := range s.tasks {
		s.wg.Add(1)
		go s.runTask(task)
		log.Printf("Scheduled task '%s' to run every %s", task.Name, task.Interval)
	}
}

// Stop gracefully shuts down the scheduler and waits for all tasks to complete
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	close(s.stopCh)
	s.wg.Wait()
	log.Println("Scheduler stopped")
}

// runTask executes a single task on its defined interval
func (s *Scheduler) runTask(task Task) {
	defer s.wg.Done()

	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Running task '%s'", task.Name)
			task.Run()
		case <-s.stopCh:
			log.Printf("Task '%s' stopped", task.Name)
			return
		}
	}
}
