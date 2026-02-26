package widgets

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/style"
	"github.com/ttpreport/ligolo-mp/v2/internal/protocol"
	"github.com/ttpreport/ligolo-mp/v2/internal/session"
)

type InterfacesWidgetElem struct {
	Name    string
	Address string
}

type InterfacesWidget struct {
	*tview.Table
	data            []protocol.NetInterface
	displayData     []InterfacesWidgetElem
	selectedSession *session.Session
	selectedFunc    func(*session.Session, InterfacesWidgetElem)
}

func NewInterfacesWidget() *InterfacesWidget {
	widget := &InterfacesWidget{
		Table: tview.NewTable(),
	}

	widget.SetSelectable(false, false)
	widget.SetBackgroundColor(style.BgColor)
	widget.SetTitle(fmt.Sprintf("[::b]%s", strings.ToUpper("interfaces")))
	widget.SetBorderColor(style.BorderColor)
	widget.SetTitleColor(style.FgColor)
	widget.SetBorder(true)

	widget.SetFocusFunc(func() {
		widget.SetSelectable(true, false)
		widget.ResetSelector()
	})
	widget.SetBlurFunc(func() {
		widget.SetSelectable(false, false)
	})

	widget.Table.SetSelectedFunc(func(row, _ int) {
		idx := row - 1
		if idx >= 0 && idx < len(widget.displayData) && widget.selectedFunc != nil && widget.selectedSession != nil {
			widget.selectedFunc(widget.selectedSession, widget.displayData[idx])
		}
	})

	return widget
}

func (widget *InterfacesWidget) SetSelectedFunc(f func(*session.Session, InterfacesWidgetElem)) {
	widget.selectedFunc = f
}

func (widget *InterfacesWidget) SetData(data []*session.Session) {
	widget.Clear()

	widget.data = nil

	if widget.selectedSession != nil {
		found := false
		for _, sess := range data {
			if sess.ID == widget.selectedSession.ID {
				widget.selectedSession = sess
				found = true
				break
			}
		}
		if !found {
			widget.selectedSession = nil
		}
	}

	for _, session := range data {
		for _, iface := range session.Interfaces.All() {
			widget.data = append(widget.data, iface)
		}
	}

	widget.Refresh()
}

func (widget *InterfacesWidget) SetSelectedSession(sess *session.Session) {
	widget.Clear()
	widget.selectedSession = sess
	widget.Refresh()
}

func (widget *InterfacesWidget) ResetSelector() {
	if len(widget.data) > 0 {
		widget.Select(1, 0) // forcing selection for highlighting to work immediately
	}
}

func (widget *InterfacesWidget) Refresh() {
	widget.displayData = nil

	headers := []string{"Name", "IP"}
	for i := 0; i < len(headers); i++ {
		header := fmt.Sprintf("[::b]%s", strings.ToUpper(headers[i]))
		widget.SetCell(0, i, tview.NewTableCell(header).SetExpansion(1).SetSelectable(false)).SetFixed(1, 0)
	}

	if widget.selectedSession != nil {
		rowId := 1
		for _, elem := range widget.selectedSession.Interfaces.All() {
			for _, IP := range elem.Addresses {
				widget.SetCell(rowId, 0, tview.NewTableCell(elem.Name))
				widget.SetCell(rowId, 1, tview.NewTableCell(IP))
				widget.displayData = append(widget.displayData, InterfacesWidgetElem{
					Name:    elem.Name,
					Address: IP,
				})

				rowId++
			}
		}
	}
}
