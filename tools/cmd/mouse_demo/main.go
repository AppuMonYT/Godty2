// License: GPLv3 Copyright: 2023, Kovid Goyal, <kovid at kovidgoyal.net>

package mouse_demo

import (
	"fmt"
	"strconv"

	"github.com/kovidgoyal/kitty/tools/tui/loop"
)

var _ = fmt.Print

func Run(args []string) (rc int, err error) {
	all_pointer_shapes := []loop.PointerShape{
		// start all pointer shapes (auto generated by gen-key-constants.py do not edit)
		loop.DEFAULT_POINTER,
		loop.TEXT_POINTER,
		loop.POINTER_POINTER,
		loop.HELP_POINTER,
		loop.WAIT_POINTER,
		loop.PROGRESS_POINTER,
		loop.CROSSHAIR_POINTER,
		loop.CELL_POINTER,
		loop.VERTICAL_TEXT_POINTER,
		loop.MOVE_POINTER,
		loop.E_RESIZE_POINTER,
		loop.NE_RESIZE_POINTER,
		loop.NW_RESIZE_POINTER,
		loop.N_RESIZE_POINTER,
		loop.SE_RESIZE_POINTER,
		loop.SW_RESIZE_POINTER,
		loop.S_RESIZE_POINTER,
		loop.W_RESIZE_POINTER,
		loop.EW_RESIZE_POINTER,
		loop.NS_RESIZE_POINTER,
		loop.NESW_RESIZE_POINTER,
		loop.NWSE_RESIZE_POINTER,
		loop.ZOOM_IN_POINTER,
		loop.ZOOM_OUT_POINTER,
		loop.ALIAS_POINTER,
		loop.COPY_POINTER,
		loop.NOT_ALLOWED_POINTER,
		loop.NO_DROP_POINTER,
		loop.GRAB_POINTER,
		loop.GRABBING_POINTER,
		// end all pointer shapes
	}
	all_pointer_shape_names := make([]string, len(all_pointer_shapes))
	col_width := 0
	for i, p := range all_pointer_shapes {
		all_pointer_shape_names[i] = p.String()
		col_width = max(col_width, len(all_pointer_shape_names[i]))
	}
	col_width += 1

	lp, err := loop.New()
	if err != nil {
		return 1, err
	}
	lp.MouseTrackingMode(loop.FULL_MOUSE_TRACKING)
	var current_mouse_event *loop.MouseEvent

	draw_screen := func() {
		lp.StartAtomicUpdate()
		defer lp.EndAtomicUpdate()
		lp.AllowLineWrapping(false)
		defer lp.AllowLineWrapping(true)
		if current_mouse_event == nil {
			lp.ClearScreen()
			lp.Println(`Move the mouse or click to see mouse events`)
			return
		}
		lp.ClearScreen()
		if current_mouse_event.Event_type == loop.MOUSE_LEAVE {
			lp.Println("Mouse has left the window")
			return
		}
		lp.Printf("Position: %d, %d (pixels)\r\n", current_mouse_event.Pixel.X, current_mouse_event.Pixel.Y)
		lp.Printf("Cell    : %d, %d\r\n", current_mouse_event.Cell.X, current_mouse_event.Cell.Y)
		lp.Printf("Type    : %s\r\n", current_mouse_event.Event_type)
		y := 3
		if current_mouse_event.Buttons != loop.NO_MOUSE_BUTTON {
			lp.Println(current_mouse_event.Buttons.String())
			y += 1
		}
		if mods := current_mouse_event.Mods.String(); mods != "" {
			lp.Printf("Modifiers: %s\r\n", mods)
			y += 1
		}
		lp.Println("Hover the mouse over the names below to see the shapes")
		y += 1

		sw := 80
		sh := 24
		if s, err := lp.ScreenSize(); err == nil {
			sw = int(s.WidthCells)
			sh = int(s.HeightCells)
		}

		num_cols := max(1, sw/col_width)
		pos := 0
		colfmt := "%-" + strconv.Itoa(col_width) + "s"
		is_on_name := false
		var ps loop.PointerShape
		for y < sh && pos < len(all_pointer_shapes) {
			is_row := y == current_mouse_event.Cell.Y
			for c := 0; c < num_cols && pos < len(all_pointer_shapes); c++ {
				name := all_pointer_shape_names[pos]
				is_hovered := false
				if is_row {
					start_x := c * col_width
					x := current_mouse_event.Cell.X
					if x < start_x+len(name) && x >= start_x {
						is_on_name = true
						is_hovered = true
						ps = all_pointer_shapes[pos]
					}
				}
				if is_hovered {
					lp.QueueWriteString("\x1b[31m")
				}
				lp.Printf(colfmt, name)
				lp.QueueWriteString("\x1b[m")
				pos++
			}
			y += 1
			lp.Println()
		}
		lp.PopPointerShape()
		if is_on_name {
			lp.PushPointerShape(ps)
		}
	}

	lp.OnInitialize = func() (string, error) {
		lp.SetWindowTitle("kitty mouse features demo")
		lp.SetCursorVisible(false)
		draw_screen()
		return "", nil
	}
	lp.OnFinalize = func() string {
		lp.SetCursorVisible(true)
		return ""
	}

	lp.OnMouseEvent = func(ev *loop.MouseEvent) error {
		current_mouse_event = ev
		draw_screen()
		return nil
	}
	lp.OnKeyEvent = func(ev *loop.KeyEvent) error {
		if ev.MatchesPressOrRepeat("esc") || ev.MatchesPressOrRepeat("ctrl+c") {
			lp.Quit(0)
		}
		return nil
	}
	lp.OnResize = func(old_size loop.ScreenSize, new_size loop.ScreenSize) error {
		draw_screen()
		return nil
	}
	err = lp.Run()
	if err != nil {
		rc = 1
	}
	return
}
