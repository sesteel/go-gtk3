package main

import "github.com/norisatir/go-gtk3/gtk3"
import "github.com/norisatir/go-gtk3/gobject"
import "time"

// Bug struct
type Bug struct {
	Fixed bool
	Number uint
	Severity string
	Description string
}

// Create slice of bugs (initialized in init func)
var Bugs []Bug

// Column constants
const (
	ColumnFixed = iota
	ColumnNumber
	ColumnSeverity
	ColumnDescription
	ColumnPulse
    ColumnIcon
	ColumnActive
	ColumnSensitive
	NumColumns
)

func CreateModel() gtk3.TreeModelLike {
	store := gtk3.NewListStore([]gobject.GType{gobject.G_TYPE_BOOLEAN,
											   gobject.G_TYPE_UINT,
											   gobject.G_TYPE_STRING,
											   gobject.G_TYPE_STRING,
											   gobject.G_TYPE_UINT,
											   gobject.G_TYPE_STRING,
											   gobject.G_TYPE_BOOLEAN,
											   gobject.G_TYPE_BOOLEAN})

	// Add data to list store
	var icon_name string
	var sensitive bool

	var iter gtk3.TreeIter

	for i,bug := range Bugs {
		if i == 1 || i == 3 {
			icon_name = "battery-caution-charging-symbolic"
		} else {
			icon_name = ""
		}

		if i == 3 {
			sensitive = false
		} else {
			sensitive = true
		}

		store.Append(&iter)
		store.SetValues(&iter,gtk3.V{ColumnFixed : bug.Fixed,
									ColumnNumber : bug.Number,
									ColumnSeverity : bug.Severity,
									ColumnDescription : bug.Description,
									ColumnPulse : uint(0),
									ColumnIcon : icon_name,
									ColumnActive : false,
									ColumnSensitive : sensitive })

	}

	return store
}

func FixedToggled(m gtk3.TreeModelLike, data ...interface{}) {
	// Data should hold additional arguments filled with closure
	if len(data) != 2 {
		return
	}

	// We need path string
	pathStr, _ := data[1].(string)

	model := m.(*gtk3.ListStore)

	var iter gtk3.TreeIter
	path := gtk3.NewTreePathFromString(pathStr)

	// Get toggled iter
	model.GetIter(&iter, path)
	fixed := model.GetValue(&iter, ColumnFixed).(bool)

	// Do something with the value
	if fixed {
		fixed = false
	} else {
		fixed = true
	}

	// set new value
	model.SetValue(&iter, ColumnFixed, fixed)

	// C code would have to free path
	// Go GC will do that for us - weeeee
}

func AddColumns(treeView *gtk3.TreeView) {
	model := treeView.GetModel()

	var renderer gtk3.CellRendererLike

	// column for fixed toggles
	renderer = gtk3.NewCellRendererToggle()
	renderer.CRenderer().Connect("toggled", FixedToggled, model)

	column := gtk3.NewTreeViewColumnWithAttributes("Fixed?", renderer, gtk3.A{{"active", ColumnFixed}})
	
	// Set this column to a fixed sizing (50 pixels)
	column.SetSizing(gtk3.GtkTreeViewColumnSizing.FIXED)
	column.SetFixedWidth(50)
	treeView.AppendColumn(column)

	// Column for bug numbers
	renderer = gtk3.NewCellRendererText()
	column = gtk3.NewTreeViewColumnWithAttributes("Bug number", renderer, gtk3.A{{"text", ColumnNumber}})
	column.SetSortColumnId(ColumnNumber)
	treeView.AppendColumn(column)

	// Column for severities
	renderer = gtk3.NewCellRendererText()
	column = gtk3.NewTreeViewColumnWithAttributes("Severity", renderer, gtk3.A{{"text", ColumnSeverity}})
	column.SetSortColumnId(ColumnSeverity)
	treeView.AppendColumn(column)

	// Column for description
	renderer = gtk3.NewCellRendererText()
	column = gtk3.NewTreeViewColumnWithAttributes("Description", renderer, gtk3.A{{"text", ColumnDescription}})
	column.SetSortColumnId(ColumnDescription)
	treeView.AppendColumn(column)

	// Column for spinner
	renderer = gtk3.NewCellRendererSpinner()
	column = gtk3.NewTreeViewColumnWithAttributes("Spinning", renderer, gtk3.A{{"pulse", ColumnPulse}, {"active", ColumnActive}})
	column.SetSortColumnId(ColumnPulse)
	treeView.AppendColumn(column)

	// Column for symbolic icon
	renderer = gtk3.NewCellRendererPixbuf()
	renderer.CRenderer().Set(gtk3.P{"follow-state" : true})
	column = gtk3.NewTreeViewColumnWithAttributes("Symbolic icon", renderer,
		gtk3.A{{"icon-name", ColumnIcon}, {"sensitive", ColumnSensitive}})
	column.SetSortColumnId(ColumnIcon)
	treeView.AppendColumn(column)
}

func init() {
	Bugs = []Bug{{false, 60482, "Normal", "scrollable notebooks and hidden tabs"},
				  { false, 60620, "Critical",   "gdk_window_clear_area (gdkwindow-win32.c) is not thread-safe" },
				  { false, 50214, "Major",      "Xft support does not clean up correctly" },
				  { true,  52877, "Major",      "GtkFileSelection needs a refresh method. " },
				  { false, 56070, "Normal",     "Can't click button after setting in sensitive" },
				  { true,  56355, "Normal",     "GtkLabel - Not all changes propagate correctly" },
				  { false, 50055, "Normal",     "Rework width/height computations for TreeView" },
				  { true, 58278, "Normal",     "gtk_dialog_set_response_sensitive () doesn't work" },
				  { false, 55767, "Normal",     "Getters for all setters" },
				  { true, 56925, "Normal",     "Gtkcalender size" },
				  { false, 56221, "Normal",     "Selectable label needs right-click copy menu" },
				  { true,  50939, "Normal",     "Add shift clicking to GtkTextView" },
				  { false, 6112,  "Enhancement","netscape-like collapsable toolbars" },
				  { false, 1,     "Normal",     "First bug :=)" }}
}


func main() {
	gtk3.Init()
	// Create Window
	window := gtk3.NewWindow(gtk3.GtkWindowType.TOPLEVEL)
	window.SetTitle("GtkListStore demo")
	window.Connect("destroy", gtk3.GoMainQuit)
	window.SetBorderWidth(8)

	vbox := gtk3.NewVBox(8)
	window.Add(vbox)

	label := gtk3.NewLabel("This is the bug list (note: not based on real data...)")
	vbox.PackStart(label, false, false, 0)

	sw := gtk3.NewScrolledWindow(nil, nil)
	sw.SetShadowType(gtk3.GtkShadow.ETCHED_IN)
	sw.SetPolicy(gtk3.GtkPolicy.NEVER, gtk3.GtkPolicy.AUTOMATIC)
	vbox.PackStart(sw, true, true, 0)

	// create tree model
	model := CreateModel()

	// create tree view
	treeview := gtk3.NewTreeViewWithModel(model)
	treeview.SetRulesHint(true)
	treeview.SetSearchColumn(ColumnDescription)
	sw.Add(treeview)

	// Add columns to treeview
	AddColumns(treeview)

	// Make spinner spin
	t := time.NewTicker(80000000)
	// Spinner func
	go func(timeout <-chan int64) {
		for _ = range timeout {
			gtk3.FreezeMain.L.Lock()
			var sIter gtk3.TreeIter
			m := model.(*gtk3.ListStore)
			m.GetIterFirst(&sIter)
			pulse := m.GetValue(&sIter, ColumnPulse).(uint)
			
			if pulse == ^uint(0) {
				pulse = 0
			} else {
				pulse++
			}

			m.SetValues(&sIter, gtk3.V{ColumnPulse : pulse, ColumnActive : true})
			gtk3.FreezeMain.L.Unlock()
		}
	}(t.C)

	window.SetDefaultSize(280, 250)
	window.ShowAll()

	gtk3.GoMain(false)
	t.Stop()
}
