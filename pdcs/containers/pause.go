package containers

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Pause pauses a pod's containers
func Pause(id string) error {
	log.Debug().Msgf("pdcs: podman container pause %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	return containers.Pause(conn, id, new(containers.PauseOptions))
}

// Unpause pauses a pod's containers
func Unpause(id string) error {
	log.Debug().Msgf("pdcs: podman container unpause %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	return containers.Unpause(conn, id, new(containers.UnpauseOptions))
}
