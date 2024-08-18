#!/usr/bin/env bash

set -xeuo pipefail

vault="${VAULT:-Sandbox}"
script_dir="$(dirname "$(realpath "$0")")"
item_tag=git-credential-op-test
op_command="op --vault '${vault}' --additional-tags ${item_tag}"
username=foo
password=bar
compose_file="$script_dir/docker-compose.yaml"

trap_exit() {
   docker compose -f "${compose_file}" down --rmi all --remove-orphans
}
trap trap_exit EXIT

op_cleanup() {
    op --vault "${vault}" item list --tags "${item_tag}" --format json \
    | jq -r '.[].id' \
    | tee /dev/tty \
    | xargs op --vault Sandbox item delete

    [[ -z "$(op_find_item)" ]]
}

git_init() {
    tmp_dir=$(mktemp -d)
    cd "${tmp_dir}" || exit 1

    git init . -b main
    git config user.email "me@example.com"
    git config user.name "me"
}

git_commit() {
    echo "$1" >> CONTENT
    git add .
    git commit -m "test"
}

git_push() {
    git push origin main --force
}

git_push_with_pass() {
    GIT_PASSWORD="$1" GIT_ASKPASS="$script_dir/askpass.bash" git push origin main --force
}

op_find_item() {
    op --vault "${vault}" item list --tags "${item_tag}" --format json | jq -r '.[].id'
}

test_username_in_url() {
    git_init
    git_commit "content"

    git remote add origin "http://${username}@localhost:8888/test"
    git config credential.helper "${op_command}"
    git_push_with_pass "bar"

    [[ -n "$(op_find_item)" ]]
}

test_username_and_password_in_url() {
    git_init
    git_commit "content"

    git remote add origin "http://${username}:${password}@localhost:8888/test"
    git config credential.helper "${op_command}"
    git_push

    [[ -n "$(op_find_item)" ]]
}

test_username_in_config() {
    git_init
    git_commit "content"

    git remote add origin "http://localhost:8888/test"
    git config credential.http://localhost:8888/test.helper "${op_command}"
    git config credential.http://localhost:8888/test.username "foo"
    git_push_with_pass "bar"

    [[ -n "$(op_find_item)" ]]
}

test_http_path() {
    git_init
    git_commit "content"

    git remote add origin "http://localhost:8888/test"
    git config credential.http://localhost:8888/test.helper "${op_command}"
    git config credential.http://localhost:8888/test.username "foo"
    git config credential.http://localhost:8888/test.useHttpPath true
    git_push_with_pass "bar"

    [[ -n "$(op_find_item)" ]]
}

test_erase_on_failure() {
    git_init
    git_commit "content"

    git remote add origin "http://localhost:8888/test"
    git config credential.helper "${op_command}"
    git config credential.username "foo"
    git_push_with_pass "bar"

    # Set the password on the item to another value.
    item_id="$(op_find_item)"
    op --vault "${vault}" item edit "${item_id}" 'password=baz'

    # Push again. Git will erase the item.
    git_commit "more content"
    git_push

    [[ -z "$(op_find_item)" ]]
}

# Run tests

docker compose -f "${compose_file}" up -d
until curl --output /dev/null --silent --head --fail "http://${username}:${password}@localhost:8888/test/info/refs"; do
    sleep 1
done

tests=(
    "test_username_in_url"
    "test_username_in_config"
    "test_username_and_password_in_url"
    "test_http_path" 
    "test_erase_on_failure"
)
for t in "${tests[@]}"; do
    op_cleanup
    (export GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=1; ${t}) || exit 1
    op_cleanup
done
