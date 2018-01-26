# StatHat Command-Line Interface

## Installation

Choose one of the following methods:

### macOS

    brew install stathat/cmd/stathat

### Binary releases

[Download latest binary releases for macOS, linux, windows](https://github.com/stathat/cmd/releases/latest)

### Compile it

Do you have a go compiler installed?

    go get github.com/stathat/cmd/stathat

## Usage

Once installed, there will be a `stathat` command.  The first time you use it, youu will need to configure your access keys
by running

    stathat setup

and following the prompts.

Then you can do things like

    stathat list

to get a list of all the stats in your account.

Or 

    stathat search "request time" 

to get a list of stats containing "request time" in their name.

You can get a summary of the data for a stat with

    stathat summary "api call"

or

    stathat summary JxJz

with a stat ID.

The `--tf` flag allows you to specify any timeframe:

    stathat summary "api call" --tf 1M1h

Datasets at any timeframe have similar flags:

    stathat dataset "devices registered"  --tf 4h5m

And all of these commands take `--json` or `--csv` flags to generate 
JSON or CSV output.
