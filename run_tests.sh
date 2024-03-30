#!/bin/bash

# Define the commands with placeholders for ports
commands=(
    "go test -race -buildvcs -uri='mongodb://root:toor@localhost:27010' ./..."
    "go test -race -buildvcs -uri='mongodb://root:toor@localhost:27011' ./..."
    "go test -race -buildvcs -uri='mongodb://root:toor@localhost:27012' ./..."
)

# Function to run a command
run_command() {
    local command="$1"
    local port="$2"
    # Replace the placeholder 'port' with the actual port number
    command="${command/port/$port}"
    # Run the command and capture its output
    output=$(eval "$command" 2>&1)
    ret=$?
    if [ $ret -ne 0 ]; then
        echo "Command '$command' failed with exit code $ret and output:"
        echo "$output"
        # If one command fails, set global failure flag
        failed=true
    else
        echo "Output of command '$command':"
        echo "$output"
    fi
}

# Function to flush output
flush_output() {
    # Flush output buffers
    echo
}

# Initialize failure flag
failed=false

# Loop through the commands and run them in parallel
for ((i = 0; i < ${#commands[@]}; i++)); do
    # Run each command in the background
    run_command "${commands[i]}" "$((i + 1))" &
done

# Wait for all background processes to finish
wait

# Flush the output
flush_output

# Check if any command failed
if [ "$failed" = true ]; then
    exit -1
else
    exit 0
fi