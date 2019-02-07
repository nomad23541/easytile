package main

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xinerama"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func getStruts(X *xgbutil.XUtil, err error) xinerama.Heads {
	// get root window and wrap in a Window Type
	root := xwindow.New(X, X.RootWin())

	// get Geometry of the root window
	rootGeom, err := root.Geometry()
	assert(err)

	// get each rectangle head (each specific monitor)
	// first check if xinerama is enable, if not use root geometry
	var heads xinerama.Heads
	if X.ExtInitialized("XINERAMA") {
		heads, err = xinerama.PhysicalHeads(X)
		assert(err)
	} else {
		heads = xinerama.Heads{rootGeom}
	}

	// done, now get any top-level windows and subtract that from the head geometry (panels)
	clients, err := ewmh.ClientListGet(X)
	assert(err)

	for _, clientId := range clients {
		strut, err := ewmh.WmStrutPartialGet(X, clientId)
		// if there are errors, then there are no struts for said client
		if err != nil {
			continue
		}

		// modify the head
		xrect.ApplyStrut(heads, uint(rootGeom.Width()), uint(rootGeom.Height()),
			strut.Left, strut.Right, strut.Top, strut.Bottom,
			strut.LeftStartY, strut.LeftEndY,
			strut.RightStartY, strut.RightEndY,
			strut.TopStartX, strut.TopEndX,
			strut.BottomStartX, strut.BottomEndX)
	}

	// finally return xrect.Rect array
	return heads
}