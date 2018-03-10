// +build darwin,amd64

package mac

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
#cgo LDFLAGS: -framework WebKit
#cgo LDFLAGS: -framework CoreImage
#cgo LDFLAGS: -framework Security
#include "bridge.h"
*/
import "C"
import (
	"net/url"
	"unsafe"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/bridge"
	"github.com/pkg/errors"
)

func handleMacOSRequest(rawurl string, p bridge.Payload, returnID string) (res bridge.Payload, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Wrap(err, "handling MacOS request failed"))
	}

	if len(returnID) != 0 {
		q := u.Query()
		q.Set("return-id", returnID)
		u.RawQuery = q.Encode()
	}

	var bres C.bridge_result
	if p == nil {
		bres = C.macosRequest(C.CString(u.String()), nil)
	} else {
		bres = C.macosRequest(C.CString(u.String()), C.CString(p.String()))
	}
	return parseBridgeResult(bres)
}

func parseBridgeResult(res C.bridge_result) (p bridge.Payload, err error) {
	if res.payload != nil {
		p = bridge.PayloadFromString(C.GoString(res.payload))
		C.free(unsafe.Pointer(res.payload))
	}

	if res.err != nil {
		err = errors.Errorf("handling MacOS request failed: %s", C.GoString(res.err))
		C.free(unsafe.Pointer(res.err))
	}
	return p, err
}

//export macosRequestResult
func macosRequestResult(rawretID *C.char, res C.bridge_result) {
	retID := C.GoString(rawretID)
	C.free(unsafe.Pointer(rawretID))

	payload, err := parseBridgeResult(res)
	driver.macos.Return(retID, payload, err)
}

//export goRequest
func goRequest(url *C.char, payload *C.char) {
	driver.golang.Request(
		C.GoString(url),
		bridge.PayloadFromString(C.GoString(payload)),
	)
}

//export goRequestWithResult
// res should be free after each call of goRequestWithResult.
func goRequestWithResult(url *C.char, payload *C.char) (res *C.char) {
	pret := driver.golang.RequestWithResponse(
		C.GoString(url),
		bridge.PayloadFromString(C.GoString(payload)),
	)

	if pret != nil {
		res = C.CString(pret.String())
	}
	return res
}

func windowHandler(h func(w *Window, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating window handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			return nil
		}

		win, ok := elem.(app.Window)
		if !ok {
			panic(errors.Errorf("creating window handler failed: element with id %v is not a window", id))
		}
		return h(win.Base().(*Window), u, p)
	}
}

func menuHandler(h func(m *Menu, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating menu handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			panic(errors.Wrap(err, "creating menu handler failed"))
		}

		menu, ok := elem.(app.Menu)
		if !ok {
			panic(errors.Errorf("creating menu handler failed: element with id %v is not a menu", id))
		}
		return h(menu.Base().(*Menu), u, p)
	}
}

func filePanelHandler(h func(panel *FilePanel, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating file panel handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			panic(errors.Wrap(err, "creating file panel handler failed"))
		}

		panel, ok := elem.(*FilePanel)
		if !ok {
			panic(errors.Errorf("creating file panel handler failed: element with id %v is not a file panel", id))
		}

		return h(panel, u, p)
	}
}

func saveFilePanelHandler(h func(panel *SaveFilePanel, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating save file panel handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			panic(errors.Wrap(err, "creating save file panel handler failed"))
		}

		panel, ok := elem.(*SaveFilePanel)
		if !ok {
			panic(errors.Errorf("creating save file panel handler failed: element with id %v is not a file panel", id))
		}

		return h(panel, u, p)
	}
}

func notificationHandler(h func(n *Notification, u *url.URL, p bridge.Payload) (res bridge.Payload)) bridge.GoHandler {
	return func(u *url.URL, p bridge.Payload) (res bridge.Payload) {
		id, err := uuid.Parse(u.Query().Get("id"))
		if err != nil {
			panic(errors.Wrap(err, "creating notification handler failed"))
		}

		var elem app.Element
		if elem, err = driver.elements.Element(id); err != nil {
			return nil
			// panic(errors.Wrap(err, "creating notification handler failed"))
		}

		notification, ok := elem.(*Notification)
		if !ok {
			panic(errors.Errorf("creating notification handler failed: element with id %v is not notification", id))
		}

		return h(notification, u, p)
	}
}

func macCall(call string) error {
	ccall := C.CString(call)
	C.macCall(ccall)
	C.free(unsafe.Pointer(ccall))
	return nil
}

//export macCallReturn
func macCallReturn(retID, ret, err *C.char) {
	driver.macRPC.Return(
		C.GoString(retID),
		C.GoString(ret),
		C.GoString(err),
	)
}
