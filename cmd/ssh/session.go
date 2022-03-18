package ssh

import (
	"fmt"
	"time"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/utils"
)

func StartDebugJobSession(debug *models.DebugJobV1Alpha, message string) error {
	c := client.NewJobsV1AlphaApi()
	job, err := c.CreateDebugJob(debug)
	utils.Check(err)

	return StartDebugSession(job, message)
}

func StartDebugProjectSession(job *models.JobV1Alpha, message string) error {
	c := client.NewJobsV1AlphaApi()
	job, err := c.CreateJob(job)
	utils.Check(err)

	return StartDebugSession(job, message)
}

func StartDebugSession(job *models.JobV1Alpha, message string) error {
	c := client.NewJobsV1AlphaApi()

	defer func() {
		fmt.Printf("\n")
		fmt.Printf("* Stopping debug session ..\n")

		err := c.StopJob(job.Metadata.Id)

		if err != nil {
			utils.Check(err)
		} else {
			fmt.Printf("* Session stopped\n")
		}
	}()

	fmt.Printf("* Waiting for debug session to boot up .")
	err := waitUntilJobIsRunning(job, func() { fmt.Printf(".") })
	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	fmt.Println()

	// Reload Job
	job, err = c.GetJob(job.Metadata.Id)
	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	// Get SSH key for job
	sshKey, err := c.GetJobDebugSSHKey(job.Metadata.Id)
	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	conn, err := NewConnectionForJob(job, sshKey.Key)
	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}
	defer conn.Close()

	fmt.Printf("* Waiting for ssh daemon to become ready .")
	err = conn.WaitUntilReady(20, func() { fmt.Printf(".") })
	fmt.Println()

	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	fmt.Println(message)

	err = conn.Session()

	if err != nil {
		fmt.Printf("\n[ERROR] %s\n", err)
		return err
	}

	return nil
}

func waitUntilJobIsRunning(job *models.JobV1Alpha, callback func()) error {
	var err error

	c := client.NewJobsV1AlphaApi()

	for {
		time.Sleep(1000 * time.Millisecond)

		job, err = c.GetJob(job.Metadata.Id)

		if err != nil {
			continue
		}

		if job.Status.State == "FINISHED" {
			return fmt.Errorf("Job '%s' has already finished.\n", job.Metadata.Id)
		}

		if job.Status.State == "RUNNING" {
			return nil
		}

		// do some processing between ticks
		callback()
	}
}
