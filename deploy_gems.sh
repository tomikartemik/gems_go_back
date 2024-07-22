#!/bin/bash

SESSION_NAME="back"
REPO_DIR="$HOME/gems_github/gems_go_back"
GIT_REPO_URL="https://github.com/tomikartemik/gems_go_back"
STATUS_URL="https://api.youngrusssia.ru/admin/change-status"

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

echo "Changing game status to false..."
change_status false

echo "Waiting for 30 seconds..."
sleep 30

echo "Killing existing tmux session (if any)"
tmux kill-session -t $SESSION_NAME

echo "Cloning/updating repository"
if [ ! -d "$REPO_DIR" ]; then
    git clone $GIT_REPO_URL $REPO_DIR
fi

cd $REPO_DIR && git pull origin main

echo "Starting new tmux session"
tmux new-session -d -s $SESSION_NAME
tmux_send "cd $REPO_DIR/cmd"
tmux_send "go run main.go"

echo "Waiting for 10 seconds..."
sleep 10

echo "Changing game status to true..."
change_status true

echo "Deployment script finished"