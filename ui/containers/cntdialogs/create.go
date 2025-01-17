package cntdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	containerCreateDialogMaxWidth = 80
	containerCreateDialogHeight   = 17
)

const (
	formFocus = 0 + iota
	categoriesFocus
	categoryPagesFocus
	containerNameFieldFocus
	containerImageFieldFocus
	containerPodFieldFocis
	containerLabelsFieldFocus
	containerRemoveFieldFocus
	containerPortExposeFieldFocus
	containerPortPublishFieldFocus
	containerPortPublishAllFieldFocus
	containerHostnameFieldFocus
	containerIPAddrFieldFocus
	containerMacAddrFieldFocus
	containerNetworkFieldFocus
	containerDNSServersFieldFocus
	containerDNSOptionsFieldFocus
	containerDNSSearchFieldFocus
	containerImageVolumeFieldFocus
	containerVolumeFieldFocus
	containerVolumeDestFocus
)

const (
	basicInfoPageIndex = 0 + iota
	networkingPageIndex
	portPageIndex
	dnsPageIndex
	volumePageIndex
)

// ContainerCreateDialog implements container create dialog
type ContainerCreateDialog struct {
	*tview.Box
	layout                       *tview.Flex
	categoryLabels               []string
	categories                   *tview.TextView
	categoryPages                *tview.Pages
	basicInfoPage                *tview.Flex
	portPage                     *tview.Flex
	networkingPage               *tview.Flex
	dnsPage                      *tview.Flex
	volumePage                   *tview.Flex
	form                         *tview.Form
	display                      bool
	activePageIndex              int
	focusElement                 int
	imageList                    []images.ImageListReporter
	podList                      []*entities.ListPodsReport
	containerNameField           *tview.InputField
	containerImageField          *tview.DropDown
	containerPodField            *tview.DropDown
	containerLabelsField         *tview.InputField
	containerRemoveField         *tview.Checkbox
	containerPortExposeField     *tview.InputField
	containerPortPublishField    *tview.InputField
	ContainerPortPublishAllField *tview.Checkbox
	containerHostnameField       *tview.InputField
	containerIPAddrField         *tview.InputField
	containerMacAddrField        *tview.InputField
	containerNetworkField        *tview.DropDown
	containerDNSServersField     *tview.InputField
	containerDNSOptionsField     *tview.InputField
	containerDNSSearchField      *tview.InputField
	containerVolumeField         *tview.DropDown
	containerVolumeDestField     *tview.InputField
	containerImageVolumeField    *tview.DropDown
	cancelHandler                func()
	createHandler                func()
}

// NewContainerCreateDialog returns new container create dialog primitive ContainerCreateDialog
func NewContainerCreateDialog() *ContainerCreateDialog {
	containerDialog := ContainerCreateDialog{
		Box:                          tview.NewBox(),
		layout:                       tview.NewFlex().SetDirection(tview.FlexRow),
		categories:                   tview.NewTextView(),
		categoryPages:                tview.NewPages(),
		basicInfoPage:                tview.NewFlex(),
		networkingPage:               tview.NewFlex(),
		dnsPage:                      tview.NewFlex(),
		portPage:                     tview.NewFlex(),
		volumePage:                   tview.NewFlex(),
		form:                         tview.NewForm(),
		categoryLabels:               []string{"Basic Information", "Network Settings", "Ports Settings", "DNS Settings", "Volumes Settings"},
		activePageIndex:              0,
		display:                      false,
		containerNameField:           tview.NewInputField(),
		containerImageField:          tview.NewDropDown(),
		containerPodField:            tview.NewDropDown(),
		containerLabelsField:         tview.NewInputField(),
		containerRemoveField:         tview.NewCheckbox(),
		containerPortExposeField:     tview.NewInputField(),
		containerPortPublishField:    tview.NewInputField(),
		ContainerPortPublishAllField: tview.NewCheckbox(),
		containerHostnameField:       tview.NewInputField(),
		containerIPAddrField:         tview.NewInputField(),
		containerMacAddrField:        tview.NewInputField(),
		containerNetworkField:        tview.NewDropDown(),
		containerDNSServersField:     tview.NewInputField(),
		containerDNSOptionsField:     tview.NewInputField(),
		containerDNSSearchField:      tview.NewInputField(),
		containerVolumeField:         tview.NewDropDown(),
		containerVolumeDestField:     tview.NewInputField(),
		containerImageVolumeField:    tview.NewDropDown(),
	}

	bgColor := utils.Styles.ImageHistoryDialog.BgColor

	containerDialog.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	containerDialog.categories.SetBackgroundColor(bgColor)
	containerDialog.categories.SetBorder(true)

	// basic information setup page
	basicInfoPageLabelWidth := 14
	// name field
	containerDialog.containerNameField.SetLabel("name:")
	containerDialog.containerNameField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerNameField.SetBackgroundColor(bgColor)
	containerDialog.containerNameField.SetLabelColor(tcell.ColorWhite)
	// image field
	containerDialog.containerImageField.SetLabel("select image:")
	containerDialog.containerImageField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerImageField.SetBackgroundColor(bgColor)
	containerDialog.containerImageField.SetLabelColor(tcell.ColorWhite)
	// pod field
	containerDialog.containerPodField.SetLabel("select pod:")
	containerDialog.containerPodField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerPodField.SetBackgroundColor(bgColor)
	containerDialog.containerPodField.SetLabelColor(tcell.ColorWhite)
	// labels field
	containerDialog.containerLabelsField.SetLabel("labels:")
	containerDialog.containerLabelsField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerLabelsField.SetBackgroundColor(bgColor)
	containerDialog.containerLabelsField.SetLabelColor(tcell.ColorWhite)
	// remove field
	containerDialog.containerRemoveField.SetLabel("remove container after exit ")
	//containerDialog.containerRemoveField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerRemoveField.SetBackgroundColor(bgColor)
	containerDialog.containerRemoveField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerRemoveField.SetChecked(true)

	// networking setup page
	networkingPageLabelWidth := 13
	// hostname field
	containerDialog.containerHostnameField.SetLabel("hostname:")
	containerDialog.containerHostnameField.SetLabelWidth(networkingPageLabelWidth)
	containerDialog.containerHostnameField.SetBackgroundColor(bgColor)
	containerDialog.containerHostnameField.SetLabelColor(tcell.ColorWhite)
	// IP field
	containerDialog.containerIPAddrField.SetLabel("IP address:")
	containerDialog.containerIPAddrField.SetLabelWidth(networkingPageLabelWidth)
	containerDialog.containerIPAddrField.SetBackgroundColor(bgColor)
	containerDialog.containerIPAddrField.SetLabelColor(tcell.ColorWhite)
	// mac field
	containerDialog.containerMacAddrField.SetLabel("MAC address:")
	containerDialog.containerMacAddrField.SetLabelWidth(networkingPageLabelWidth)
	containerDialog.containerMacAddrField.SetBackgroundColor(bgColor)
	containerDialog.containerMacAddrField.SetLabelColor(tcell.ColorWhite)
	// network field
	containerDialog.containerNetworkField.SetLabel("network:")
	containerDialog.containerNetworkField.SetLabelWidth(networkingPageLabelWidth)
	containerDialog.containerNetworkField.SetBackgroundColor(bgColor)
	containerDialog.containerNetworkField.SetLabelColor(tcell.ColorWhite)

	// ports setup page
	portPageLabelWidth := 15
	// publish field
	containerDialog.containerPortPublishField.SetLabel("publish ports:")
	containerDialog.containerPortPublishField.SetLabelWidth(portPageLabelWidth)
	containerDialog.containerPortPublishField.SetBackgroundColor(bgColor)
	containerDialog.containerPortPublishField.SetLabelColor(tcell.ColorWhite)
	// expose field
	containerDialog.containerPortExposeField.SetLabel("expose ports:")
	containerDialog.containerPortExposeField.SetLabelWidth(portPageLabelWidth)
	containerDialog.containerPortExposeField.SetBackgroundColor(bgColor)
	containerDialog.containerPortExposeField.SetLabelColor(tcell.ColorWhite)
	// publish all field
	containerDialog.ContainerPortPublishAllField.SetLabel("publish all ")
	containerDialog.ContainerPortPublishAllField.SetLabelWidth(portPageLabelWidth)
	containerDialog.ContainerPortPublishAllField.SetBackgroundColor(bgColor)
	containerDialog.ContainerPortPublishAllField.SetLabelColor(tcell.ColorWhite)
	containerDialog.ContainerPortPublishAllField.SetChecked(false)

	// dns setup page
	dnsPageLabelWidth := 13
	// hostname field
	containerDialog.containerDNSServersField.SetLabel("DNS servers:")
	containerDialog.containerDNSServersField.SetLabelWidth(dnsPageLabelWidth)
	containerDialog.containerDNSServersField.SetBackgroundColor(bgColor)
	containerDialog.containerDNSServersField.SetLabelColor(tcell.ColorWhite)
	// IP field
	containerDialog.containerDNSOptionsField.SetLabel("DNS options:")
	containerDialog.containerDNSOptionsField.SetLabelWidth(dnsPageLabelWidth)
	containerDialog.containerDNSOptionsField.SetBackgroundColor(bgColor)
	containerDialog.containerDNSOptionsField.SetLabelColor(tcell.ColorWhite)
	// mac field
	containerDialog.containerDNSSearchField.SetLabel("DNS search:")
	containerDialog.containerDNSSearchField.SetLabelWidth(dnsPageLabelWidth)
	containerDialog.containerDNSSearchField.SetBackgroundColor(bgColor)
	containerDialog.containerDNSSearchField.SetLabelColor(tcell.ColorWhite)

	// volume setup page
	volumePageLabelWidth := 14
	// volume
	containerDialog.containerVolumeField.SetLabel("Volume:")
	containerDialog.containerVolumeField.SetLabelWidth(volumePageLabelWidth)
	containerDialog.containerVolumeField.SetBackgroundColor(bgColor)
	containerDialog.containerVolumeField.SetLabelColor(tcell.ColorWhite)

	// volume
	containerDialog.containerVolumeDestField.SetLabel("Volume Dest:")
	containerDialog.containerVolumeDestField.SetLabelWidth(volumePageLabelWidth)
	containerDialog.containerVolumeDestField.SetBackgroundColor(bgColor)
	containerDialog.containerVolumeDestField.SetLabelColor(tcell.ColorWhite)

	// image volume
	containerDialog.containerImageVolumeField.SetLabel("Image volume:")
	containerDialog.containerImageVolumeField.SetLabelWidth(volumePageLabelWidth)
	containerDialog.containerImageVolumeField.SetBackgroundColor(bgColor)
	containerDialog.containerImageVolumeField.SetLabelColor(tcell.ColorWhite)

	// category pages
	containerDialog.categoryPages.SetBackgroundColor(bgColor)
	containerDialog.categoryPages.SetBorder(true)

	// form
	containerDialog.form.SetBackgroundColor(bgColor)
	containerDialog.form.AddButton("Cancel", nil)
	containerDialog.form.AddButton("Create", nil)
	containerDialog.form.SetButtonsAlign(tview.AlignRight)

	containerDialog.layout.AddItem(tview.NewBox().SetBackgroundColor(bgColor), 1, 0, true)
	containerDialog.setupLayout()
	containerDialog.layout.SetBackgroundColor(bgColor)
	containerDialog.layout.SetBorder(true)
	containerDialog.layout.SetTitle("PODMAN CONTAINER CREATE")
	containerDialog.layout.AddItem(containerDialog.form, dialogs.DialogFormHeight, 0, true)

	containerDialog.setActiveCategory(0)
	return &containerDialog
}

func (d *ContainerCreateDialog) setupLayout() {
	bgColor := utils.Styles.ImageHistoryDialog.BgColor

	emptySpace := func() *tview.Box {
		box := tview.NewBox()
		box.SetBackgroundColor(bgColor)
		return box
	}

	// basic info page
	d.basicInfoPage.SetDirection(tview.FlexRow)
	d.basicInfoPage.AddItem(d.containerNameField, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerImageField, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerPodField, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerLabelsField, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerRemoveField, 1, 0, true)
	d.basicInfoPage.SetBackgroundColor(bgColor)

	// network settigs page
	d.networkingPage.SetDirection(tview.FlexRow)
	d.networkingPage.AddItem(d.containerHostnameField, 1, 0, true)
	d.networkingPage.AddItem(emptySpace(), 1, 0, true)
	d.networkingPage.AddItem(d.containerIPAddrField, 1, 0, true)
	d.networkingPage.AddItem(emptySpace(), 1, 0, true)
	d.networkingPage.AddItem(d.containerMacAddrField, 1, 0, true)
	d.networkingPage.AddItem(emptySpace(), 1, 0, true)
	d.networkingPage.AddItem(d.containerNetworkField, 1, 0, true)
	d.networkingPage.SetBackgroundColor(bgColor)

	// port settigs page
	d.portPage.SetDirection(tview.FlexRow)
	d.portPage.AddItem(d.containerPortPublishField, 1, 0, true)
	d.portPage.AddItem(emptySpace(), 1, 0, true)
	d.portPage.AddItem(d.ContainerPortPublishAllField, 1, 0, true)
	d.portPage.AddItem(emptySpace(), 1, 0, true)
	d.portPage.AddItem(d.containerPortExposeField, 1, 0, true)
	d.portPage.SetBackgroundColor(bgColor)

	// dns settigs page
	d.dnsPage.SetDirection(tview.FlexRow)
	d.dnsPage.AddItem(d.containerDNSServersField, 1, 0, true)
	d.dnsPage.AddItem(emptySpace(), 1, 0, true)
	d.dnsPage.AddItem(d.containerDNSOptionsField, 1, 0, true)
	d.dnsPage.AddItem(emptySpace(), 1, 0, true)
	d.dnsPage.AddItem(d.containerDNSSearchField, 1, 0, true)
	d.dnsPage.SetBackgroundColor(bgColor)

	// volume settigs page
	d.volumePage.SetDirection(tview.FlexRow)
	d.volumePage.AddItem(d.containerVolumeField, 1, 0, true)
	d.volumePage.AddItem(emptySpace(), 1, 0, true)
	d.volumePage.AddItem(d.containerVolumeDestField, 1, 0, true)
	d.volumePage.AddItem(emptySpace(), 1, 0, true)
	d.volumePage.AddItem(d.containerImageVolumeField, 1, 0, true)
	d.volumePage.SetBackgroundColor(bgColor)

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[basicInfoPageIndex], d.basicInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[networkingPageIndex], d.networkingPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[portPageIndex], d.portPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[dnsPageIndex], d.dnsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[volumePageIndex], d.volumePage, true, true)

	// add it to layout.
	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(d.categories, layoutWidth+6, 0, true)
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)

	d.layout.AddItem(layout, 0, 1, true)

}

// Display displays this primitive
func (d *ContainerCreateDialog) Display() {
	d.display = true
	d.initData()
	d.focusElement = categoryPagesFocus
}

// IsDisplay returns true if primitive is shown
func (d *ContainerCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ContainerCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus
func (d *ContainerCreateDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ContainerCreateDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	// form has focus
	case formFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = categoriesFocus // category text view
				d.Focus(delegate)
				d.form.SetFocus(0)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				//d.pullSelectHandler()
				return nil
			}
			return event
		})
		delegate(d.form)
	// category text view
	case categoriesFocus:
		d.categories.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = categoryPagesFocus // category page view
				d.Focus(delegate)
				return nil
			}
			if event.Key() == tcell.KeyDown {
				d.nextCategory()
			}
			if event.Key() == tcell.KeyUp {
				d.previousCategory()
			}
			return event
		})
		delegate(d.categories)
	// basic info page
	case containerNameFieldFocus:
		delegate(d.containerNameField)
	case containerImageFieldFocus:
		delegate(d.containerImageField)
	case containerPodFieldFocis:
		delegate(d.containerPodField)
	case containerLabelsFieldFocus:
		delegate(d.containerLabelsField)
	case containerRemoveFieldFocus:
		delegate(d.containerRemoveField)
	// networking page
	case containerHostnameFieldFocus:
		delegate(d.containerHostnameField)
	case containerIPAddrFieldFocus:
		delegate(d.containerIPAddrField)
	case containerMacAddrFieldFocus:
		delegate(d.containerMacAddrField)
	case containerNetworkFieldFocus:
		delegate(d.containerNetworkField)
	// port page
	// networking page
	case containerPortPublishFieldFocus:
		delegate(d.containerPortPublishField)
	case containerPortPublishAllFieldFocus:
		delegate(d.ContainerPortPublishAllField)
	case containerPortExposeFieldFocus:
		delegate(d.containerPortExposeField)
	// dns page
	case containerDNSServersFieldFocus:
		delegate(d.containerDNSServersField)
	case containerDNSOptionsFieldFocus:
		delegate(d.containerDNSOptionsField)
	case containerDNSSearchFieldFocus:
		delegate(d.containerDNSSearchField)
	// volume page
	case containerVolumeFieldFocus:
		delegate(d.containerVolumeField)
	case containerVolumeDestFocus:
		delegate(d.containerVolumeDestField)
	case containerImageVolumeFieldFocus:
		delegate(d.containerImageVolumeField)
	// category page
	case categoryPagesFocus:
		delegate(d.categoryPages)
	}

}

// InputHandler returns input handler function for this primitive
func (d *ContainerCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container create dialog: event %v received", event.Key())
		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()
			return
		}
		if d.basicInfoPage.HasFocus() {
			if handler := d.basicInfoPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setBasicInfoPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.networkingPage.HasFocus() {
			if handler := d.networkingPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setNetworkSettingsPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.portPage.HasFocus() {
			if handler := d.portPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setPortPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.dnsPage.HasFocus() {
			if handler := d.dnsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setDNSSettingsPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.volumePage.HasFocus() {
			if handler := d.volumePage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setVolumeSettingsPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.categories.HasFocus() {
			if categroryHandler := d.categories.InputHandler(); categroryHandler != nil {
				categroryHandler(event, setFocus)
				return
			}
		}
		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				if event.Key() == tcell.KeyEnter {
					enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
					if enterButton.HasFocus() {
						d.createHandler()
					}
				}
				formHandler(event, setFocus)
				return
			}
		}

	})
}

// SetRect set rects for this primitive.
func (d *ContainerCreateDialog) SetRect(x, y, width, height int) {

	if width > containerCreateDialogMaxWidth {
		emptySpace := (width - containerCreateDialogMaxWidth) / 2
		x = x + emptySpace
		width = containerCreateDialogMaxWidth
	}

	if height > containerCreateDialogHeight {
		emptySpace := (height - containerCreateDialogHeight) / 2
		y = y + emptySpace
		height = containerCreateDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ContainerCreateDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function
func (d *ContainerCreateDialog) SetCancelFunc(handler func()) *ContainerCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetCreateFunc sets form create button selected function
func (d *ContainerCreateDialog) SetCreateFunc(handler func()) *ContainerCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return d
}

func (d *ContainerCreateDialog) setActiveCategory(index int) {
	d.activePageIndex = index
	d.categories.Clear()
	var ctgList []string
	alignedList, _ := utils.AlignStringListWidth(d.categoryLabels)
	for i := 0; i < len(alignedList); i++ {
		if i == index {
			ctgList = append(ctgList, fmt.Sprintf("[white:blue:b]-> %s ", alignedList[i]))
			continue
		}
		ctgList = append(ctgList, fmt.Sprintf("[-:-:-]   %s ", alignedList[i]))
	}
	d.categories.SetText(strings.Join(ctgList, "\n"))

	// switch the page
	d.categoryPages.SwitchToPage(d.categoryLabels[index])
}

func (d *ContainerCreateDialog) nextCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex < len(d.categoryLabels)-1 {
		activePage = activePage + 1
		d.setActiveCategory(activePage)
		return
	}
	d.setActiveCategory(0)
}

func (d *ContainerCreateDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage = activePage - 1
		d.setActiveCategory(activePage)
		return
	}
	d.setActiveCategory(len(d.categoryLabels) - 1)
}

func (d *ContainerCreateDialog) initData() {
	// get available images
	imgList, _ := images.List()
	d.imageList = imgList
	imgOptions := []string{""}
	for i := 0; i < len(d.imageList); i++ {
		if d.imageList[i].ID == "<none>" {
			imgOptions = append(imgOptions, d.imageList[i].ID)
			continue
		}
		imgname := d.imageList[i].Repository + ":" + d.imageList[i].Tag
		imgOptions = append(imgOptions, imgname)
	}

	// get available pods
	podOptions := []string{""}
	podList, _ := pods.List()
	d.podList = podList
	for i := 0; i < len(podList); i++ {
		podOptions = append(podOptions, podList[i].Name)
	}

	// get available networks
	networkOptions := []string{""}
	networkList, _ := networks.List()
	for i := 0; i < len(networkList); i++ {
		networkOptions = append(networkOptions, networkList[i][0])
	}

	// get available volumes
	imageVolumeOptions := []string{"", "ignore", "tmpfs", "anonymous"}
	volumeOptions := []string{""}
	volList, _ := volumes.List()
	for i := 0; i < len(volList); i++ {
		volumeOptions = append(volumeOptions, volList[i].Name)
	}

	d.setActiveCategory(0)
	d.containerNameField.SetText("")
	d.containerImageField.SetOptions(imgOptions, nil)
	d.containerImageField.SetCurrentOption(0)
	d.containerPodField.SetOptions(podOptions, nil)
	d.containerPodField.SetCurrentOption(0)
	d.containerLabelsField.SetText("")
	d.containerRemoveField.SetChecked(false)
	d.containerPortPublishField.SetText("")
	d.ContainerPortPublishAllField.SetChecked(false)
	d.containerPortExposeField.SetText("")
	d.containerHostnameField.SetText("")
	d.containerIPAddrField.SetText("")
	d.containerMacAddrField.SetText("")
	d.containerNetworkField.SetOptions(networkOptions, nil)
	d.containerNetworkField.SetCurrentOption(0)
	d.containerVolumeField.SetOptions(volumeOptions, nil)
	d.containerVolumeField.SetCurrentOption(0)
	d.containerVolumeDestField.SetText("")
	d.containerImageVolumeField.SetOptions(imageVolumeOptions, nil)
	d.containerImageVolumeField.SetCurrentOption(0)
}

func (d *ContainerCreateDialog) setPortPageNextFocus() {
	if d.containerPortPublishField.HasFocus() {
		d.focusElement = containerPortPublishAllFieldFocus
	} else if d.ContainerPortPublishAllField.HasFocus() {
		d.focusElement = containerPortExposeFieldFocus
	} else {
		d.focusElement = formFocus
	}
}

func (d *ContainerCreateDialog) setBasicInfoPageNextFocus() {
	if d.containerNameField.HasFocus() {
		d.focusElement = containerImageFieldFocus
	} else if d.containerImageField.HasFocus() {
		d.focusElement = containerPodFieldFocis
	} else if d.containerPodField.HasFocus() {
		d.focusElement = containerLabelsFieldFocus
	} else if d.containerLabelsField.HasFocus() {
		d.focusElement = containerRemoveFieldFocus
	} else {
		d.focusElement = formFocus
	}
}

func (d *ContainerCreateDialog) setNetworkSettingsPageNextFocus() {
	if d.containerHostnameField.HasFocus() {
		d.focusElement = containerIPAddrFieldFocus
	} else if d.containerIPAddrField.HasFocus() {
		d.focusElement = containerMacAddrFieldFocus
	} else if d.containerMacAddrField.HasFocus() {
		d.focusElement = containerNetworkFieldFocus
	} else {
		d.focusElement = formFocus
	}
}

func (d *ContainerCreateDialog) setDNSSettingsPageNextFocus() {
	if d.containerDNSServersField.HasFocus() {
		d.focusElement = containerDNSOptionsFieldFocus
	} else if d.containerDNSOptionsField.HasFocus() {
		d.focusElement = containerDNSSearchFieldFocus
	} else {
		d.focusElement = formFocus
	}
}

func (d *ContainerCreateDialog) setVolumeSettingsPageNextFocus() {
	if d.containerVolumeField.HasFocus() {
		d.focusElement = containerVolumeDestFocus
	} else if d.containerVolumeDestField.HasFocus() {
		d.focusElement = containerImageVolumeFieldFocus
	} else {
		d.focusElement = formFocus
	}
}

// ContainerCreateOptions returns new network options
func (d *ContainerCreateDialog) ContainerCreateOptions() containers.CreateOptions {
	var (
		labels           = make(map[string]string)
		imageID          string
		podID            string
		dnsServers       []string
		dnsOptions       []string
		dnsSearchDomains []string
		publish          []string
		expose           []string
		volume           string
		imageVolume      string
	)
	log.Info().Msg(d.containerLabelsField.GetText())
	for _, label := range strings.Split(d.containerLabelsField.GetText(), " ") {
		if label != "" {
			split := strings.Split(label, "=")
			if len(split) == 2 {
				key := split[0]
				value := split[1]
				if key != "" && value != "" {
					labels[key] = value
				}
			}
		}
	}
	selectedImageIndex, _ := d.containerImageField.GetCurrentOption()
	if len(d.imageList) > 0 && selectedImageIndex > 0 {
		imageID = d.imageList[selectedImageIndex-1].ID
	}
	selectedPodIndex, _ := d.containerPodField.GetCurrentOption()
	if len(d.podList) > 0 && selectedPodIndex > 0 {
		podID = d.podList[selectedPodIndex-1].Id
	}

	// ports
	for _, p := range strings.Split(d.containerPortPublishField.GetText(), " ") {
		if p != "" {
			publish = append(publish, p)
		}
	}
	for _, e := range strings.Split(d.containerPortExposeField.GetText(), " ") {
		if e != "" {
			expose = append(expose, e)
		}
	}
	// DNS setting
	for _, dns := range strings.Split(d.containerDNSServersField.GetText(), " ") {
		if dns != "" {
			dnsServers = append(dnsServers, dns)
		}
	}
	for _, do := range strings.Split(d.containerDNSOptionsField.GetText(), " ") {
		if do != "" {
			dnsOptions = append(dnsOptions, do)
		}
	}
	for _, ds := range strings.Split(d.containerDNSSearchField.GetText(), " ") {
		if ds != "" {
			dnsSearchDomains = append(dnsSearchDomains, ds)
		}
	}
	_, volume = d.containerVolumeField.GetCurrentOption()
	_, imageVolume = d.containerImageVolumeField.GetCurrentOption()

	_, network := d.containerNetworkField.GetCurrentOption()
	opts := containers.CreateOptions{
		Name:            d.containerNameField.GetText(),
		Image:           imageID,
		Pod:             podID,
		Labels:          labels,
		Remove:          d.containerRemoveField.IsChecked(),
		Hostname:        d.containerHostnameField.GetText(),
		MacAddress:      d.containerMacAddrField.GetText(),
		IPAddress:       d.containerIPAddrField.GetText(),
		Network:         network,
		Publish:         publish,
		Expose:          expose,
		PublishAll:      d.ContainerPortPublishAllField.IsChecked(),
		DNSServer:       dnsServers,
		DNSOptions:      dnsOptions,
		DNSSearchDomain: dnsSearchDomains,
		Volume:          volume,
		VolumeDest:      d.containerVolumeDestField.GetText(),
		ImageVolume:     imageVolume,
	}
	return opts
}
