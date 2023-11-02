# Simple File Manager

---

## Disclaimer:
This is a simple file manager I made to learn more about Go and it is not perfect! 

## Features:
- Front-end which shows a simple overview of the files saved with some management options
- Back-end which can handle:
  - File upload
  - File download
  - File deletion
- Endpoints can also be easily made to require authorization (see the AuthMiddleware function in main.go)

## Pre-requisites:
- Golang (tested on 1.21.3)
- Git (or just download the zip from Github directly)

## How to use:
- Clone the repo (`git clone git@github.com:aesthetic0001/FileUploader.git`
- Install the dependencies (`go get -d ./...`)
- Run the server through the run.sh script (`sh run.sh`)
- Front-end is available at `localhost:8080` by default
- Default permissions only allow uploading and delete if a key is provided in the Authorization header, you can create them in the auto-generated data/apikeys.json folder. Key file should be in the format of:
```json
{
  "keys": {
    "key1": true,
    "key2": true,
    ...
  }
}
```

## Contributions:
Feel free to make a pull request if you want to add something or fix something!