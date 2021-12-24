# shorten
Simple link shortener.

### Command Line Options
-  -base-url string
        Base url for shortened result (default "localhost:8090")
-  -cert-file string
        location of cert file
-  -key-file string
        location of key file
-  -log-path string
        Logs location (default "./logs.txt")
-  -port string
        port to host on (default "8090")
-  -shortened-length int
        The length of the shortned url (default 5)
### How to Run
#### Option 1: From source
- `make vendor`
- `make run`

#### Option 2: Build it
- `make vendor`
- `make build` - will build for current system architecture. 
- `make build-linux` - will build Linux distributable
- `make build-pi` - will build Raspberry Pi compatible distributable
- You will find the executable in the `./bin` folder.

#### Option 3: Docker
Linux images:
- `docker run -d jaeg/shorten:latest`

Raspberry pi images:
- `docker run -d jaeg/shorten:latest-pi`

### API

### /shorten
-  `Post`
    - Body: `{ "link": "link to shorten"}`
    - Returns: `{"link":"short link"}`

### /\<short url>
-  `Get`
    - Body: None
    - Returns: 
        - Short url exists: Status 301 to long url
        - Short url doesn't exist: `{"error":"Invalid url"}`
