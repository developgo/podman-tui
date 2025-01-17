package images

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/rs/zerolog/log"
)

// Pull pulls image from registry
func Pull(name string) error {
	log.Debug().Msgf("pdcs: podman image pull %s", name)

	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	_, err = images.Pull(conn, name, new(images.PullOptions).WithQuiet(true))
	if err != nil {
		return err
	}

	return nil
}
