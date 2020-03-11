package common

import (
	"errors"

	"github.com/google/uuid"
)

// Job async job
type Job struct {
	ID     string
	Status string
	Err    error
}

// NewJob job constructure
func NewJob() *Job {
	u1 := uuid.New()
	return &Job{ID: u1.String(), Status: "running", Err: nil}
}

// GetStatus get job status
func (j *Job) GetStatus() (string, error) {
	if j != nil {
		return j.Status, j.Err
	}
	return "", nil
}

// SetFinish set job finish
func (j *Job) SetFinish() {
	if j != nil {
		j.Status = "finish"
		j.Err = nil
	}
}

// SetErr set job err
func (j *Job) SetErr(err error) {
	if j != nil {
		j.Status = "err"
		j.Err = err
	}
}

// CacheJob store job to db
func (j *Job) CacheJob() (string, error) {
	var cache = make(map[string]interface{})
	cache["id"] = j.ID
	cache["status"] = j.Status
	if j.Err != nil {
		cache["err"] = j.Err.Error()
	} else {
		cache["err"] = nil
	}
	return j.ID, CacheJob(cache)
}

// GetCacheJob get job from db
func (j *Job) GetCacheJob() error {
	data, err := GetJob(j.ID)
	if err != nil {
		return err
	}
	j.Status = data["status"]
	j.Err = errors.New(data["err"])
	return nil
}
