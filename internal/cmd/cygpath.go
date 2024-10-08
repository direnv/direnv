package cmd

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

// pathMutex is used to synchronise writes to the PATH environment variable.
var pathMutex sync.Mutex

// winExec is a windows approximation of the syscall.Exec function.
//
// In Unix environments, syscall.Exec replaces the current process with
// commandPath. Since we can't do that in Windows, we approximate that
// behaviour.
func winExec(commandPath string, args []string, env []string) error {
	// https://man7.org/linux/man-pages/man2/execve.2.html
	// In the execve and syscall.Exec interface, commandPath is also args[0].
	// Let's strip that from args for this function.
	if len(args) > 0 {
		args = args[1:]
	}

	logDebug("winExec: %s %v", commandPath, args)

	cmd := exec.Command(commandPath, args...)
	cmd.Env = env

	// Set the standard input, output, and error to the current process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logDebug("winExec: %s: could not start command: %v", commandPath, err)
		return err
	}

	logDebug("winExec: %s: waiting for result", commandPath)
	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				logDebug("winExec: %s: exit error: %v", commandPath, err)
				os.Exit(status.ExitStatus())
			}
		}
		logDebug("winExec: %s: non-exit error: %v", commandPath, err)
		os.Exit(1)
	}

	logDebug("winExec: success: %s %v", commandPath, args)
	os.Exit(0)

	return nil
}

// isWinPathList returns true if pathList looks like a Windows path list.
func isWinPathList(pathList string) bool {
	return strings.Contains(pathList, ";")
}

// isWinPath returns true if path looks like a Windows path.
//
// NOTE: This will return true for a Unix path list like /bla:/bla
func isWinPath(path string) bool {
	return strings.Contains(path, ":\\") || strings.Contains(path, ":/")
}

// errCygpath is returned when cygpath is not available.
var errCygpath = errors.New("cygpath is not available")

// cygpath represents the data and logic needed for cygpath environments
// (MSYS2/Cygwin).
//
// Also see the stdlib.sh for more cygpath usage.
//
// This functionality was originally added in the following PR:
// https://github.com/direnv/direnv/pull/1291
type cygpath struct {
	// context is the environment in which cygpath is running (cygwin, msys2).
	context string
}

// newCygpath will return nil if direnv is not launched from a cygpath environment.
//
// This function is executed at every direnv invocation, so be mindful of its
// performance impact, especially to the typical non-cygpath environment. That
// means, exit early and perform the most costly operations at the end.
func newCygpath() *cygpath {

	if runtime.GOOS != "windows" {
		// This is not Windows, so we won't have cygpath.
		return nil
	}

	// Check the format of the PATH environment variable. If it is in Unix
	// format, then direnv is not running in a working context.
	if p := os.Getenv("PATH"); !isWinPathList(p) {
		// In some cases, the incoming PATH environment variable is in Unix format,
		// which breaks the exec.LookPath function, which means we can't find
		// cygpath. If we can't find cygpath, then we can't convert the incoming
		// PATH to Windows format.
		// TODO: Can we correct this situation without having cygpath? Maybe the
		// Windows location of cygpath can be derived using heuristics.
		logDebugInit("newCygpath: expected Windows PATH format but got Unix PATH format")
		return nil
	}

	if _, err := exec.LookPath("cygpath"); err != nil {
		// We require cygpath for more stable conversion of select environment
		// variables. The stdlib.sh also contains calls to cygpath.
		// TODO: A future PR can try to implement cygpath logic directly in Go,
		// and then have the stdlib.sh call into direnv to resolve paths, e.g.
		// direnv unix_path "$@" and direnv unix_path_list "$@".
		logDebugInit("newCygpath: %v", errCygpath)
		return nil
	}

	context := "cygwin"
	if msystem := os.Getenv("MSYSTEM"); len(msystem) > 0 {
		context = "msys2"
	}

	logDebugInit("newCygpath: context=%s", context)

	return &cygpath{
		context: context,
	}
}

// Cygpath is the cygpath environment instance. It can be nil if direnv is not
// running in a cygpath environment.
var Cygpath = newCygpath()

// execCygpath calls cygpath and returns the output.
//
// TODO: Try to reduce all calls to cygpath by first using heuristic methods to
// check if cygpath is required for the conversion.
func (m cygpath) execCygpath(args ...string) (string, error) {

	cmd := exec.Command("cygpath", args...)

	if len(os.Args) > 0 {
		logDebugInit("execCygpath: %v: cygpath %.100s...", os.Args[1:], strings.Join(cmd.Args[1:], " "))
	} else {
		logDebugInit("execCygpath: cygpath %.100s...", strings.Join(cmd.Args[1:], " "))
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// tryUnixPathList converts the Windows pathList to a Unix path list using
// cygpath. If the cygpath conversion fails, it returns the original.
func (m cygpath) tryUnixPathList(pathList string) string {
	if !isWinPathList(pathList) {
		// This does not look like a Windows path list, so let's skip the
		// cygpath call.
		logDebugInit("tryUnixPathList: already a Unix path, keeping original")
		return pathList
	}
	p, err := m.execCygpath("-u", "-p", pathList)
	if err == nil {
		return p
	}
	logDebugInit("tryUnixPathList: could not convert path list using cygpath, keeping original: %s", pathList)
	return pathList
}

// initExportEnv converts paths into the format they are required.
//
// If initExportEnv is called during application startup, then the final state of the
// exported variables will be determined by this function.
func (m cygpath) initExportEnv(env Env) Env {

	// Conversion of environment variables is not necessary for every direnv
	// command, e.g. direnv version does not care about environment variables.
	// All commands that use actionWithConfig probably need their paths
	// corrected. Commands that use actionSimple with _ Env don't care about the
	// environment (or about the config) and probably don't need it fixed.
	// Commands that use env variables that are not path-like (e.g.
	// DIRENV_WATCHES), also don't need to be fixed. Commands that are marked as
	// private (Private: True) are probably used by the stdlib.sh, which means
	// they will expect Unix paths.

	if len(os.Args) > 0 {
		cmd := os.Args[1]
		switch cmd {
		case
			"watch",
			"stdlib",
			"version":
			return env
		}
	}

	// PATH
	// If PATH is not converted, then the "space dir" bash test fails with a
	// "command not found" error. We convert the PATH that will be exported,
	// into Unix form, but the process path via os.Getenv("PATH") will remain in
	// the form received (probably Windows form).
	if p := env["PATH"]; len(p) > 0 {
		env["PATH"] = m.tryUnixPathList(p)
	}

	// Some environments also have a "Path" (as opposed to a "PATH")
	// environment.
	if p := env["Path"]; len(p) > 0 {
		env["Path"] = m.tryUnixPathList(p)
	}

	return env
}

// LookPath functions like exec.LookPath but allows specifying the PATH as
// pathenv.
//
// This Windows version of the function opts in to the built-in Go
// functionality.
//
// In some cases file and/or pathenv gets passed in using Unix path format. We
// convert these paths to Windows path format as expected by the
// exec.LookPath function.
//
// Also see:
// - https://cs.opensource.google/go/go/+/refs/tags/go1.22.4:src/os/exec/lp_windows.go
// - https://cs.opensource.google/go/go/+/refs/tags/go1.22.4:src/os/exec/lp_windows_test.go
func (m cygpath) LookPath(file string, pathenv string) (string, error) {

	logDebug("LookPath: %s: lookup requested", file)

	// pathenv must be in the Windows format.

	if !isWinPathList(pathenv) {
		if p, err := m.execCygpath("-w", "-p", pathenv); err == nil {
			pathenv = p
		}
	}

	// file must be in the Windows format.

	if !isWinPath(file) {
		if p, err := m.execCygpath("-w", file); err == nil {
			file = p
		}
	}

	logDebug("LookPath: %s: lookup ready", file)

	pathMutex.Lock()
	defer pathMutex.Unlock()

	// Save the current PATH to restore it later
	originalPath, found := os.LookupEnv("PATH")
	if !found {
		return "", fs.ErrNotExist
	}

	// Set the PATH environment variable to the provided pathenv, and restore it
	// when this function exits.
	os.Setenv("PATH", pathenv)
	defer os.Setenv("PATH", originalPath)

	// Now, call exec.LookPath with the temporarily set PATH.
	logDebug("LookPath: %s: lookup initiated", file)
	return exec.LookPath(file)
}

// Exec is a windows approximation of the syscall.Exec function.
func (m cygpath) Exec(commandPath string, args []string, env Env) error {
	if pathenv := env["PATH"]; !isWinPathList(pathenv) {
		if p, err := m.execCygpath("-w", "-p", pathenv); err == nil {
			env["PATH"] = p
		}
	}
	return winExec(commandPath, args, env.ToGoEnv())
}

// tryWindowsPath ensures that the incoming path can be used by Go.
//
// In an MSYS2 cygpath environment, we typically get paths in their Windows
// form. We use this as an optimization hint, so that we don't waste time
// calling out to cygpath.
//
// In a Cygwin cygpath environment, we typically get paths in their Unix form.
func (m cygpath) tryWindowsPath(path string) string {
	if len(path) == 0 {
		return path
	}
	if isWinPath(path) {
		return path
	}
	if p, err := m.execCygpath("-w", path); err == nil {
		logDebug("tryWindowsPath: %s -> %s", path, p)
		return p
	}
	logDebug("tryWindowsPath: %s: could not fix the path", path)
	return path
}

// initTempEnv converts a select set of paths so that they will work with the
// Go runtime. A function is returned that can revert the changes.
//
// NOTE: These are temporary, since we don't want to export them to the shell
// with our modifications. It is only used temporarily by the Go runtime.
func (m cygpath) initTempEnv(env Env) func() {

	// XDG_CACHE_HOME is used by LoadConfig (xdg package) and required to be in
	// Windows format.
	xdgCacheHome := env["XDG_CACHE_HOME"]
	if len(xdgCacheHome) > 0 {
		env["XDG_CACHE_HOME"] = m.tryWindowsPath(xdgCacheHome)
	}

	// XDG_CONFIG_HOME is set by direnv-test-common.sh and used by LoadConfig
	// (xdg package), and required to be in Windows format.
	xdgConfigHome := env["XDG_CONFIG_HOME"]
	if len(xdgConfigHome) > 0 {
		env["XDG_CONFIG_HOME"] = m.tryWindowsPath(xdgConfigHome)
	}

	// XDG_DATA_HOME is set by direnv-test-common.sh and used by LoadConfig
	// (xdg package), and required to be in Windows format.
	xdgDataHome := env["XDG_DATA_HOME"]
	if len(xdgDataHome) > 0 {
		env["XDG_DATA_HOME"] = m.tryWindowsPath(xdgDataHome)
	}

	// DIRENV_CONFIG is set by direnv-test-common.sh and used by LoadConfig and
	// required to be in Windows format.
	direnvConfig := env["DIRENV_CONFIG"]
	if len(direnvConfig) > 0 {
		env["DIRENV_CONFIG"] = m.tryWindowsPath(direnvConfig)
	}

	// DIRENV_BASH is set by direnv-test-common.sh and used by LoadConfig and
	// required to be in Windows format.
	direnvBash := env["DIRENV_BASH"]
	if len(direnvBash) > 0 {
		env["DIRENV_BASH"] = m.tryWindowsPath(direnvBash)
	}

	// DIRENV_FILE is used by LoadConfig and required to be in Windows format.
	// Also see RCFromPath and RCFromEnv.
	direnvFile := env["DIRENV_FILE"]
	if len(direnvFile) > 0 {
		env["DIRENV_FILE"] = m.tryWindowsPath(direnvFile)
	}

	// HOME is used by LoadConfig and the xdg package and required to be in
	// Windows format.
	home := env["HOME"]
	if len(home) > 0 {
		env["HOME"] = m.tryWindowsPath(home)
	}

	// DIRENV_DIR is not really used anywhere, so let's do nothing with it.

	// DIRENV_DUMP_FILE_PATH is not used by the runtime using the env object,
	// but it is used using the os.Getenv function (see cmd_dump.go).

	return func() {
		if len(xdgCacheHome) > 0 {
			env["XDG_CACHE_HOME"] = xdgCacheHome
		}
		if len(xdgConfigHome) > 0 {
			env["XDG_CONFIG_HOME"] = xdgConfigHome
		}
		if len(xdgDataHome) > 0 {
			env["XDG_DATA_HOME"] = xdgDataHome
		}
		if len(direnvConfig) > 0 {
			env["DIRENV_CONFIG"] = direnvConfig
		}
		if len(direnvBash) > 0 {
			env["DIRENV_BASH"] = direnvBash
		}
		if len(direnvFile) > 0 {
			env["DIRENV_FILE"] = direnvFile
		}
		if len(home) > 0 {
			env["HOME"] = home
		}
	}
}

// InitExportEnv provides an opportunity to fix environment variables in certain
// environments (e.g. MSYS2/Cygwin on Windows).
func InitExportEnv(env Env) Env {
	if Cygpath != nil {
		return Cygpath.initExportEnv(env)
	}
	return env
}

// FixPath tries to fix the input path according to the requirements of the
// runtime environment.
func FixPath(path string) string {
	if Cygpath != nil {
		return Cygpath.tryWindowsPath(path)
	}
	return path
}

// InitTempEnv converts a select set of paths so that they will work with the
// Go runtime. A function is returned that can undo the changes.
func InitTempEnv(env Env) func() {
	if Cygpath != nil {
		return Cygpath.initTempEnv(env)
	}
	return func() {}
}
