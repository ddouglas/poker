#!/bin/bash

# Set the root directory
root_directory="./.builds"

# Check if the root directory exists
if [ ! -d "$root_directory" ]; then
  echo "Root directory does not exist: $root_directory"
  exit 1
fi

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o .builds/ -tags lambda.norpc ./cmd/...

# Loop through files in the root directory
for file in "$root_directory"/*; do
  # Check if the item is a file
  if [ -f "$file" ]; then
    # Extract the file name without the path
    filename=$(basename "$file")

    # Create a subdirectory with the same name as the file (without extension)
    directory_name="${filename%.*}"

    # Check if the directory already exists, and if not, create it
    if [ ! -d "$root_directory/$directory_name-handler" ]; then
      mkdir "$root_directory/$directory_name-handler"
    fi

    # Move the file into the subdirectory
    mv "$file" "$root_directory/$directory_name-handler/bootstrap"

    cd "$root_directory/$directory_name-handler"
    ls -la .

    zip "$directory_name.zip" ./bootstrap

    mv "$directory_name.zip" ../

    cd ../

    echo "$root_directory/$directory_name-handler/"

    rm -r "./$directory_name-handler/"

    aws-vault exec --no-session ots -- aws lambda update-function-code --function-name $directory_name-handler --zip-file "fileb://$directory_name.zip"

    rm "$directory_name.zip"
  fi
done

echo "Files moved to subdirectories successfully!"
