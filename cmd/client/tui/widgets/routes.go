package widgets

import (
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/style"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/utils"
	"github.com/ttpreport/ligolo-mp/v2/internal/route"
	"github.com/ttpreport/ligolo-mp/v2/internal/session"
)

type RoutesWidget struct {
	tview.Table
	data            []*RoutesWidgetElem
	selectedSession *session.Session
	selectedFunc    func(string)
}

func NewRoutesWidget() *RoutesWidget {
	widget := &RoutesWidget{
		Table: *tview.NewTable(),
		data:  nil,
	}

	widget.Table.SetSelectable(false, false)
	widget.Table.SetBackgroundColor(style.BgColor)
	widget.Table.SetTitle(fmt.Sprintf("[::b]%s", strings.ToUpper("routes")))
	widget.Table.SetBorderColor(style.BorderColor)
	widget.Table.SetTitleColor(style.FgColor)
	widget.Table.SetBorder(true)

	widget.SetFocusFunc(func() {
		widget.SetSelectable(true, false)
		widget.ResetSelector()
	})
	widget.SetBlurFunc(func() {
		widget.SetSelectable(false, false)
	})

	return widget
}

func (widget *RoutesWidget) SetSelectedSession(sess *session.Session) {
	widget.selectedSession = sess
	widget.Refresh()
}

func (widget *RoutesWidget) FetchElem(row int) *RoutesWidgetElem {
	id := max(0, row-1)
	if len(widget.data) > id {
		return widget.data[id]
	}

	return nil
}

func (widget *RoutesWidget) FetchRow(sess *session.Session) int {
	for row, elem := range widget.data {
		if elem.Session.ID == sess.ID {
			return row + 1
		}
	}

	return 0
}

func (widget *RoutesWidget) SetSelectionChangedFunc(f func(*RoutesWidgetElem)) {
	widget.Table.SetSelectionChangedFunc(func(row, _ int) {
		item := widget.FetchElem(row)
		if item != nil {
			f(item)
		}
	})
}

func (widget *RoutesWidget) SetSelectedFunc(f func(*RoutesWidgetElem)) {
	widget.Table.SetSelectedFunc(func(row, _ int) {
		item := widget.FetchElem(row)
		if item != nil {
			f(item)
		}
	})
}

func (widget *RoutesWidget) SetData(data []*session.Session) {
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
		for _, route := range session.Tun.Routes.All() {
			widget.data = append(widget.data, NewRoutesWidgetElem(route, session))
		}
	}

	widget.Refresh()
	widget.ResetSelector()
}

func (widget *RoutesWidget) ResetSelector() {
	if len(widget.data) > 0 {
		row := 1
		if widget.selectedSession != nil {
			foundRow := widget.FetchRow(widget.selectedSession)
			if foundRow > 0 {
				row = foundRow
			} else {
				widget.selectedSession = nil
			}
		}
		widget.Select(row, 0)
	}
}

func (widget *RoutesWidget) Refresh() {
	headers := []string{"Session", "Route", "Loopback", "Priority", ""}
	for i := 0; i < len(headers); i++ {
		header := fmt.Sprintf("[::b]%s", strings.ToUpper(headers[i]))
		widget.SetCell(0, i, tview.NewTableCell(header).SetExpansion(1).SetSelectable(false)).SetFixed(1, 0)
	}

	sort.Sort(ByRouteOrder(widget.data))

	rowId := 1
	for _, elem := range widget.data {
		if elem.IsSelected(widget.selectedSession) {
			elem.Highlight(true)
		} else {
			elem.Highlight(false)
		}

		widget.SetCell(rowId, 0, elem.Name())
		widget.SetCell(rowId, 1, elem.Cidr())
		widget.SetCell(rowId, 2, elem.IsLoopback())
		widget.SetCell(rowId, 3, elem.Priority())
		widget.SetCell(rowId, 4, elem.Status().SetSelectable(false).SetAlign(tview.AlignCenter))

		rowId++

	}
}

type RoutesWidgetElem struct {
	Route   *route.Route
	Session *session.Session
	bgcolor tcell.Color
}

func NewRoutesWidgetElem(route *route.Route, session *session.Session) *RoutesWidgetElem {
	return &RoutesWidgetElem{
		Route:   route,
		Session: session,
		bgcolor: style.BgColor,
	}
}

func (elem *RoutesWidgetElem) IsSelected(sess *session.Session) bool {
	if sess == nil {
		return false
	}

	if sess.ID != elem.Session.ID {
		return false
	}

	return true
}

func (elem *RoutesWidgetElem) Highlight(h bool) {
	if h {
		elem.bgcolor = style.HighlightColor
	} else {
		elem.bgcolor = style.BgColor
	}
}

func (elem *RoutesWidgetElem) Name() *tview.TableCell {
	val := elem.Session.GetName()
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

func (elem *RoutesWidgetElem) IsLoopback() *tview.TableCell {
	val := utils.HumanBool(elem.Route.IsLoopback)
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

func (elem *RoutesWidgetElem) Cidr() *tview.TableCell {
	val := elem.Route.Cidr
	return tview.NewTableCell(val.String()).SetBackgroundColor(elem.bgcolor)
}

func (elem *RoutesWidgetElem) Status() *tview.TableCell {
	val := "⚑"
	if !elem.Session.IsConnected {
		return tview.NewTableCell(val).SetTextColor(tcell.ColorRed)
	}

	if !elem.Session.IsRelaying {
		return tview.NewTableCell(val).SetTextColor(tcell.ColorBlue)
	}

	return tview.NewTableCell(val).SetTextColor(tcell.ColorGreen)
}

func (elem *RoutesWidgetElem) Priority() *tview.TableCell {
	val := fmt.Sprintf("%d", elem.Route.Metric)
	return tview.NewTableCell(val).SetBackgroundColor(elem.bgcolor)
}

type ByRouteOrder []*RoutesWidgetElem

func (sorter ByRouteOrder) Len() int      { return len(sorter) }
func (sorter ByRouteOrder) Swap(i, j int) { sorter[i], sorter[j] = sorter[j], sorter[i] }
func (sorter ByRouteOrder) Less(i, j int) bool {
	ipCmp := sorter.compareIPs(sorter[i].Route.Cidr.IP, sorter[j].Route.Cidr.IP)
	if ipCmp != 0 {
		return ipCmp < 0
	}

	onesI, _ := sorter[i].Route.Cidr.Mask.Size()
	onesJ, _ := sorter[j].Route.Cidr.Mask.Size()
	if onesI != onesJ {
		return onesI > onesJ
	}

	if sorter[i].Route.Metric != sorter[j].Route.Metric {
		return sorter[i].Route.Metric < sorter[j].Route.Metric
	}

	return sorter[i].Session.ID > sorter[j].Session.ID
}

func (ByRouteOrder) compareIPs(a, b net.IP) int {
	a = a.To16()
	b = b.To16()
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
