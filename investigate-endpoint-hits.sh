if [ -z "$1" ]; then
  echo "Usage: $0 <ffuf-results.json>" # Need to include the output filename
  exit 1
fi

FILE=$1

if [ ! -f "$FILE" ]; then
  echo "Error: File '$FILE' not found."
  exit 1
fi

jq -r '.results[].url' "$FILE" | while read url; do
  echo "=== GET $url ==="
  curl -s "$url"
  echo -e "\n\n"
done
