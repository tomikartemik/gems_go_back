#!/bin/bash

SESSION_NAME="back"
REPO_DIR="$HOME/artem/gems_go_back"

tmux_send() {
    tmux send-keys -t $SESSION_NAME "$1" C-m
}

tmux kill-session -t $SESSION_NAME

cd $REPO_DIR && git pull origin main

# Запуск нового экземпляра в новой tmux сессии
tmux new-session -d -s $SESSION_NAME
tmux_send "cd $REPO_DIR/cmd"
tmux_send "go run main.go"