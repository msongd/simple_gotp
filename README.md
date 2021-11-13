# simple_gotp

A very simple TOTP (Time based OTP - compatible with Google Authenticator app) management web app. Only one binary to deploy

## Model
 - Username 1 <-> n Token
 - Admin user can create / delete username
 - User can create token
 - User can choose active token
 - User can add aliases.

## Authenticate / Authorization
 - Use external OpenID Connect for authentication (currently only keycloak is supported), username is from field preferred_username in jwt
 - Use config claim_realm_role to specify role of user to allow access
 - Use config claim_role_admin to specify admin role of user
 - Admin can create/delete username
 - After logging (having jwt), username is obtained from field preferred_username in jwt, and then used as a filter to select same username

## Backend
 - Save to json file when exit

## Endpoint
 - API endpoints are listed in router.go and are pretty self explain

  | VERB | URL  | DATA | Notes |
  | ---- | ---- | ---- | ----- |
  | DELETE | /user/{user}/alias/{alias} | | |
	| POST   | /user/{user}/alias | { alias: "alias" } | create new alias "alias" |
	| GET | /user/{user} | | |
	| DELETE | /user/{user} | | |
	| POST | /user/{user} | { active_token: "tokenId" } | set active token |
	| GET | /user | | |
	| POST | /user | { username: "newUsername" } | |
  | POST | /qr/{user}/{token} | {} | obtain QR image (contain secret) of a token, return { img: "base64encoded" } |
	| GET | /otp/{user}/{token} | | |
	| GET | /otp | | get all current totp of users |
	| POST | /token/{user}/import | { url: "newTokenUrl" } | import existing TOTP token |
	| DELETE | /token/{user}/{token} | | |
	| GET | /token/{user} | | |
	| POST | /token/{user} | { issuer: "newTokenIssuer"} | create new TOTP with issuer name "newTokenIssuer" |

## Usage
 - `simple_gotp -f config.json`

## Sample config with explaination
```
{
  "listen":":8080", // ip:port to listen
  "log_dir": "./log", // where to store log file: app.log and access.log, auto rotated.
  "data_file": "data/data.json", // where to store/load data as a json file (db)
  "frontend_dir": "./frontend", // where to serve alternate front end
  "use_embedded_frontend": true, // use embedded front end 
  "ssl_key":"", // tls private key for https
  "ssl_cert":"", // cert for https
  "keycloak_cfg" : { // if configured, user must login to access /auth/xxxx endpoints
    "auth_url": "https://xxx/auth", // url to redirect to for login
    "realm":"My-Realm", //keycloak realm name
    "client_id":"simple_gotp_test", //keycloak client id
    "jwk_url": "https://xxx/auth/realms/My-Realm/protocol/openid-connect/certs", //jwk location
    "claim_iss": "https://xxx/auth/realms/My-Realm", // claim issuer of jwt
    "claim_realm_role": "otp_user", // keycloak role for users
    "claim_role_admin": "otp_admin" // keycloak role for admin (admin should also have user role)
  }
}
```
