package containers

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/rs/zerolog/log"
)

func (cnt *Containers) runCommand(cmd string) {
	switch cmd {
	case "create":
		cnt.createDialog.Display()
	case "diff":
		cnt.diff()
	case "inspect":
		cnt.inspect()
	case "kill":
		cnt.kill()
	case "logs":
		cnt.logs()
	case "pause":
		cnt.pause()
	case "prune":
		cnt.cprune()
	case "rename":
		cnt.rename()
	case "port":
		cnt.port()
	case "rm":
		cnt.rm()
	case "start":
		cnt.start()
	case "stop":
		cnt.stop()
	case "top":
		cnt.top()
	case "unpause":
		cnt.unpause()
	}
}

func (cnt *Containers) create() {
	createOpts := cnt.createDialog.ContainerCreateOptions()
	if createOpts.Name == "" || createOpts.Image == "" {
		cnt.errorDialog.SetText("container name or image name is empty")
		cnt.errorDialog.Display()
		return
	}
	cnt.progressDialog.SetTitle("container create in progress")
	cnt.progressDialog.Display()
	create := func() {
		warnings, err := containers.Create(createOpts)
		cnt.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: containers create %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
			return
		}
		if len(warnings) > 0 {
			cnt.messageDialog.SetTitle("CONTAINER CREATE WARNINGS")
			cnt.messageDialog.SetText(strings.Join(warnings, "\n"))
			cnt.messageDialog.Display()
		}
	}
	go create()
}

func (cnt *Containers) diff() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to display diff")
		cnt.errorDialog.Display()
		return
	}
	data, err := containers.Diff(cnt.selectedID)
	if err != nil {
		log.Error().Msgf("view: containers %s", err.Error())
		cnt.errorDialog.SetText(err.Error())
		cnt.errorDialog.Display()
		return
	}
	cnt.messageDialog.SetTitle("podman container diff")
	cnt.messageDialog.SetText(strings.Join(data, "\n"))
	cnt.messageDialog.Display()
}

func (cnt *Containers) inspect() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to inspect")
		cnt.errorDialog.Display()
		return
	}
	data, err := containers.Inspect(cnt.selectedID)
	if err != nil {
		log.Error().Msgf("view: containers %s", err.Error())
		cnt.errorDialog.SetText(err.Error())
		cnt.errorDialog.Display()
		return
	}
	cnt.messageDialog.SetTitle("podman container inspect")
	cnt.messageDialog.SetText(data)
	cnt.messageDialog.Display()
}

func (cnt *Containers) kill() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to kill")
		cnt.errorDialog.Display()
		return
	}
	cnt.progressDialog.SetTitle("container kill in progress")
	cnt.progressDialog.Display()
	kill := func(id string) {
		err := containers.Kill(id)
		cnt.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: containers %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
			return
		}
	}
	go kill(cnt.selectedID)
}

func (cnt *Containers) logs() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to display logs")
		cnt.errorDialog.Display()
		return
	}
	logs, err := containers.Logs(cnt.selectedID)
	if err != nil {
		log.Error().Msgf("view: containers %s", err.Error())
		cnt.errorDialog.SetText(err.Error())
		cnt.errorDialog.Display()
		return
	}
	cntLogs := strings.Join(logs, "\n")
	cntLogs = strings.ReplaceAll(cntLogs, "[", "")
	cntLogs = strings.ReplaceAll(cntLogs, "]", "")
	cnt.messageDialog.SetTitle("podman container logs")
	cnt.messageDialog.SetText(cntLogs)
	cnt.messageDialog.TextScrollToEnd()
	cnt.messageDialog.Display()
}

func (cnt *Containers) pause() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to pause")
		cnt.errorDialog.Display()
		return
	}
	cnt.progressDialog.SetTitle("container pause in progress")
	cnt.progressDialog.Display()
	pause := func(id string) {
		err := containers.Pause(id)
		cnt.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: containers %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
			return
		}
	}
	go pause(cnt.selectedID)
}

func (cnt *Containers) port() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to display port")
		cnt.errorDialog.Display()
		return
	}
	data, err := containers.Port(cnt.selectedID)
	if err != nil {
		log.Error().Msgf("view: containers %s", err.Error())
		cnt.errorDialog.SetText(err.Error())
		cnt.errorDialog.Display()
		return
	}
	cnt.messageDialog.SetTitle("podman container port")
	cnt.messageDialog.SetText(strings.Join(data, "\n"))
	cnt.messageDialog.Display()
}

func (cnt *Containers) cprune() {
	cnt.confirmDialog.SetTitle("podman container prune")
	cnt.confirmData = "prune"
	cnt.confirmDialog.SetText("Are you sure you want to remove all unused containers ?")
	cnt.confirmDialog.Display()
}

func (cnt *Containers) prune() {
	cnt.progressDialog.SetTitle("container purne in progress")
	cnt.progressDialog.Display()
	prune := func() {
		errData, err := containers.Prune()
		cnt.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: containers %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
			return
		}
		if len(errData) > 0 {
			cnt.errorDialog.SetText(strings.Join(errData, "\n"))
			cnt.errorDialog.Display()
		}

	}
	go prune()
}

func (cnt *Containers) rename() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to rename")
		cnt.errorDialog.Display()
		return
	}
	cnt.cmdInputDialog.SetTitle("podman container rename")
	description := fmt.Sprintf("[white::]container name : [black::]%s[white::]\ncontainer ID   : [black::]%s", cnt.selectedName, cnt.selectedID)
	cnt.cmdInputDialog.SetDescription(description)
	cnt.cmdInputDialog.SetSelectButtonLabel("rename")
	cnt.cmdInputDialog.SetLabel("target name")
	cnt.cmdInputDialog.SetSelectedFunc(func() {
		newName := cnt.cmdInputDialog.GetInputText()
		cnt.cmdInputDialog.Hide()
		err := containers.Rename(cnt.selectedID, newName)
		if err != nil {
			log.Error().Msgf("view: containers %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
		}
	})
	cnt.cmdInputDialog.Display()
}

func (cnt *Containers) rm() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to remove")
		cnt.errorDialog.Display()
		return
	}
	cnt.confirmDialog.SetTitle("podman container remove")
	cnt.confirmData = "rm"
	description := fmt.Sprintf("Are you sure you want to remove following container ? \n\nCONTAINER ID : %s", cnt.selectedID)
	cnt.confirmDialog.SetText(description)
	cnt.confirmDialog.Display()
}

func (cnt *Containers) remove() {
	cnt.progressDialog.SetTitle("container remove in progress")
	cnt.progressDialog.Display()
	remove := func(id string) {
		err := containers.Remove(id)
		cnt.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: containers %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
			return
		}
	}
	go remove(cnt.selectedID)
}

func (cnt *Containers) start() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to start")
		cnt.errorDialog.Display()
		return
	}
	cnt.progressDialog.SetTitle("container start in progress")
	cnt.progressDialog.Display()
	start := func(id string) {
		err := containers.Start(id)
		cnt.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: containers %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
			return
		}
	}
	go start(cnt.selectedID)
}

func (cnt *Containers) stop() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to stop")
		cnt.errorDialog.Display()
		return
	}
	cnt.progressDialog.SetTitle("container stop in progress")
	cnt.progressDialog.Display()
	stop := func(id string) {
		err := containers.Stop(id)
		cnt.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: containers %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
			return
		}
	}
	go stop(cnt.selectedID)
}

func (cnt *Containers) top() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to display top")
		cnt.errorDialog.Display()
		return
	}
	data, err := containers.Top(cnt.selectedID)
	if err != nil {
		log.Error().Msgf("view: containers %s", err.Error())
		cnt.errorDialog.SetText(err.Error())
		cnt.errorDialog.Display()
		return
	}
	cnt.topDialog.UpdateResults(data)
	cnt.topDialog.Display()
}

func (cnt *Containers) unpause() {
	if cnt.selectedID == "" {
		cnt.errorDialog.SetText("there is no container to unpause")
		cnt.errorDialog.Display()
		return
	}
	cnt.progressDialog.SetTitle("container unpause in progress")
	cnt.progressDialog.Display()
	unpause := func(id string) {
		err := containers.Unpause(id)
		cnt.progressDialog.Hide()
		if err != nil {
			log.Error().Msgf("view: containers %s", err.Error())
			cnt.errorDialog.SetText(err.Error())
			cnt.errorDialog.Display()
			return
		}
	}
	go unpause(cnt.selectedID)
}
