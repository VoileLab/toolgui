//go:build js && wasm

package tgwasm

import (
	"strconv"
	"syscall/js"

	"github.com/voilelab/toolgui/toolgui/tgframe"
)

type folderPick struct {
	session *tgframe.Session
	id      string
	write   bool
}

// Promises run off the JS callback stack. Validation and applying a grant can
// wait for a running page; blocking the callback would stop its awaited IO.
func (b *bridge) folderPromise(fn func() (any, error)) js.Value {
	executor := js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve, reject := args[0], args[1]
		b.runs.do(func() {
			value, err := fn()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New(err.Error()))
			} else {
				resolve.Invoke(value)
			}
		})
		return nil
	})
	defer executor.Release()
	return js.Global().Get("Promise").New(executor)
}

func (b *bridge) jsBeginFolderPick(_ js.Value, args []js.Value) any {
	session := b.host.Session()
	return b.folderPromise(func() (any, error) {
		if session == nil {
			return nil, ErrNoSession
		}
		if b.host.Session() != session {
			return nil, tgframe.ErrSessionReplaced
		}
		if len(args) < 2 {
			return nil, tgframe.ErrNotOnPage
		}
		id := args[0].String()
		if err := session.ValidateFolderPick(id); err != nil {
			return nil, err
		}
		b.nextFolderPick++
		token := strconv.FormatUint(b.nextFolderPick, 10)
		b.folderPicks[token] = folderPick{session: session, id: id, write: args[1].Bool()}
		return token, nil
	})
}

func (b *bridge) jsCancelFolderPick(_ js.Value, args []js.Value) any {
	if len(args) > 0 {
		delete(b.folderPicks, args[0].String())
	}
	return nil
}

func (b *bridge) jsFolderPicked(_ js.Value, args []js.Value) any {
	return b.folderPromise(func() (any, error) {
		if len(args) < 2 {
			return nil, tgframe.ErrNotOnPage
		}
		token := args[0].String()
		pick, ok := b.folderPicks[token]
		delete(b.folderPicks, token)
		if !ok {
			return nil, tgframe.ErrSessionReplaced
		}
		handle := args[1]
		if handle.Type() != js.TypeObject || handle.IsNull() || handle.Get("kind").String() != "directory" {
			return nil, tgframe.ErrNotOnPage
		}
		// Only a complete, successful grant enters state, and only its old session.
		return nil, b.host.HandleFolder(pick.session, &tgframe.EventFolder{
			ID: pick.id, Folder: &browserFolder{root: handle, write: pick.write},
		})
	})
}
