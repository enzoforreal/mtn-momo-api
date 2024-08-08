#!/bin/bash

# Function to display the ASCII Art
function show_ascii_art {
    cat << "EOF"
 __  __          __  __          
|  \/  |        |  \/  |         
| \  / |  ___   | \  / |  ___    
| |\/| | / _ \  | |\/| | / _ \   
| |  | || (_) | | |  | || (_) |  
|_|  |_| \___/  |_|  |_| \___/   
EOF
}

function show_help {
    show_ascii_art
    echo
    echo "momo is a tool for managing the momo API."
    echo
    echo "Usage:"
    echo "    momo <command> [arguments]"
    echo
    echo "The commands are:"
    echo "    start       start the server"
    echo "    test        run integration tests with optional UUID"
    echo "    update      update dependencies"
    echo
    echo "Use 'momo help <command>' for more information about a command."
}

function start_server {
    echo "Starting the server..."
    cd "$(dirname "$0")/.."
    go run example/main.go
}

function run_tests {
    echo "Running integration tests..."
    script_dir="$(dirname "$0")/../tests/integration"
    cd "$script_dir" || exit 1

    auth_req_id="$1"
    if [ -z "$auth_req_id" ]; then
        echo "No auth_req_id provided, will be generated inside the script."
    else
        echo "auth_req_id passed to tests: $auth_req_id"
    fi

    # Array of scripts to run in order
    scripts=(
        "create-api-user.sh"
        "create-api-key.sh"
        "create-oauth2-token.sh"
        "get-account-balance.sh"
        "get-auth-token.sh"
        "request-to-pay.sh"  
        "payment-status.sh"
    )

    # Loop over scripts and execute them
    for script in "${scripts[@]}"; do
        echo "Running $script..."
        bash "$script"

        # Check for errors after running each script
        if [ $? -ne 0 ]; then
            echo "Error: $script failed to execute properly."
            exit 1
        fi
    done

    # Verify presence of necessary files
    if [ ! -f /tmp/X-Reference-Id-requesttopay ]; then
      echo "Error: /tmp/X-Reference-Id-requesttopay is missing."
      exit 1
    fi
}

function update_deps {
    echo "Updating Go dependencies..."
    cd "$(dirname "$0")/.."
    go get -u ./...
    go mod tidy
}

if [[ $# -eq 0 ]]; then
    show_help
    exit 0
fi

case "$1" in
    start)
        start_server
        ;;
    test)
        run_tests "$2"
        ;;
    update)
        update_deps
        ;;
    help)
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        show_help
        ;;
esac
