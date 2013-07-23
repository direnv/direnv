function _direnv
  if test $argv[1] = 0 -a -n "$__fish_direnv"
  	set -l direnv (eval $__fish_direnv export fish)
    eval $direnv
  end
end
#
# The following functions add support for a directory history
#

function cd --description "Change directory"

	# Skip history in subshells
	if status --is-command-substitution
		builtin cd $argv
  	set -l cd_status $status
    _direnv $status
		return $cd_status
	end

	# Avoid set completions
	set -l previous $PWD

	if test $argv[1] = - ^/dev/null
		if test "$__fish_cd_direction" = next ^/dev/null
			nextd
		else
			prevd
		end
    	set -l cd_status $status
      _direnv $cd_status
		return $cd_status
	end

	builtin cd $argv[1]
	set -l cd_status $status

	if test $cd_status = 0 -a "$PWD" != "$previous"
		set -g dirprev $dirprev $previous
		set -e dirnext
		set -g __fish_cd_direction prev
	end

  _direnv $cd_status

	return $cd_status
end

