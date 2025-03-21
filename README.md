# Uptime Kuma widget for Waybar

This is a simple program to display a summary of [Uptime Kuma](https://github.com/louislam/uptime-kuma) status in a Waybar module.

It displays a green checkmark if all monitors are up, or the number of monitors Up, Pending or Down in green, yellow and red numbers respectively.

You can also use this program to output detailed monitor statuses in json format via the cli with `--format=json`.
Zero dependencies, this program only uses the go standard library.

## Usage with Waybar

![](docs/assets/waybar.png)

In your `waybar/config` file, add the following:
```json
"custom/kuma-waybar": {
    "exec": "kuma-waybar --format=waybar --env=$HOME/.config/waybar/.kuma-waybar.env",
    "interval": 60,
    "on-click": "kuma-waybar open --env=$HOME/.config/waybar/.kuma-waybar.env",
    "max-length": 40,
    "format": "🐻 {}"
},
```

Note: The default format is ANSI for the CLI, specify --format=waybar in the `exec` command.

Clicking on the module will open the Uptime Kuma dashboard in your default browser.
Interval is in seconds, so it will update automatically every 60 seconds.

Choose what order to place the custom element:
Here is an example of placing it to the left in between some other elements:
```json
"modules-left": [
    "custom/exit",
    "custom/kuma-waybar",
    "sway/workspaces"
],
```

See the [Waybar repo](https://github.com/Alexays/Waybar) for more information on Waybar config.

## Installation

### From Releases

```bash
wget https://github.com/WebTender/kuma-waybar/releases/download/v1.0.1/kuma-waybar-linux_x86_64 -O ~/.local/bin/kuma-waybar`
chmod +x ~/.local/bin/kuma-waybar
sudo cp /home/brandon/.local/bin/kuma-waybar /usr/local/bin
```
> Check [Releases](https://github.com/WebTender/kuma-waybar/releases) to substitute the latest download link
> Assumes `~/.local/bin/` is in your $PATH

### Install from Source
```bash
git clone https://github.com/WebTender/kuma-waybar.git
cd kuma-waybar

# Builds and installs at /usr/local/bin/kuma-waybar
make install

# Optionally remove the source code
cd ../
rm -rf kuma-waybar

# Now you can run the binary
kuma-waybar --help
```
Tip: For security you can review the [install make script](./Makefile) before running it.

## Configuration

Default configuration can be set in `~/.kuma-waybar.env`:
```env
UPTIME_KUMA_API_KEY=your-api-key
UPTIME_KUMA_BASE_URL=https://your-uptime-kuma-instance.com
```

> You may also provide a local `.env` in the working directory
> You can also provide the `--env` argument to the script to point to the `.env` file.
> E.g. `kuma-waybar --env=$HOME/.config/waybar/.kuma-waybar.env`

**Where to get the API Key?**

You can get your API from `/settings/api-keys` of your Uptime Kuma instance or under Settings -> API Keys in Uptime Kuma.

Supports a few `--format` options:
- `--format=plain` - Outputs the uptime summary in no color formatting.
- `--format=waybar` - Outputs the uptime summary in a format that Waybar can display.
- `--format=ansi` - (default) Outputs the uptime summary in a format that is easy to read in the cli.
- `--format=json` - Outputs the uptime details in a json format.
- `--format=jsonp` - Outputs the uptime details in a json format with indentation.

## CLI Usage
```bash
kuma-waybar
```

![](docs/assets/partial-down.png)

See also [Full JSON Output](#full-json-output) for more detailed monitor status.

You may want an alias for CLI usage, you could add the following to your `.bashrc` or `.zshrc`:
```bash
alias kuma='kuma-waybar'
```
Then simply:
```bash
kuma
kuma list
```

## Full JSON Output

```bash
kuma-waybar --format=json > kuma-status.json
```
Tip: You can pipe into `jq` to parse the json output.

## Unit Tests

```bash
make test
# or
go test ./..
```

## Uptime Kuma Supported Version

This project is tested with both v1 and v2 beta of [Uptime Kuma](https://github.com/louislam/uptime-kuma)

Tested versions:
- [1.23.16](https://github.com/louislam/uptime-kuma/releases/tag/1.23.16)
- [2.0.0-beta.1](https://github.com/louislam/uptime-kuma/releases/tag/2.0.0-beta.1)

## License

MIT License

This script was written for our own purposes. Feel free to use without any restrictions. There is no warranty or guarantee of any kind.

## Contributing

PRs will be considered.

If you make use of this project, consider a BTC donation:
`BC1QZ3QQLZ5LK89DRKAU5N4HG47R9GFEMLYAKUCAMD`
