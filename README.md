# gh-app

gh-app is a GitHub CLI extension that simplifies managing GitHub Apps by providing commands to fetch app details, list installations, and generate authentication tokens.

## Features

- Fetch app details to save along private key
- List installations for app
- Generate app JWT using private key
- Generate access token for app installation
- Customizable store path (defaults to `~/.gh-app.yaml`)

## Installation

Use `gh extension install` to install this extension: 

```shell
gh extension install github.com/cursethevulgar/gh-app
```

## Usage

`gh app list`

`gh app save -s <app slug> -p <private key path> [-a <application id>]`

`gh app jwt -s <app slug>`

`gh app installations -s <app slug>`

`gh app token -s <app slug> -i <installation id>`

