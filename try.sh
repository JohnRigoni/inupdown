

sesh="inupdown"
cd ~/code/ho/inupdown

tmux new-session -s "$sesh" -d -x "$(tput cols)" -y "$(tput lines)"
tmux send-keys -t "$sesh".0 'nv .' C-m
tmux split-pane -t "$sesh".0 -v -l "14%"
tmux send-keys -t "$sesh".1 'go run main.go' C-m
tmux split-pane -t "$sesh".1 -h -l "50%"
tmux send-keys -t "$sesh".2 'templ generate -watch' C-m
tmux split-pane -t "$sesh".2 -v -l "50%"
tmux send-keys -t "$sesh".3 'npx tailwindcss -i ./tailwind.css -o ./assets/tstyle.css --watch' C-m

tmux a -t "$sesh"
# tmux send-keys -t "$sesh".0 'htop'
# split-window -v -p 32 \; 

# \
# send-keys 'htop' C-m \;  \
# split-window -v -p 16 \; \
# send-keys 'htop' C-m \;
