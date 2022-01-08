# update-plex-ipv6-access-url
DynDNS-like script for keeping your Plex IPv6 custom access URL up to date, automating the IPv6 workaround for Plex as described on [Reddit](https://www.reddit.com/r/PleX/comments/b82opu/plex_remote_access_over_ipv6/).

## Features
- determine IPv6 address for a specified interface
- update Plex config with plex.direct-domain using current IPv6 address
- restart Plex service

## Prerequisite
You need to initially add a custom server access URL manually, which the script will then keep up to date. Here's how you can easily find everything you need:

1. Navigate to any media item in your library as the admin user
2. Click `...` (More) next to the edit pen to see more options
3. Click "Get Info" at the very bottom of the list
4. At the bottom right of the "Media info" dialogue, click "View XML"
5. Take note of the domain in the new tab, you should something like `https://[some-ip-address].[server-id].plex.direct:32400/`
6. Determine your server's current IPv6 address (either directly from the server or via your router)
7. Expand the IPv6 address to it's uncompressed/full state (using a tool like [this one](https://dnschecker.org/ipv6-expand.php), which turns a shortened IPv6 such as `2606:4700::6810:84e5` into `2606:4700:0000:0000:0000:0000:6810:84e5`)
8. Replace all colons in the expanded IPv6 address with dashes
9. Put everything together as `https://["dashed"-ipv6-address].[serverid].plex.direct:32400/` the "dashed" IPv6
10. Add the url as a "Custom server access URL" in your Plex server settings unter "Network"

If you are unsure where to find the "Media info" dialogue in order to get your server id, follow [this Plex support article](https://support.plex.tv/articles/201998867-investigate-media-information-and-formats/).

__Example__:

IPv6 address reported by server: `2606:4700::6810:84e5`

Expanded IPv6 address: `2606:4700:0000:0000:0000:0000:6810:84e5`

Custom server access URL to add: `https://2606-4700-0000-0000-0000-0000-6810-84e5.055e11e51095c8b4f16c572691ff8113.plex.direct:32400/`

## Setup
In order to use the script, you just need to download it. However, if you want it to automatically restart the Plex service in case the IPv6 address changed, you need to either comment out or add a command to restart the service (depending on your host OS).

## Command line arguments
The script requires two positional arguments:

Position|Description|Required
--------|-----------|--------
1       |Name of network interface to use for IPv6 access|Yes
2       |Path to Plex config (`Preferences.xml`)|Yes

If you are unsure where to find your Plex data directory (which contains the `Preferences.xml`), please refer to [this Plex support article](https://support.plex.tv/articles/202915258-where-is-the-plex-media-server-data-directory-located/).

## Usage

A simple example: You are running Plex on an Ubuntu server and set Plex up to listen on the `ens18` interface. Your Plex library resides in the default location, which is `/var/lib/plexmediaserver/Library/Application Support/Plex Media Server`. Assuming you are currently in the directory you placed the script in, you would run the script like so:

```bash
./update-plex-ipv6-access-url.sh ens18 "/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Preferences.xml"
```

**Please note:** You need to either quote the path to the config or escape the contained spaces.

In order to automate the process, create a cronjob or other type of scheduled task in order to run the script regularly. Just keep in mind that your Plex server will be unavailable for a few seconds when the script restarts the Plex service after an IPv6 address change, so choose times to run the script accordingly.