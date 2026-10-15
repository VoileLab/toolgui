package tgwails

import (
	goruntime "runtime"

	"github.com/voilelab/toolgui/toolgui/tgframe"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
)

// nativeMenu translates the tree [tgframe.App.SetMenu] declared into the one
// Wails hangs off the window. Same declaration the web menubar draws from:
// only the thing drawing it changes.
//
// It returns nil where there is nothing to draw -- an app with no menu, on a
// platform that adds no standard ones -- leaving the window without a
// menubar. An optional fullscreen callback adds the built-in View menu.
func nativeMenu(frameless bool, nodes []*tgframe.MenuNode,
	click func(id string), fullscreen ...func()) *menu.Menu {

	return buildMenu(goruntime.GOOS, frameless, nodes, click, fullscreen...)
}

// buildMenu is [nativeMenu] with the platform named, so a test can build the
// menu of a platform it is not running on.
func buildMenu(goos string, frameless bool, nodes []*tgframe.MenuNode,
	click func(id string), fullscreen ...func()) *menu.Menu {

	var toggle func()
	if len(fullscreen) > 0 {
		toggle = fullscreen[0]
	}
	roles := roleItems(goos, frameless)
	if len(roles) == 0 && len(nodes) == 0 && toggle == nil {
		return nil
	}

	native := menu.NewMenu()
	for _, item := range roles {
		// Wails' macOS Window role includes its own fullscreen command with
		// Ctrl+Cmd+F. Omit that role only when supplying our runtime toggle.
		if toggle != nil && item.Role == menu.WindowMenuRole {
			continue
		}
		native.Append(item)
	}

	appendNodes(native, nodes, click)

	if toggle != nil {
		shortcut := keys.Key("f11")
		if goos == "darwin" {
			shortcut = keys.Combo("f", keys.ControlKey, keys.CmdOrCtrlKey)
		}
		// Application shortcuts take precedence. Keep the command clickable
		// without registering a second action for an existing accelerator.
		if fullscreenShortcutTaken(native, goos) {
			shortcut = nil
		}
		var view *menu.Menu
		for _, item := range native.Items {
			if item.Type == menu.SubmenuType && item.Label == "View" {
				view = item.SubMenu
				break
			}
		}
		if view == nil {
			view = native.AddSubmenu("View")
		}
		view.AddText("Toggle Fullscreen", shortcut,
			func(*menu.CallbackData) { toggle() })
	}

	return native
}

// fullscreenShortcutTaken searches all app submenus, comparing native
// modifiers on macOS so declaration order does not affect collision checks.
func fullscreenShortcutTaken(native *menu.Menu, goos string) bool {
	for _, item := range native.Items {
		if a := item.Accelerator; a != nil {
			if goos == "darwin" {
				if a.Key == "f" && keys.ToMacModifier(a) == keys.NSEventModifierFlagControl|keys.NSEventModifierFlagCommand {
					return true
				}
			} else if a.Key == "f11" && len(a.Modifiers) == 0 {
				return true
			}
		}
		if item.SubMenu != nil && fullscreenShortcutTaken(item.SubMenu, goos) {
			return true
		}
	}
	return false
}

// roleItems returns the standard menus that go in front of the app's own
// entries.
//
// Roles are macOS only -- Wails v2 says so, and the constants for them are
// commented out everywhere else -- so nothing goes in front on the other
// platforms. macOS needs them: Cmd+C, Cmd+V and Cmd+Q are actions of the
// standard App and Edit menus rather than of the webview, so a window whose
// menu replaces them has none of those shortcuts.
//
// This is what Wails itself puts in when no menu is set, Frameless rule
// included, so declaring a menu does not cost the window a standard one.
func roleItems(goos string, frameless bool) []*menu.MenuItem {
	if goos != "darwin" {
		return nil
	}

	items := []*menu.MenuItem{menu.AppMenu(), menu.EditMenu()}
	if !frameless {
		// The Window menu's entries do nothing on a frameless window.
		items = append(items, menu.WindowMenu())
	}

	return items
}

// appendNodes writes nodes into native, recursing through submenus.
//
// Text, separator and submenu are the three kinds a [tgframe.Menu] can hold,
// and SetMenu turns away a tree carrying anything else, so an unknown node is
// nothing to report on here -- it is dropped.
func appendNodes(native *menu.Menu, nodes []*tgframe.MenuNode,
	click func(id string)) {

	for _, node := range nodes {
		switch node.Type {
		case tgframe.MenuNodeSeparator:
			native.AddSeparator()
		case tgframe.MenuNodeText:
			// The id is the one the app declared, prefix and all, so the run
			// handling the click reads it back with tgframe.MenuClicked just
			// as it does the web menubar's.
			native.AddText(node.Label, accelerator(node.Accelerator),
				func(*menu.CallbackData) {
					click(node.ID)
				})
		case tgframe.MenuNodeSubmenu:
			appendNodes(native.AddSubmenu(node.Label), node.Children, click)
		}
	}
}

// accelerator translates the combination the item declared into the one Wails
// hangs off the menu item, from where the OS dispatches it. This is the whole
// of the desktop's half of the feature: nothing here listens for a keystroke.
//
// The spelling tgframe normalizes to is one Wails parses, so a parse that
// fails is an accelerator tgframe let through and this build does not know --
// the item still belongs in the menu, so it goes in without one rather than
// taking the menu down with it.
func accelerator(accel string) *keys.Accelerator {
	if accel == "" {
		return nil
	}

	parsed, err := keys.Parse(accel)
	if err != nil {
		return nil
	}

	return parsed
}
