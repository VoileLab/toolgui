package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"syscall"
	"unicode"
)

const desktopImport = "github.com/voilelab/toolgui/toolgui-wails"

type goEnvironment struct {
	GOOS, GOARCH, CGO_ENABLED, GOFLAGS string
}

func buildOrRun(mode string, opts buildOptions, stdout, stderr io.Writer) error {
	ctx := context.Background()
	if mode == "run" {
		var stop context.CancelFunc
		ctx, stop = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
		defer stop()
	}
	var env goEnvironment
	if err := goJSON(ctx, &env, "env", "-json", "GOOS", "GOARCH", "CGO_ENABLED", "GOFLAGS"); err != nil {
		return err
	}
	if err := checkTarget(env, runtime.GOOS, runtime.GOARCH); err != nil {
		return err
	}
	tags, ldflags, err := compilerFlags(opts, env, pkgConfigExists)
	if err != nil {
		return err
	}
	app, desktopDir, err := findApp(ctx, opts.pkg, tags)
	if err != nil {
		return err
	}
	if err := checkAssets(desktopDir); err != nil {
		return err
	}
	out := opts.out
	if mode == "run" {
		dir, err := os.MkdirTemp("", "toolgui-desktop-*")
		if err != nil {
			return fmt.Errorf("create temporary build directory: %w", err)
		}
		defer os.RemoveAll(dir)
		out = filepath.Join(dir, executableName(app.ImportPath, env.GOOS))
	} else if out == "" {
		out = filepath.Join("build", "bin", executableName(app.ImportPath, env.GOOS))
	}
	out, err = filepath.Abs(out)
	if err != nil {
		return err
	}
	if stat, err := os.Stat(out); err == nil && stat.IsDir() {
		return fmt.Errorf("output %s is a directory; -o must name an executable", out)
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	args := []string{"build", "-tags", strings.Join(tags, ","), "-o", out}
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	args = append(args, opts.pkg)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = compilerEnvironment(os.Environ(), env.GOOS)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}
	if mode == "build" {
		fmt.Fprintln(stdout, out)
		return nil
	}
	cmd = exec.CommandContext(ctx, out, opts.appArgs...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, stdout, stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("app: %w", err)
	}
	return nil
}

func checkTarget(env goEnvironment, hostOS, hostArch string) error {
	if env.GOOS != hostOS || env.GOARCH != hostArch {
		return fmt.Errorf("cross compilation is not supported: target %s/%s, host %s/%s; unset GOOS/GOARCH or build on the target machine", env.GOOS, env.GOARCH, hostOS, hostArch)
	}
	if !slices.Contains([]string{"darwin", "linux", "windows"}, env.GOOS) {
		return fmt.Errorf("unsupported desktop OS %q", env.GOOS)
	}
	if env.GOOS != "windows" && env.CGO_ENABLED != "1" {
		return errors.New("desktop builds on macOS and Linux need cgo; set CGO_ENABLED=1 and install a C compiler and the webview development dependencies")
	}
	return nil
}

// GOFLAGS uses Go's quoted-field syntax, not shell expansion. Preserve its
// tags when adding production, and its ldflags unless -ldflags overrides them.
func inheritedFlags(value string) (tags, ldflags string, err error) {
	for value = strings.TrimSpace(value); value != ""; value = strings.TrimSpace(value) {
		var field string
		if value[0] == '\'' || value[0] == '"' {
			end := strings.IndexByte(value[1:], value[0])
			if end < 0 {
				return "", "", errors.New("unclosed quote in GOFLAGS")
			}
			field, value = value[1:end+1], value[end+2:]
		} else {
			end := strings.IndexFunc(value, unicode.IsSpace)
			if end < 0 {
				end = len(value)
			}
			field, value = value[:end], value[end:]
		}
		if strings.HasPrefix(field, "-tags=") {
			tags = strings.TrimPrefix(field, "-tags=")
		}
		if strings.HasPrefix(field, "-ldflags=") {
			ldflags = strings.TrimPrefix(field, "-ldflags=")
		}
	}
	return tags, ldflags, nil
}

func compilerFlags(opts buildOptions, env goEnvironment, exists func(...string) bool) ([]string, string, error) {
	inheritedTags, ldflags, err := inheritedFlags(env.GOFLAGS)
	if err != nil {
		return nil, "", err
	}
	tags := []string{"desktop", "production"}
	for _, tag := range strings.FieldsFunc(inheritedTags+","+opts.tags, func(r rune) bool { return r == ',' || unicode.IsSpace(r) }) {
		if tag == "dev" || tag == "bindings" {
			return nil, "", fmt.Errorf("tag %q conflicts with a production desktop build", tag)
		}
		if !slices.Contains(tags, tag) {
			tags = append(tags, tag)
		}
	}
	if opts.ldflags != "" {
		ldflags = opts.ldflags
	}
	if env.GOOS == "windows" && !opts.console {
		ldflags = strings.TrimSpace(ldflags + " -H=windowsgui")
	}
	if env.GOOS != "linux" {
		if opts.webkit != "auto" || slices.Contains(tags, "webkit2_41") {
			return nil, "", errors.New("-webkit and webkit2_41 are Linux-only")
		}
		return tags, ldflags, nil
	}
	webkit := opts.webkit
	if slices.Contains(tags, "webkit2_41") {
		if webkit == "4.0" {
			return nil, "", errors.New("-webkit 4.0 conflicts with tag webkit2_41")
		}
		webkit = "4.1"
	}
	if webkit == "auto" {
		switch {
		case exists("gtk+-3.0", "webkit2gtk-4.1"):
			webkit = "4.1"
		case exists("gtk+-3.0", "webkit2gtk-4.0"):
			webkit = "4.0"
		default:
			return nil, "", errors.New("GTK3/WebKitGTK development packages not found by pkg-config; on Debian/Ubuntu install pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev")
		}
	} else if !exists("gtk+-3.0", "webkit2gtk-"+webkit) {
		return nil, "", fmt.Errorf("pkg-config cannot find gtk+-3.0 and webkit2gtk-%s; install the matching development packages or select another -webkit version", webkit)
	}
	if webkit == "4.1" && !slices.Contains(tags, "webkit2_41") {
		tags = append(tags, "webkit2_41")
	}
	return tags, ldflags, nil
}

func pkgConfigExists(packages ...string) bool {
	name := os.Getenv("PKG_CONFIG")
	if name == "" {
		name = "pkg-config"
	}
	return exec.Command(name, append([]string{"--exists"}, packages...)...).Run() == nil
}

// Append to CGO_LDFLAGS rather than replacing the user's library search paths
// and flags. The other Apple frameworks are already linked by Wails' cgo code.
func compilerEnvironment(environ []string, targetOS string) []string {
	result := slices.Clone(environ)
	if targetOS != "darwin" {
		return result
	}
	const key = "CGO_LDFLAGS="
	for i, field := range result {
		if strings.HasPrefix(field, key) {
			result[i] = field + " -framework UniformTypeIdentifiers"
			return result
		}
	}
	return append(result, key+"-framework UniformTypeIdentifiers")
}

type listedPackage struct {
	Dir, Name, ImportPath string
	Match                 []string
	Error                 *packageError
	DepsErrors            []packageError
}

type packageError struct{ Err string }

func findApp(ctx context.Context, pkg string, tags []string) (listedPackage, string, error) {
	// -e still reports Dir for a package with missing embedded assets. That
	// lets us explain how to fix them before go build prints an embed error.
	cmd := exec.CommandContext(ctx, "go", "list", "-e", "-deps", "-json=Dir,Name,ImportPath,Match,Error,DepsErrors", "-tags", strings.Join(tags, ","), pkg)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	data, err := cmd.Output()
	if err != nil {
		return listedPackage{}, "", fmt.Errorf("go list: %w\n%s", err, strings.TrimSpace(stderr.String()))
	}
	return readPackages(bytes.NewReader(data))
}

func readPackages(reader io.Reader) (listedPackage, string, error) {
	var app listedPackage
	var desktopDir string
	count := 0
	decoder := json.NewDecoder(reader)
	for {
		var pkg listedPackage
		err := decoder.Decode(&pkg)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return app, "", fmt.Errorf("decode go list: %w", err)
		}
		if len(pkg.Match) > 0 {
			app = pkg
			count++
		}
		if pkg.ImportPath == desktopImport {
			desktopDir = pkg.Dir
		}
	}
	if count != 1 {
		return app, "", fmt.Errorf("expected one main package, matched %d", count)
	}
	if app.Error != nil {
		return app, "", fmt.Errorf("load app: %s", app.Error.Err)
	}
	if app.Name != "main" {
		return app, "", fmt.Errorf("package %s is %q; expected main", app.ImportPath, app.Name)
	}
	if desktopDir == "" {
		if len(app.DepsErrors) > 0 {
			return app, "", fmt.Errorf("load app dependencies: %s", app.DepsErrors[0].Err)
		}
		return app, "", fmt.Errorf("package %s does not import %s", app.ImportPath, desktopImport)
	}
	return app, desktopDir, nil
}

func checkAssets(desktopDir string) error {
	dir := filepath.Join(desktopDir, "frontend", "dist")
	index, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil || len(bytes.TrimSpace(index)) == 0 || bytes.Contains(index, []byte("toolgui: assets not built")) {
		return fmt.Errorf("desktop frontend assets are missing or placeholders in %s; use a tagged toolgui-wails release, or run `task asset_wails` from the ToolGUI checkout root", dir)
	}
	var hasJS bool
	err = filepath.WalkDir(dir, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(name, ".js") {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			hasJS = hasJS || info.Size() > 0
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("read desktop frontend: %w", err)
	}
	if !hasJS {
		return fmt.Errorf("desktop frontend JavaScript is missing in %s; use a tagged toolgui-wails release, or run `task asset_wails` from the ToolGUI checkout root", dir)
	}
	return nil
}

func executableName(importPath, targetOS string) string {
	name := path.Base(importPath)
	if targetOS == "windows" {
		name += ".exe"
	}
	return name
}

func goJSON(ctx context.Context, value any, args ...string) error {
	cmd := exec.CommandContext(ctx, "go", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	data, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("go %s: %w\n%s", args[0], err, strings.TrimSpace(stderr.String()))
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("decode go %s: %w", args[0], err)
	}
	return nil
}
