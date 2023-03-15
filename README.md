# cheapel

CLI tool to find cheapest period of electricity in the next ~36 hours, using Tibber's API. Can also optionally send a notification to your device via
Pushbullet.

## Installation

Create a config file, cheapel-config.yaml:

```
tibber-token: <your Tibber API token>
pb-token: <optional, Pushbullet token for notifications>
pb-device: <optional, device name to notify>
```

Install the binary and point to the config file.

## Usage

```
$ cheapel [-hours 3] [-notify] [-config ~/path/to/config.yaml]
```

