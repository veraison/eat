#!/usr/bin/env bash
set -ueo pipefail

if [[ "$( type -p cbor2diag.rb )" == "" ]]; then
	echo -e "\033[0;31merror\033[0m: cbor2diag.rb no found in PATH"
	exit 1
fi

this_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
dest_dir=$( realpath "$( dirname "$this_dir" )" )

for src_file in "$this_dir"/*.diag; do
	dest_file="$dest_dir"/$( basename "${src_file%.diag}.cbor" )

	echo "compiling ${dest_file}..."
	diag2cbor.rb < "$src_file" > "$dest_file"
done

echo "done."
