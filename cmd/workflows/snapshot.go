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

func CreateSnapshot(projectName, label, archiveName string) {
	archiveContent, err := getArchiveContent(archiveName)
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

	fmt.Printf("%s\n", string(body))
}

func getArchiveContent(archiveName string) ([]byte, error) {
	if archiveName == "" {
		var err error
		archiveName, err = createArchive()

		if err != nil {
			return nil, err
		}
	}

	archive, err := ioutil.ReadFile(archiveName)
	return archive, err
}

// FIXME Respect .gitignore file
func createArchive() (string, error) {
	archiveFileName := "/tmp/snapshot.tgz"
	cmd := exec.Command("rm", "-f", archiveFileName)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("removing old archive file failed '%s'", err)
	}

	files, err := filepath.Glob("*")
	if err != nil {
		return "", fmt.Errorf("finding files to archive failed '%s'", err)
	}

	args := append([]string{"czf", archiveFileName}, files...)

	cmd = exec.Command("/bin/tar", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("creating archive file failed '%s'\n%s", out, err)
	}

	return archiveFileName, nil
}
