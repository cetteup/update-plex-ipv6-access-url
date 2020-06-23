#!/bin/bash
# Retrieve host IPv6 address and update Plex custom access url accordingly
# based on workaround posted by Pikey18 on the Plex subreddit: https://www.reddit.com/r/PleX/comments/b82opu/plex_remote_access_over_ipv6/

# Get IPv6 address of given interface (command adapted from: https://superuser.com/a/1057290)
IPv6=`/sbin/ip -6 addr show dev "$1" | grep inet6 | awk -F '[ \t]+|/' '{print $3}' | grep -v ^::1 | grep -v ^fe80`
if [ -n "$IPv6" ]; then
	echo "Got IPv6 address: $IPv6"
	# Format IPv6 for Plex (replace : with -)
	PlexFormatIPv6=`echo "$IPv6" | sed -e "s/:/-/g"`
	if ! grep -q "customConnections\=\"https\:\/\/$PlexFormatIPv6" "$2" ; then
		echo "Current IPv6 does not match config, updating config"
		# Replace old IPv6 with new one
		sed -i -e "s/customConnections\=\"https\:\/\/[a-fA-F0-9\-]*/customConnections\=\"https:\/\/$PlexFormatIPv6/" "$2"
		# Restart Plex service, uncomment line for your server's os or add your own
		# systemctl restart plexmediaserver # systemd Linux distributions (Ubuntu, Debian, ...)
		# synoservice --restart pkgctl-Plex\ Media\ Server # Synology DiskStations (details: https://tech.setepontos.com/2018/03/25/control-synology-dsm-services-via-terminal-ssh/)
	else
		echo "Current IPv6 matches config, exiting"
	fi
	exit 0
else
	echo "No IPv6 address found, exiting"
	exit 1
fi
