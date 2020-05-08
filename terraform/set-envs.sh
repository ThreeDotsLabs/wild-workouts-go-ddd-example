#!/bin/bash
readonly env_prefix=TF_VAR_
readonly example_file=.env.example
readonly target_file=.env

answer=""

function ask() {
    local var="$1"
    local example="$2"
    local current="$3"

    local name=${var/$env_prefix/}

    if [ -n "$current" ]; then
        hint="current: $current"
    else
        hint="example: $example"
    fi

    answer=""

    while [ -z "$answer" ]; do
        echo -ne "\t$name [$hint]: " 
        read -r answer </dev/tty

        if [ -z "$answer" ] && [ -n "$current" ]; then
            answer="$current"
            break
        fi
    done
}

if [ -f "$target_file" ]; then
    source "$target_file"
fi

if [ -z "$TF_VAR_user" ]; then
    readonly current_user=$(gcloud auth list --filter="status:ACTIVE" --format="value(account)")

    echo "Active gcloud account: $current_user"
    echo
    echo "Hit Enter if you want to continue with this account."
    echo "Otherwise, cancel with Ctrl+C and login with:"
    echo -e "\tgcloud auth login"
    echo
    echo "Or set active account with:"
    echo -e "\tgcloud config set account 'ACCOUNT'"
    echo
    echo "Remember to add application-default credentials with:"
    echo -e "\tgcloud auth application-default login"

    read

    TF_VAR_user="$current_user"
fi

if [ -z "$TF_VAR_billing_account" ]; then
    readonly current_billing_account=$(gcloud beta billing accounts list --filter="open=true" --format="value(displayName)" --limit 1)

    echo "Found billing account: $current_billing_account"
    TF_VAR_billing_account="$current_billing_account"
fi

echo "Fill all required parameters:"

new_content=""

while IFS= read -r line; do
    var=$(echo "$line" | cut -d= -f1)
    example=$(echo "$line" | cut -d= -f2)
    current="${!var}"

    ask "$var" "$example" "$current"

    new_content+="export $var=\"$answer\"\n"
done < <(cut -d' ' -f2 "$example_file")

echo
echo "Saving to $target_file:"
echo -en "$new_content" | tee "$target_file"
