# noglucrawl

Crawler to extract particular information from specified websites. 

It's probably misnamed since its really a scraper of certain data.

Configuration is in the configuration file, which can be `json`, `toml`, `yaml`, and `xml`.

## Logging
Logging is currently provided by Seelog and its configuration file is `seelog.xml`. This may change in the future as Seelog is a bit much and causes problems in testing. However it is very convenient for debugging, which is why it hasn't been replaced yet.

It's replacement will be `log` and a wrapped `log` for verbose output. Such a switch would also make it possbile to support logging to tmp on startup and either discarding the log or moving it to the specified log destination after all configuration has been loaded, including command-line or environment variables.

## Usage
Don't please. It's not currently usable.
