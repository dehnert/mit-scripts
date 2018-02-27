#!/bin/sh
## Joe Presbrey <presbrey@mit.edu>
## Quentin Smith <quentin@mit.edu>
## Mitchell Berger <mitchb@mit.edu>
## SIPB Scripts LVS Firewall marks

iptables -F -t mangle

# Create a table for regular scripts hosts
iptables -t mangle -N scripts 2>/dev/null || :

# scripts-vhosts.mit.edu
iptables -A PREROUTING -t mangle -d 18.181.0.46 -j scripts
iptables -A PREROUTING -t mangle -d 18.4.86.46 -j scripts
# scripts.mit.edu
iptables -A PREROUTING -t mangle -d 18.181.0.43 -j scripts
iptables -A PREROUTING -t mangle -d 18.4.86.43 -j scripts
# scripts-cert.mit.edu
iptables -A PREROUTING -t mangle -d 18.181.0.50 -j scripts
iptables -A PREROUTING -t mangle -d 18.4.86.50 -j scripts

# Send Apache-bound traffic to FWM 2 (load-balanced)
iptables -A scripts -t mangle -m tcp -m multiport -p tcp --dports 80,443,444 -j MARK --set-mark 2
# Send SMTP-bound traffic to FWM 3 (load-balanced)
iptables -A scripts -t mangle -m tcp -p tcp --dport 25 -j MARK --set-mark 3
# Send finger-bound traffic to FWM 255 (the LVS director itself)
iptables -A scripts -t mangle -m tcp -p tcp --dport 78:79 -j MARK --set-mark 255
# Send everything else to FWM 1 (primary)
iptables -A scripts -t mangle -m mark --mark 0 -j MARK --set-mark 1

# webzephyr.mit.edu is special because its SMTP needs to always go to the primary (FWM 1)
iptables -A PREROUTING -t mangle -m tcp -m multiport -p tcp -d 18.181.0.49 --dports 80,443,444 -j MARK --set-mark 2
iptables -A PREROUTING -t mangle -m mark --mark 0 -d 18.181.0.49 -j MARK --set-mark 1

# scripts-primary.mit.edu goes to the primary (FWM 1) on all ports
iptables -A PREROUTING -t mangle -d 18.181.0.182 -j MARK --set-mark 1
iptables -A PREROUTING -t mangle -d 18.4.86.182 -j MARK --set-mark 1

# sipb.mit.edu acts like regular scripts for the web ports, everything else goes to i-hate-penguins.xvm.mit.edu (FWM 4)
iptables -A PREROUTING -t mangle -m tcp -m multiport -p tcp -d 18.181.0.29 --dports 80,443,444 -j MARK --set-mark 2
iptables -A PREROUTING -t mangle -m tcp -m multiport -p tcp -d 18.4.86.29 --dports 80,443,444 -j MARK --set-mark 2
# Also send port 25 there too because the IP is shared with rtfm.mit.edu (fix this after renaming the machine)
#iptables -A PREROUTING -t mangle -m tcp -m multiport -p tcp -d 18.181.0.29 --dports 20,21,25 -j MARK --set-mark 4
# All else to i-hate-penguins
iptables -A PREROUTING -t mangle -m mark --mark 0 -d 18.181.0.29 -j MARK --set-mark 4
iptables -A PREROUTING -t mangle -m mark --mark 0 -d 18.4.86.29 -j MARK --set-mark 4