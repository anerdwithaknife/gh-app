# gh-app
Extension for gh cli to simplify handling of GitHub apps

## Features

- [x] Fetch app details
- [x] Save app to database (sqlite)
- [x] List apps
- [x] List installations
- [x] Generate app JWT
- [ ] Generate installation token
- [ ] Customizable db path

## Usage

`gh app list`

`gh app save -s <app slug> -p <private key path>`

`gh app jwt -s <app slug>`

`gh app installations -s <app slug>`

`gh app token -s <app slug>`

