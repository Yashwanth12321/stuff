# Url shortner that uses go as server
- go - server
- gin - go http framwork
- postgres- main db for storing url mappings
- redis- cache frequent url clicks with expiry

# installation
```
  git clone https://github.com/Yashwanth12321/goshortner
  go mod download
  cd goshortner
```

- have postgres docker image running at default port
- have redis docker image at default port

  ```
  go run main.go
  ```

  API Endpoints Summary
  
Endpoint

Method

Description

/shorten

POST

Shortens a long URL

/:shortURL

GET

Redirects to the original URL

/analytics

GET

Fetches URL usage statistics

# for ui
```
  cd ui-admin
  npm i
  npm run dev

```
