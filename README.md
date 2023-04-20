# update-plex-ipv6-access-url

DynDNS-like tool for keeping your Plex IPv6 custom access URL up to date, automating the IPv6 workaround for Plex as described on [Reddit](https://www.reddit.com/r/PleX/comments/b82opu/plex_remote_access_over_ipv6/).

## Features

- determine IPv6 address for a specified interface
- update Plex settings with plex.direct-domain using current IPv6 address

## Command line arguments

If any required command line argument is omitted, the tool will prompt you to provide input at runtime.

| Name      | Description                                                                                                                                            | Required               |
|-----------|--------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------|
| address   | Plex server's address in format http\[s\]://host:port                                                                                                  | Yes                    |
| interface | Name of network interface to use for IPv6 access                                                                                                       | Yes                    |
| token     | Plex access token (X-Plex-Token) [How to find](https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/)               | If config is not given |
 | config    | Path to Plex config (Preferences.xml) [How to find](https://support.plex.tv/articles/202915258-where-is-the-plex-media-server-data-directory-located/) | No                     |

## Usage

A simple example: You are running Plex on an Ubuntu server and set Plex up to listen on the `ens18` interface. Your Plex library resides in the default location, which is `/var/lib/plexmediaserver/Library/Application Support/Plex Media Server`. Assuming you are currently in the directory you placed the script in, you would run the script like so:
```bash
./update-plex-ipv6-access-url -address http://localhost:32400 -interface ens18 -config "/var/lib/plexmediaserver/Library/Application Support/Plex Media Server/Preferences.xml"
```

Plex does not store the config on disk on all platforms. On those platforms, you need to manually specify an access token, as it cannot be read from Preferences.xml.
```powershell
.\update-plex-ipv6-access-url.exe -address http://localhost:32400 -interface Ethernet -token your-X-Plex-Token
```

To automate the process, create a cronjob or other type of scheduled task in order to run the script regularly.

For testing or one-time use, you can also run the tool without any command line arguments and provide required input at runtime.
```commandline
$ ./update-plex-ipv6-access-url
Enter the Plex server's address in format 'http[s]://host:port': http://localhost:32400
Enter the name of network interface to use for IPv6 access: eth0
Enter a Plex access token (X-Plex-Token): your-X-Plex-Token
```
