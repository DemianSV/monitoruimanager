# Polyxena Monitoring System UIManager (MonitorUIManager) 2.0.0

The **Polyxena Monitoring System** project is an attempt to create a simple infrastructure and application monitoring system based on Zabbix's best practices while addressing its main weaknesses in terms of scalability and data storage.  
**UIManager** is a component that implements a full-fledged user WEB interface for the monitoring system's functions.

<img width="960" alt="Monitor01" src="https://github.com/user-attachments/assets/1d523339-7a69-4218-a129-64d3b5a3d07b" />

## Installing UIManager
It is recommended to install **UIManager** after installing **DataManager**.

### Preparing the Environment
Download the latest release or clone the source code from the **master** branch:
```sh
git clone https://github.com/DemianSV/monitoruimanager.git
```

### Installing UIManager from a Docker Image
**Building and Running**  
Navigate to the directory with the cloned repository.

Configure the environment variables and launch parameters in the *docker-compose.yml* file.

Build the image and start **UIManager** by running:
```sh
docker compose up --build -d
```

Verify the result:
```sh
docker ps -a
docker logs monitoruimanager01
```
Try opening **UIManager** in a browser. For a local installation, the address will be: http://localhost:7007. If installed on a remote server, adjust the URL accordingly.

**Stopping**  
```sh
docker compose down
```

### Installing UIManager Directly on a Server
**Building and Running**  
Install the latest version of GoLang (if not already installed) from the official website (https://go.dev).  
Install the latest version of NodeJS from the official website (https://nodejs.org/en/download).

Navigate to the directory with the cloned repository.

Build **UIManager**:
```sh
cd frontend
npm install -D vite
npm install
npm run build
cp -r dist/assets ../
cp dist/index.html ../
cd ../
go build
```

Configure the parameters in the *monitoruimanager.json* file.

Run **UIManager**:
```sh
./monitoruimanager
# or
./monitoruimanager -log file
```

### Configuration Parameters
Configuration can be done in two ways:
1. Via the configuration file (uimanager.json)
2. Via environment variables (for the Docker option)

>**Environment variables have higher priority and will always override parameters from the configuration file!**

>**Variable naming convention:**  
>**MONITOR**: Indicates the variable belongs to the monitoring system,  
>**HTTP**: Group of variables for the HTTP scope,  
>**HOST**: Configuration parameter.

**MONITOR_TLS_CERTPATH**: Path to the certificate for UI Manager connections,  
**MONITOR_TLS_KEYPATH**: Path to the private key for UI Manager connections,  
**MONITOR_TLS_KEYPASSWORD**: The password for the private key,  
**MONITOR_TLS_CAPATH**: Path to the CA certificate for HTTP server connections.  

**MONITOR_UIMANAGER_MODULEID**: Reserved module ID for UIManager,  
**MONITOR_UIMANAGER_SPACEID**: Reserved space ID for UIManager,  
**MONITOR_UIMANAGER_DATAMANAGERURL**: Array of DataManager API URLs.

**MONITOR_HTTP_HOST**: Address for incoming HTTP server connections (API, Swagger),  
**MONITOR_HTTP_PORT**: Port for incoming HTTP server connections,  
**MONITOR_HTTP_READTIMEOUT**: HTTP server read timeout,  
**MONITOR_HTTP_WRITETIMEOUT**: HTTP server write timeout,  
**MONITOR_HTTP_RL**: Requests Limit, HTTP server incoming request limit per second,  
**MONITOR_HTTP_TLS**: Enable TLS for HTTP server (true/false),  
**MONITOR_HTTP_TLS_SKIPVERIFY**: Disable certificate verification for TLS connections (true/false),
**MONITOR_HTTP_TRUSTPROXY**: Trust `X-Forwarded-Proto` from a reverse proxy (true/false; only enable when direct access to UIManager is restricted to trusted proxies).

**MONITOR_JWT_SECRET**: Persistent JWT signing secret (at least 32 bytes; keep the same value across application instances),
**MONITOR_JWT_ALLOW_INSECURE_HTTP**: Allow password operations over plain HTTP (testing only; defaults to false).

When TLS terminates at a reverse proxy, keep `MONITOR_JWT_ALLOW_INSECURE_HTTP=false` and set `MONITOR_HTTP_TRUSTPROXY=true`. The proxy must overwrite `X-Forwarded-Proto` with `$scheme`, and network/firewall rules must prevent clients from reaching UIManager except through that trusted proxy. This setting trusts the header; it does not authenticate the proxy by itself.

**MONITOR_DB_HOST**: Array of database node addresses, comma-separated (e.g., 127.0.0.1:9042,127.0.0.1:9043),  
**MONITOR_DB_USERNAME**: Database username,  
**MONITOR_DB_PASSWORD**: Database password,  
**MONITOR_DB_TLS**: Enable TLS for database cluster interactions (true/false),  
**MONITOR_DB_TLS_SKIPVERIFY**: Disable certificate verification for TLS connections (true/false),  
**MONITOR_DB_KEYSPACE**: KeySpace for the connection,  
**MONITOR_DB_CONSISTENCY**: Consistency parameter for database data,  
**MONITOR_DB_CONSISTENCYREAD**: Consistency parameter for database queries,  
**MONITOR_DB_TIMEOUT**: Query execution timeout before forced termination.

**Configuration via uimanager.json File**  
>Parameter groups and parameters in the file are identical to those in the environment variables.  

>When using the configuration file, remember that its parameters have lower priority than environment variables and will be overridden by them!  

>Combining parameters from different sources is allowed, keeping the above in mind.  

```json
{
	"TLS": {
		"CERTPATH": "",
		"KEYPATH": "",
		"KEYPASSWORD": "",
		"CAPATH": ""
	},
	"UIMANAGER": {
		"MODULЕID": "623a57c2-3df5-4287-ada5-82f7e4a0b5db",
		"SPACEID": "17f0bd20-41cf-4801-a481-ff721a41fa93",
		"DATAMANAGERURL": ["https://172.30.0.10:6006"]
	},
	"HTTP": {
		"HOST": "0.0.0.0",
		"PORT": "7007",
		"READTIMEOUT": 30,
		"WRITETIMEOUT": 30,
		"RL": 100,
		"TLS": false,
		"SKIPVERIFY": true
	},
	"JWT": {
		"SECRET": "you-jwt-secret",
		"ALLOWINSECUREHTTP": false
	},
	"DB": {
		"HOSTS": ["127.0.0.1:19042", "127.0.0.1:29042"],
		"USERNAME": "monitor",
		"PASSWORD": "",
		"TLS": true,
		"SKIPVERIFY": true,
		"KEYSPACE": "monitor",
		"CONSISTENCY": "Quorum",
		"CONSISTENSYREAD": "One",
		"TIMEOUT": 30000
	}
}
```

### Access and refresh tokens
The UI uses the following authentication endpoints:

- `POST /auth` accepts the `username` and original `password` form fields over HTTPS, returns a 15-minute access token in JSON, and sets a 7-day refresh token in an `HttpOnly`, `SameSite=Strict` cookie.
- `POST /auth/refresh` rotates the refresh token and returns a new access token.
- `POST /auth/logout` revokes the presented refresh token and clears its cookie.

Protected API routes accept access tokens only in the `Authorization: Bearer` header. Access tokens must contain the `token_use=access` claim. The legacy single-token `POST /auth` behavior and `GET /api/v1/exit` endpoint have been removed; `POST /auth` now issues the access/refresh token pair. Legacy JWT cookies are not accepted for API authentication.

Passwords are sent as entered by the user and must be protected in transit by HTTPS. For TLS termination at a reverse proxy, the proxy must remove any client-supplied `X-Forwarded-Proto` header and set its own value to `https`. Do not trust this header from untrusted clients. The `MONITOR_JWT_ALLOW_INSECURE_HTTP=true` setting bypasses this transport requirement and is only for isolated HTTP test environments; never enable it in production.

New and changed passwords are stored using Argon2id on the server (minimum 8 Unicode characters, maximum 1024 UTF-8 bytes). Existing client-side SHA-256 password values are verified during login and upgraded to Argon2id after a successful login. This migration is lazy: accounts that do not log in remain on the legacy hash until they next authenticate or have their password reset.

Before using these endpoints, manually apply [`migrations/001_create_auth_refresh_sessions.cql`](./migrations/001_create_auth_refresh_sessions.cql) and [`migrations/002_create_auth_refresh_sessions_by_user.cql`](./migrations/002_create_auth_refresh_sessions_by_user.cql) with the application's configured Cassandra keyspace selected. The configured database user must be allowed to read, insert, update, and delete rows in both tables.

Changing a user's password revokes their outstanding refresh sessions. Existing access tokens remain valid until their short (15-minute) expiry.

The user lookup table is introduced after the initial refresh-session migration. Before deploying code that uses it, either allow pre-existing refresh sessions to expire (up to 7 days) or invalidate them during a maintenance window; sessions created before the lookup table was populated cannot be found by username.

Set `MONITOR_JWT_SECRET` to a persistent random value of at least 32 bytes, shared by all application instances. If it is omitted, the server generates a temporary secret; tokens will not survive a restart and access-token verification will not work consistently across instances.

## Using Grafana (1.0.5, 1.1.0)
Create a user in the Cassandra database, for example **Grafana**, and give out the right only to read for the table **monitor.raw_data**.  
Install Grafana in any convenient way using the documentation from the official website (https://grafana.com/docs/grafana/latest/setup-grafana/installation/).  
Turn on and configure Data Sources Plugin: **Apache Cassandra Datasource for Grafana**.

<img width="1453" height="243" alt="Снимок экрана 2025-08-30 в 18 25 56" src="https://github.com/user-attachments/assets/ac9cc4f9-ea1e-4f4b-a314-7f161624c361" />

Set up a plugin to connect to a database with an previously established UZ.  
Create a dashboard on the metrics stored in the database using an example of a request:

```sql
SELECT space_id, value, totimestamp(maxtimeuuid(event_time)), space_description FROM monitor.raw_data WHERE space_id IN (65c6b051-10fb-4bd7-8c04-7fe478e55d13, 34bcd935-8e7a-4b79-b76b-352cf3ece91f) AND metric = 'system.cpu.util' AND synt_key IN ('${__from:date:YYYY.MM}', '${__to:date:YYYY.MM}') AND event_time > $__from and event_time < $__to
```

<img width="1134" height="874" alt="Снимок экрана 2025-09-13 в 10 31 56" src="https://github.com/user-attachments/assets/a2143064-11a2-4fc7-aa2e-f34a35e063c5" />

<img width="1469" height="820" alt="Снимок экрана 2025-09-13 в 10 26 15" src="https://github.com/user-attachments/assets/b2ac7f9f-b7b8-454f-af06-ed785bb50f3b" />