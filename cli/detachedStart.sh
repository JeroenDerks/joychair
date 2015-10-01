screen -S "joychair" -d -m
screen -r "joychair" -X stuff $'while true; do ./joychair default.toml; sleep 1; done\n'