package cmd

import (
	"regexp"

	"github.com/direnv/direnv/v2/gzenv"
)

// State contains direnv's state for an environment.
// It gets stored in the environment variable DIRENV_STATE
type State struct {
	Hooks Hooks
}

// Marshal marshals state to the gzenv format
func (state *State) Marshal() string {
	return gzenv.Marshal(state)
}

// UnmarshalState unmarshals state from the gzenv format
func UnmarshalState(from string) (*State, error) {
	if from == "" {
		return newState(), nil
	}

	state := new(State)
	return state, gzenv.Unmarshal(from, state)
}

// MakeState creates State based on the environment variables set in env
func MakeState(env Env) *State {
	state := newState()
	hooks := state.Hooks

	regex := regexp.MustCompile(DIRENV_HOOK_PREFIX + `(.+)_([^_]+)`)
	for envVarName, envVarValue := range env {
		matches := regex.FindStringSubmatch(envVarName)
		if matches == nil {
			continue
		}
		hookName := matches[1]
		shellName := matches[2]

		hooks.Set(hookName, envVarValue, shellName)
	}

	return state
}

func newState() *State {
	return &State{
		Hooks: make(Hooks),
	}
}
