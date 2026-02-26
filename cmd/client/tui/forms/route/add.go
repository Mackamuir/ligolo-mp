package route

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/forms"
)

type AddRouteForm struct {
	tview.Flex
	form      *tview.Form
	submitBtn *tview.Button
	cancelBtn *tview.Button
}

var (
	add_route_cidr = forms.FormVal[string]{
		Hint: "A CIDR that will be routed via this session.\n\nExample:\n10.10.5.0/24",
	}

	add_route_metric = forms.FormVal[int]{
		Hint: "A priority of the route in case of same route for different sessions. Lower number has more priority",
	}

	add_route_loopback = forms.FormVal[bool]{
		Hint: "If checked, specified CIDR will address the machine running the agent itself, i.e. localhost. Use this instead of port forwarding.",
	}
)

func NewAddRouteFormWithCIDR(cidr string) *AddRouteForm {
	add_route_cidr.Last = cidr
	return NewAddRouteForm()
}

func NewAddRouteForm() *AddRouteForm {
	form := &AddRouteForm{
		Flex:      *tview.NewFlex(),
		form:      tview.NewForm(),
		submitBtn: tview.NewButton("Submit"),
		cancelBtn: tview.NewButton("Cancel"),
	}

	hintBox := tview.NewTextView()
	hintBox.SetTitle("HINT")
	hintBox.SetTitleAlign(tview.AlignCenter)
	hintBox.SetBorder(true)
	hintBox.SetBorderPadding(1, 1, 1, 1)

	form.form.SetTitle("Add route").SetTitleAlign(tview.AlignCenter)
	form.form.SetBorder(true)
	form.form.SetButtonsAlign(tview.AlignCenter)

	cidrField := tview.NewInputField()
	cidrField.SetLabel("CIDR")
	cidrField.SetText(add_route_cidr.Last)
	cidrField.SetFocusFunc(func() {
		hintBox.SetText(add_route_cidr.Hint)
	})
	cidrField.SetChangedFunc(func(text string) {
		add_route_cidr.Last = text
	})
	form.form.AddFormItem(cidrField)

	metricField := tview.NewInputField()
	metricField.SetLabel("Priority")
	metricField.SetAcceptanceFunc(tview.InputFieldInteger)
	metricField.SetText(fmt.Sprint(add_route_metric.Last))
	metricField.SetFocusFunc(func() {
		hintBox.SetText(add_route_metric.Hint)
	})
	metricField.SetChangedFunc(func(text string) {
		val, err := strconv.Atoi(text)
		if err == nil {
			add_route_metric.Last = val
		}
	})
	form.form.AddFormItem(metricField)

	loopbackField := tview.NewCheckbox()
	loopbackField.SetLabel("Loopback")
	loopbackField.SetChecked(add_route_loopback.Last)
	loopbackField.SetFocusFunc(func() {
		hintBox.SetText(add_route_loopback.Hint)
	})
	loopbackField.SetChangedFunc(func(checked bool) {
		add_route_loopback.Last = checked
	})
	loopbackField.SetBlurFunc(func() {
		hintBox.Clear()
	})
	form.form.AddFormItem(loopbackField)

	form.form.AddButton("Submit", nil)
	form.form.AddButton("Cancel", nil)

	formFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form.form, 11, 1, true).
		AddItem(hintBox, 8, 1, false)

	form.Flex.AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(formFlex, 0, 1, true).
			AddItem(nil, 0, 1, false),
			0, 1, true).
		AddItem(nil, 0, 1, false)

	return form
}

func (form *AddRouteForm) GetID() string {
	return "addroute_form"
}

func (form *AddRouteForm) SetSubmitFunc(f func(string, int, bool)) {
	btnId := form.form.GetButtonIndex("Submit")
	submitBtn := form.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(func() {
		f(add_route_cidr.Last, add_route_metric.Last, add_route_loopback.Last)
	})
}

func (form *AddRouteForm) SetCancelFunc(f func()) {
	btnId := form.form.GetButtonIndex("Cancel")
	submitBtn := form.form.GetButton(btnId)
	submitBtn.SetSelectedFunc(f)
}
