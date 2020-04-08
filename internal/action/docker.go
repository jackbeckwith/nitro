package action

import (
	"fmt"

	"github.com/craftcms/nitro/validate"
)

// CreateDatabaseContainer is responsible for the creation of a new Docker database and will
// assign a volume and port based on the arguments. Validation of port collisions should occur
// outside of this func and this will only validate engines and versions.
func CreateDatabaseContainer(name, engine, version string, port int) (*Action, error) {
	if err := validate.DatabaseEngineAndVersion(engine, version); err != nil {
		return nil, err
	}

	// get the container path and port based on the engine
	var containerPath string
	var containerPort int
	switch engine {
	case "postgres":
		containerPort = 5432
		containerPath = "/var/lib/postgresql/data"
	default:
		containerPort = 3306
		containerPath = "/var/lib/mysql"
	}

	// create the volumeMount path using the engine, version, and port
	volume := containerVolume(engine, version, port)
	volumeMount := fmt.Sprintf("%s:%s", volume, containerPath)

	// build the container name based on engine, version, and port
	containerName := containerName(engine, version, port)

	// create the port mapping
	portMapping := fmt.Sprintf("%d:%d", port, containerPort)

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "docker", "run", "-v", volumeMount, "--name", containerName, "-d", "--restart=always", "-p", portMapping},
	}, nil
}

func CreateDatabaseVolume(name, engine, version string, port int) (*Action, error) {
	if err := validate.DatabaseEngineAndVersion(engine, version); err != nil {
		return nil, err
	}

	volume := containerVolume(engine, version, port)

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "docker", "volume", "create", volume},
	}, nil
}

func containerName(engine, version string, port int) string {
	return fmt.Sprintf("%s_%s_%d", engine, version, port)
}

func containerVolume(engine, version string, port int) string {
	return fmt.Sprintf("%s_%s_%d", engine, version, port)
}
