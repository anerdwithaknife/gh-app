# gh-app

gh-app is a GitHub CLI extension that simplifies managing GitHub Apps by providing commands to fetch app details, list installations, and generate authentication tokens.

## Features

- Fetch app details to save along private key
- List installations for app
- Generate app JWT
- Generate access token for app installation
- Generate user token via OAuth flow
- Customizable store path (defaults to `~/.gh-app.yaml`)

## Installation

Use `gh extension install` to install this extension: 

```shell
gh extension install github.com/anerdwithaknife/gh-app
```

## Stored app credentials

> [!CAUTION]
> All information (including the private key) is saved in a plain text yaml file and should be treated with the same care as the original PEM files.

The default path is `~/.gh-app.yaml` which can be overridden via the `GH_APP_STORE_PATH` environment variable. 

When setting the path via environment variable, the full path including the desired filename must be used, i.e. `GH_APP_STORE_PATH=~/.config/gh/gh-app.yaml` (file can have any name).

## Usage

### View saved apps

`gh app list` or `gh app ls`

### Save app

`gh app save -s <app slug> -p <private key path> [-a <application id>]`

### Generate JWT token

`gh app jwt -s <app slug>`

### View all installations

`gh app installations -s <app slug>`

### Generate installation token

`gh app token -s <app slug> -i <installation id>`

`gh app token -s <app slug> -o <org name>`

### Initiate OAuth flow

`gh app oauth -s <app slug>`

`gh app oauth -s <app slug> -p <port>`