package cmd

// Hooks contains all hooks in the form: map[shellName -> map[hookName -> hook]]
type Hooks map[string]map[string]string

// Get gets a hook
func (hooks Hooks) Get(hookName string, shellName string) string {
	hooksForShell := hooks[shellName]
	if hooksForShell != nil {
		return hooksForShell[hookName]
	}

	return ""
}

// Set sets a hook
func (hooks Hooks) Set(hookName string, hook string, shellName string) {
	hooksForShell := hooks[shellName]
	if hooksForShell == nil {
		hooksForShell = map[string]string{}
		hooks[shellName] = hooksForShell
	}

	hooksForShell[hookName] = hook
}
