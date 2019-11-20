## cnotifyd - LXD Fanotify

#### Get CT PID:
`lxc info ct01 | grep Pid`

#### Add watcher:
`echo '{"Action": "Add", "PID": 27547, "Name": "ct01"}' | socat - UNIX-CONNECT:/var/run/cnotifyd.socket`

#### Remove watcher:
`echo '{"Action": "Remove", "PID": 27547, "Name": "ct01"}' | socat - UNIX-CONNECT:/var/run/cnotifyd.socket`

#### Configure LXD hook with TRACE log-level:
`echo -e "lxc.hook.start-host = /opt/lxd/cnotifyd -hook\nlxc.hook.post-stop = /opt/lxd/cnotifyd -hook\nlxc.log.level = TRACE" | lxc config set ct01 raw.lxc -`

#### Configure LXD hook with default log-level:
`echo -e "lxc.hook.start-host = /opt/lxd/cnotifyd -hook\nlxc.hook.post-stop = /opt/lxd/cnotifyd -hook | lxc config set ct01 raw.lxc -`

#### Add hook config to `/usr/share/lxc/config/common.conf.d/00-fanotify.conf`:
`cat <<EOF > /usr/share/lxc/config/common.conf.d/00-fanotify.conf
lxc.hook.start-host = /opt/lxd/cnotifyd -hook
lxc.hook.post-stop = /opt/lxd/cnotifyd -hook
EOF`