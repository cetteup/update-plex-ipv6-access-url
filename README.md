# update-plex-ipv6-access-url
DynDNS-like script for keeping your Plex IPv6 custom access URL up to date, automating the IPv6 workaround for Plex as described on [Reddit](https://www.reddit.com/r/PleX/comments/b82opu/plex_remote_access_over_ipv6/).

## Features
- determine IPv6 address for a specified interface
- update Plex config with plex.direct-domain using current IPv6 address
- restart Plex service

## Setup
In order to use the script, you just need to download it. However, if you want it to automatically restart the Plex service in case the IPv6 address changed, you need to either comment out or add a command to restart the service (depending on your host OS).

## Command line arguments
The script requires two positional arguments:

Position|Description|Required
--------|-----------|--------
1       |Name of network interface to use for IPv6 access|Yes
2       |Path to Plex config (`Preferences.xml`)|Yes

## Usage

A simple example: You are running Plex on an Ubuntu server and set Plex up to listen on the `ens18` interface. Your Plex library resides in the default location, which is `/var/lib/plexmediaserver/Library/Application Support/Plex Media Server`. Assuming you are currently in the directory you placed the script in, you would run the script like so:

```bash
./update-plex-ipv6-access-url.sh ens18 "/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Preferences.xml"
```

**Please note:** You need to either quote the path to the config or escape the contained spaces.

In order to automate the process, create a cronjob or other type of scheduled task in order to run the script regularly. Just keep in mind that your Plex server will be unavailable for a few seconds when the script restarts the Plex service after an IPv6 address change, so choose times to run the script accordingly.