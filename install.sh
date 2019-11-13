#!/usr/bin/env bash
#
# A good old bash | curl script for direnv.
#
set -euo pipefail

{ # Prevent execution if this script was only partially downloaded

  log() {
    echo "[installer] $*" >&2
  }

  die() {
    log "$@"
    exit 1
  }

  at_exit() {
    ret=$?
    if [[ $ret -gt 0 ]]; then
      log "the script failed with error $ret.\n" \
        "\n" \
        "To report installation errors, submit an issue to\n" \
        "    https://github.com/direnv/direnv/issues/new/choose"
    fi
    exit "$ret"
  }
  trap at_exit EXIT

  kernel=$(uname -s | tr "[:upper:]" "[:lower:]")
  case "$(uname -m)" in
    x86_64)
      machine=amd64
      ;;
    i686 | i386)
      machine=386
      ;;
    *)
      die "Machine $(uname -m) not supported by the installer.\n" \
        "Go to https://direnv for alternate installation methods."
      ;;
  esac
  log "kernel=$kernel machine=$machine"

  : "${use_sudo:=}"
  : "${bin_path:=}"

  if [[ -z "$bin_path" ]]; then
    log "looking for a writeable PATH"
    for path in $(echo "$PATH" | tr ':' '\n'); do
      if [[ -w $path ]]; then
        bin_path=$path
        break
      fi
    done
  fi
  if [[ -z "$bin_path" ]]; then
    die "did not find a writeable path in $PATH"
  fi
  echo "bin_path=$bin_path"

  log "looking for a download URL"
  download_url=$(
    curl -fL https://api.github.com/repos/direnv/direnv/releases/latest \
    | grep browser_download_url \
    | cut -d '"' -f 4 \
    | grep "direnv.$kernel.$machine"
  )
  echo "download_url=$download_url"

  log "downloading"
  curl -o "$bin_path/direnv" -fL "$download_url"
  chmod +x "$bin_path/direnv"

  cat <<DONE

The direnv binary is now avaible in:

    $bin_path/direnv

The last step is to configure your shell to use it. For example for bash, add
the following lines at the end of your ~/.bashrc:

    eval "\$(direnv hook bash)"

Then restart the shell.

For other shells, see https://direnv.net/docs/hook.html

Thanks!
DONE
}
