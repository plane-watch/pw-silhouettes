#!/usr/bin/env bash
set -euo pipefail
shopt -s nullglob

DELAY=5   # 50ms per frame
LOOP=0    # 0 = infinite

# Find unique bases that have frame suffixes like "-1", "-2", etc
bases=$(
  ls *-[0-9]*.png 2>/dev/null \
    | sed -E 's/-[0-9]+\.png$//' \
    | sort -u
)

for base in $bases; do
  frames=( "${base}"-[0-9]*.png )
  [[ ${#frames[@]} -gt 0 ]] || continue

  # Sort frames numerically by the number after the last dash, e.g. SF34-12.png -> 12
  frames_sorted=$(
    printf '%s\n' "${frames[@]}" \
      | awk -F'-' '{
          n=$NF
          sub(/\.png$/, "", n)
          print n "\t" $0
        }' \
      | sort -n -k1,1 \
      | cut -f2-
  )

  out="${base}.gif"
  echo "Creating $out ..."

  # shellcheck disable=SC2086
  convert -delay "$DELAY" -loop "$LOOP" $frames_sorted -layers Optimize "$out"
done

