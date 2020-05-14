#!/usr/bin/env zsh

session="bima"
sessionexists=$(tmux ls | grep $session)

if [ "$sessionexists" = "" ]; then
  # Start New Session with our name
  tmux new-session -d -s $session

  tmux rename-window -t 0 'editor'
  tmux send-keys -t 'editor' 'vim .'  C-m

  tmux new-window -t ${session}:1 -n 'docker'
  tmux send-keys -t 'docker' 'cd server/bima; docker-compose up' C-m

  tmux select-window -t ${session}:1
  tmux split-window -h -t ${session}:1
  tmux send-keys -t ${session}:1 'cd server/bima; iex -S mix phx.server' C-m
else
  tmux attach -t ${session}:0
fi
