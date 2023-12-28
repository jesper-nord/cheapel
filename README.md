# cheapel

CLI tool to find cheapest period of electricity in the next ~36 hours, using Tibber's API. 
Optionally send a notification with the result to the Tibber app on your device.

## Usage
```bash
$ cheapel -token <tibber-api-token> [-hours 3] [-notify]
```

```bash
$ cheapel -help
    -token string
        Tibber API token
    -hours int
        Length of period to check for (hours) (default 1)
    -notify
        Send a Tibber notification with result
```
