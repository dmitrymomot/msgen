package jobs

import (
	"github.com/gocraft/work"
)

// Queue{{ index .CustomOptions "jobTitle" }} ...
const Queue{{ index .CustomOptions "jobTitle" }} string = "{{ index .CustomOptions "jobName" }}"

// {{ index .CustomOptions "jobTitle" }} job handler
func (c *Worker) {{ index .CustomOptions "jobTitle" }}(job *work.Job) error {
	return nil
}
