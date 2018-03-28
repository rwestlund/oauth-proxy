# oauth-proxy
A simple proxy for your OAuth redirect URLs.

# The Problem
What do you do when your app has many different domains as redirect URLs?  This
is an issue, for example, when you have a subdomain per customer in a
multi-instance architecture.

You probably have a long list of redirect URLs with each OAuth provider you
use. When you deploy a new instance, you have to go add that URL to each OAuth
provider; sooner or later, you forget one and things break.

# The Solution
Hack the `state` parameter! The spec says it must be returned to you verbatim.
Instead of simply providing your CSRF token, URL-encode a JSON object like
this:

```json
{
    "token": "CSRF TOKEN",
    "host": "deployment23.myapp.com"
}
```

Then set the redirect URL for all your OAuth requests to where you have
`oauth-proxy` running. It will inspect the `host` parameter and redirect the
returning OAuth response to the domain you specify. Now you only need to
register one redirect URL with your OAuth providers!

You can call the token whatever you want, and include more fields if you like.
`oauth-proxy` only cares about `host`.

# Example
1. Run `oauth-proxy` at `oauth.myapp.com`
2. When someone starts OAuth at `deployment23.myapp.com`, use `oauth.myapp.com`
   as the hostname in the `redirect_url`
3. Set `state` to a URL-encoded JSON string:
   `{"token":"CSRF TOKEN","host":"deployment23.myapp.com"}`
4. When the provider returns the code, the client's browser will go to
   `oauth.myapp.com`, which will send them back to `deployment23.myapp.com`


# Building
Check out this repo in your `$GOPATH`. Run `go build`. Run the binary.

# TLS
`oauth-proxy` does not do TLS (yet). Run it behind a reverse proxy like NGINX.

# License
BSD-2-Clause
