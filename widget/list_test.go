package widget

import (
	"image/color"
	"testing"

	"github.com/blizzy78/ebitenui/event"

	"github.com/matryer/is"
)

func TestList_SelectedEntry_Initial(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	list := newList(t,
		ListOpts.WithEntries(entries),

		ListOpts.WithEntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),

		ListOpts.WithEntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			is.Fail() // event fired without previous action
		}))

	is.Equal(list.SelectedEntry(), nil)
}

func TestList_SetSelectedEntry(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.WithEntries(entries),

		ListOpts.WithEntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),

		ListOpts.WithEntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	list.SetSelectedEntry(entries[1])
	event.FireDeferredEvents()

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	list.SetSelectedEntry(entries[1])
	event.FireDeferredEvents()

	is.Equal(numEvents, 1)
}

func TestList_EntrySelectedEvent_User(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.WithEntries(entries),

		ListOpts.WithEntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),

		ListOpts.WithEntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(numEvents, 1)
}

func TestList_EntrySelectedEvent_User_AllowReselect(t *testing.T) {
	is := is.New(t)

	entries := []interface{}{"first", "second", "third"}

	var eventArgs *ListEntrySelectedEventArgs
	numEvents := 0

	list := newList(t,
		ListOpts.WithEntries(entries),

		ListOpts.WithEntryLabelFunc(func(e interface{}) string {
			return e.(string)
		}),

		ListOpts.WithAllowReselect(),

		ListOpts.WithEntrySelectedHandler(func(args *ListEntrySelectedEventArgs) {
			eventArgs = args
			numEvents++
		}))

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	leftMouseButtonClick(list.buttons[1], t)

	is.Equal(eventArgs.Entry, entries[1])
	is.Equal(eventArgs.PreviousEntry, entries[1])
	is.Equal(list.SelectedEntry(), entries[1])

	is.Equal(numEvents, 2)
}

func newList(t *testing.T, opts ...ListOpt) *List {
	t.Helper()

	l := NewList(
		append(opts, []ListOpt{
			ListOpts.WithImage(&ScrollContainerImage{
				Idle:     newNineSliceEmpty(t),
				Disabled: newNineSliceEmpty(t),
				Mask:     newNineSliceEmpty(t),
			}),

			ListOpts.WithEntryFontFace(loadFont(t)),

			ListOpts.WithEntryColor(&ListEntryColor{
				Unselected:                   color.Transparent,
				Selected:                     color.Transparent,
				DisabledUnselected:           color.Transparent,
				DisabledSelected:             color.Transparent,
				UnselectedBackground:         color.Transparent,
				SelectedBackground:           color.Transparent,
				DisabledUnselectedBackground: color.Transparent,
				DisabledSelectedBackground:   color.Transparent,
			}),
		}...)...)

	event.FireDeferredEvents()
	render(l, t)
	return l
}

func listEntryButtons(l *List) []*Button {
	return l.buttons
}
