#!/bin/bash

SESSION_NAME="back"
REPO_DIR="$HOME/gems_go_back"
GIT_REPO_URL="https://github.com/tomikartemik/gems_go_back"
STATUS_URL="https://api.dododrop.ru/admin/change-status"

echo "Starting deployment script"

# Функция для выполнения команд в tmux
tmux_send() {
    tmux send-keys -t $SESSION_NAME "$1" C-m
}

# Функция для отправки запросов для изменения статуса
change_status() {
    local status=$1
    curl -X POST "$STATUS_URL" \
        -H "Content-Type: application/json" \
        -d "{\"status\": $status}"
}

change_status false

#echo "Waiting for 180 seconds..."
#sleep 180

tmux kill-session -t $SESSION_NAME

echo "Cloning/updating repository"
if [ ! -d "$REPO_DIR" ]; then
    git clone $GIT_REPO_URL $REPO_DIR
fi

cd $REPO_DIR && git pull origin main

tmux new-session -d -s $SESSION_NAME
tmux_send "cd cmd"
tmux_send "go run main.go"

sleep 10

change_status true