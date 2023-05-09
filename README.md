# go-oauth2-example
This example will show you how you can create an OAuth2 server. In this case, only for playing with the code I choose the Fiber framework, but the truth is that the oauth2 library that I used here was written to be used with the standard net/http and that is not the case with Fiber.

These are the URLs available in this project and the purpose of each one:

Is only to show that the server is running right:
`
http://localhost:3000/
`

Will return the credentials (the client ID and the client secret) necessary to get the token:
`
http://localhost:3000/credentials
`

Will return the token related to a client ID (you need pass also the client secret):
`
http://localhost:3000/token?grant_type=client_credentials&client_id=CLIENT_ID&client_secret=CLIENT_SECRET
`

Example: (in this case your client ID is: 1081abb7, and your client secret is: 3f64ebed)
`
http://localhost:3000/token?grant_type=client_credentials&client_id=1081abb7&client_secret=3f64ebed
`

This URL is protected, you only should have access if you have the right token:
`
http://localhost:3000/protected?access_token=TOKEN
`
Example: (in this case your token is: NWE1YJQ0NMUTMZVIMY0ZNMRLLWJJNJKTNMI5MMZIMJRKMZUX)
`
http://localhost:3000/protected?access_token=NWE1YJQ0NMUTMZVIMY0ZNMRLLWJJNJKTNMI5MMZIMJRKMZUX
`

If you want to test the server, you need to clone this repository and in the root folder run the server with this command:
```bash
go run ./main.go
```

This example was written for this post on my blog: ['Build your own oauth2 server'](https://jackgris.github.io/goscrapy-blog/post/build-your-own-oauth2-server/)
