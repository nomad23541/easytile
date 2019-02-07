package main

import (
	"log"
	"math"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type RectIndex struct {
	Rect xrect.Rect
	Index int
}

func main() {
	// connect to the xgbutil's X server
	X, err := xgbutil.NewConn()
	assert(err)

	// also connect to xgb's X server, since they xgbutil can't be used in xgb's functions
	Xg, err := xgb.NewConn()
	assert(err)

	// initialize keybind
	keybind.Initialize(X)

	// right half of screen
	keypress(Xg, X, err, "Control-Shift-KP_6", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, (currentHead.X() + currentHead.Width()) - (currentHead.Width() / 2), currentHead.Y(), (currentHead.Width() / 2), currentHead.Height())
		assert(err)
	})

	// left half of screen
	keypress(Xg, X, err, "Control-Shift-KP_4", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, currentHead.X(), currentHead.Y(), (currentHead.Width() / 2), currentHead.Height())
		assert(err)
	})

	// top half of screen
	keypress(Xg, X, err, "Control-Shift-KP_8", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, currentHead.X(), currentHead.Y(), currentHead.Width(), (currentHead.Height() / 2))
		assert(err)
	})

	// bottom half of screen
	keypress(Xg, X, err, "Control-Shift-KP_2", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, currentHead.X(), (currentHead.Height() / 2) + currentHead.Y(), currentHead.Width(), (currentHead.Height() / 2))
		assert(err)
	})

	// left top corner
	keypress(Xg, X, err, "Control-Shift-KP_7", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, currentHead.X(), currentHead.Y(), (currentHead.Width() / 2), (currentHead.Height() / 2))
		assert(err)
	})

	// left bottom corner
	keypress(Xg, X, err, "Control-Shift-KP_1", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, currentHead.X(), (currentHead.Height() / 2) + currentHead.Y(), (currentHead.Width() / 2), (currentHead.Height() / 2))
		assert(err)
	})

	// top right corner
	keypress(Xg, X, err, "Control-Shift-KP_9", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, (currentHead.X() + currentHead.Width()) - (currentHead.Width() / 2), currentHead.Y(), (currentHead.Width() / 2), (currentHead.Height() / 2))
		assert(err)
	})

	// bottom right corner
	keypress(Xg, X, err, "Control-Shift-KP_3", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, (currentHead.X() + currentHead.Width()) - (currentHead.Width() / 2), (currentHead.Height() / 2) + currentHead.Y(), (currentHead.Width() / 2), (currentHead.Height() / 2))
		assert(err)
	})

	// maximize window
	keypress(Xg, X, err, "Control-Shift-KP_5", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		err = ewmh.MoveresizeWindow(X, window, currentHead.X(), currentHead.Y(), currentHead.Width(), currentHead.Height())
		assert(err)
	})

	// move window to next head
	keypress(Xg, X, err, "Control-Shift-KP_Enter", true, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {})

	// FOR DEBUGGING, exit application
	/*
	keypress(Xg, X, err, "Control-Shift-KP_0", false, func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect) {
		os.Exit(0)
	})
	*/

	log.Println("easytile started")
	xevent.Main(X)
}

func keypress(Xg *xgb.Conn, X *xgbutil.XUtil, err error, combo string, toNewHead bool, fn func(X *xgbutil.XUtil, window xproto.Window, currentHead xrect.Rect)) {
	cb := keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
		// get active window id
		windowId, err := ewmh.ActiveWindowGet(X)
		assert(err)

		// active xproto window
		window := xproto.Window(windowId)

		// remove maxmized states from window, resize won't work if fullscreen
		err = ewmh.WmStateReq(X, window, 0, "_NET_WM_STATE_MAXIMIZED_VERT")
		assert(err)
		err = ewmh.WmStateReq(X, window, 0, "_NET_WM_STATE_MAXIMIZED_HORZ")
		assert(err)

		// wrap this window into a new window to get it's width & height
		wrap := xwindow.New(X, window)
		geom, err := wrap.Geometry()
		assert(err)
		decX, decY, width, height := int(geom.X() + 1), int(geom.Y() + 1), int(geom.Width()), int(geom.Height())

		// now get this window's x & y position relative to the root window
		coord, err := xproto.TranslateCoordinates(Xg, window, X.RootWin(), 0, 0).Reply()
		assert(err)
		x, y := int(coord.DstX), int(coord.DstY)

		// compare X coord to all of the heads and determine closest one
		var currentHead xrect.Rect
		var currentIndex int
		heads := getStruts(X, err)
		viableHeads := []RectIndex{}
		for i, head := range heads {
			// determine which head the active window is in
			if x < head.X() + head.Width() && x + width > head.X() && y < head.Y() + head.Height() && y + height > head.Y() {
				viableHeads = append(viableHeads, RectIndex{head, i})
			}
		}

		// get the closest head by comparing the two head's center x and the window's center x
		centerOfWindow := (width / 2) + x
		distance := math.Abs(float64((viableHeads[0].Rect.Width() / 2) + viableHeads[0].Rect.X()) - float64(centerOfWindow))
		var closestIndex int
		for i := range viableHeads {
			idistance := math.Abs(float64((viableHeads[i].Rect.Width() / 2) + viableHeads[i].Rect.X()) - float64(centerOfWindow))
			if(idistance < distance) {
				closestIndex = i
				distance = idistance
			}
		}

		currentHead = viableHeads[closestIndex].Rect
		currentIndex = viableHeads[closestIndex].Index

		// update currentHead's width & height to offset the wm's decorations
		currentHead.WidthSet(currentHead.Width() - decX)
		currentHead.HeightSet(currentHead.Height() - decY)

		// if true, move to new head instead of moving/resizing
		if toNewHead {
			nextIndex := currentIndex + 1
			if nextIndex >= heads.Len() {
				nextIndex = 0
			}

			nextHead := heads[nextIndex]
			err = ewmh.MoveresizeWindow(X, window, nextHead.X(), nextHead.Y(), width, height)
			assert(err)
		}

		if currentHead == nil {
			log.Println("Could not find a current head")
		} else if !toNewHead {
			fn(X, window, currentHead)
		}
	})

	err = cb.Connect(X, X.RootWin(), combo, true)
	assert(err)
}
