package workflows

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func CreateSnapshot(projectName, label string) {
	archiveContent, err := createArchive()
	utils.Check(err)

	projectID := utils.GetProjectId(projectName)
	log.Printf("Project ID: %s\n", projectID)

	if label == "" {
		label = "snapshot"
	}
	log.Printf("Label: %s\n", label)

	c := client.NewWorkflowV1AlphaApi()
	body, err := c.CreateSnapshotWf(projectID, label, archiveContent)
	utils.Check(err)

	fmt.Printf("PPL_ID: %s\n", string(body))
}

// FIXME Respect .gitignore file
func createArchive() ([]byte, error) {
	archiveFileName := "/tmp/snapshot.tgz"
	cmd := exec.Command("rm", "-f", archiveFileName)
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("removing old archive file failed '%s'", err)
	}

	files, err := filepath.Glob("*")
	if err != nil {
		return nil, fmt.Errorf("finding files to archive failed '%s'", err)
	}

	args := append([]string{"czf", archiveFileName}, files...)

	cmd = exec.Command("/bin/tar", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("creating archive file failed '%s'\n%s", out, err)
	}

	archive, err := ioutil.ReadFile(archiveFileName)
	return archive, err
}
