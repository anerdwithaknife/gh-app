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

## Environment

`GH_APP_ID` should contain the GitHub app id (not client id)

`GH_APP_PRIVATE_KEY` should contain the full contents of the pem file


## Tables

```
Apps
    Name        string
    Slug        string
    AppId       int
    ClientId    string
    PrivateKey  string
```

```
Installations
    Id              int
    TargetId        int
    TargetType      string
    AccountName     string
    AppId           int
```

Apps.Installations = Installations[]

## Usage

`gh app store <app slug>` (fetching + manual mode)

`gh app jwt <app slug|app id|client id>`

`gh app installations <app slug>`

`gh app token <installation id|app slug+org>`

`gh app oauth <installation id|app slug+org>` (render links/callback)