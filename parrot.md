192.168.4.75


### telnet 192.168.42.1


* /tmp/arducopter --serial0 udp:192.168.4.75:14550 --serial3 /dev/ttyPA1



* cp /data/ftp/internal_000/arducopter /usr/bin/
* chmod +x /usr/bin/arducopter
* /usr/bin/arducopter --help


# this works
* /usr/bin/arducopter --serial0 udp:192.168.4.75:14550 --serial3 /dev/ttyPA1 --defaults /data/ftp/internal_000/ardupilot/* bebop.parm


# Run this on the drone's Telnet session before starting arducopter
* route add -net 192.168.4.0 netmask 255.255.255.0 gw 192.168.42.1

# finally got it working with this----
* /usr/bin/arducopter --serial0 udp:192.168.42.2:14550 --serial3 /dev/ttyPA1 --defaults /data/ftp/internal_000/ardupilot/* * * bebop.parm
