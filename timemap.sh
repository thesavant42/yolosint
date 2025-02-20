#!/bin/bash
# timemap.sh - get a list of URL paths known to Wayback Machine
# timemap.sh <hostname>
# json and url list for hostname will be timestamped in <pwd>
#
# Check if hostname argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <hostname>"
  exit 1
fi

# Assign the argument to a variable
HOSTNAME="$1"
TIMESTAMP=$(date +"%Y%m%d%H%M%S")

# Define filenames
JSON_FILE="${HOSTNAME}_${TIMESTAMP}.json"
URL_FILE="${HOSTNAME}_${TIMESTAMP}_urls.txt"

# Construct Wayback Machine API URL
WAYBACK_URL="https://web.archive.org/web/timemap/json?url=${HOSTNAME}&matchType=prefix&collapse=urlkey&output=json&fl=original%2Cmimetype%2Ctimestamp%2Cendtimestamp%2Cgroupcount%2Cuniqcount&filter=!statuscode%3A%5B45%5D..&limit=10000&_=$(date +%s)"

# Fetch data from Wayback Machine and save to JSON file
echo "Fetching data from: $WAYBACK_URL"
curl -s "$WAYBACK_URL" -o "$JSON_FILE"

# Check if the file was successfully downloaded
if [ ! -s "$JSON_FILE" ]; then
  echo "Error: Failed to retrieve data or empty result."
  rm -f "$JSON_FILE"
  exit 1
fi

# Extract URLs using jq and save to the URL file
jq -r ".[1:] | map(.[0]) | .[]" "$JSON_FILE" > "$URL_FILE"

# Output result
echo "JSON data saved to: $JSON_FILE"
echo "Extracted URLs saved to: $URL_FILE"
echo "First 10 URLs:"
head -n 10 "$URL_FILE"
