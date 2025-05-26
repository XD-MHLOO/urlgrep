# urlgrep

Filter URLs like grep, but smarter — by matching specific URL components.

---

## Install 📦 
```bash
go install github.com/XD-MHLOO/urlgrep@latest
```
You can also download pre-compiled binaries from the [Releases page](https://github.com/XD-MHLOO/urlgrep/releases).
## Syntax 🔧
```
urlgrep <mode>[:not] <regex>
```
- \<mode\>  the part of the URL to match (e.g. domain, path, ext, port, key, etc.)
- :not — (Optional) Negates the match
- \<regex\>  a regular expression pattern to match


## Usage 🛠️
urlgrep reads URLs from stdin and filters them by matching different components using regex.

```
▶ cat urls.txt
https://sub.example.co.uk:1234/user/uid123/profile
https://sub.example.com/shop/electronics/laptops/?sort=new_price&rate=5
http://dev.sub.example.com/api/v1/users/?page=1
https://admin:pass@example.com/a/b/c/debug.php.bak?user=admin
https://anotherexample.com/admin/dashboard/howdly.gif?debug=true#debugPanel
```

### scheme
Match the URL scheme (e.g., http, https, ftp, etc.)

Example: Find URLs using HTTPS
```
▶ cat urls.txt | urlgrep scheme "^https$"
https://sub.example.co.uk:1234/user/uid123/profile
https://sub.example.com/shop/electronics/laptops/?sort=new_price&rate=5
https://admin:pass@example.com/a/b/c/debug.php.bak?user=admin
https://anotherexample.com/admin/dashboard/howdly.gif?debug=true#debugPanel
```


### domain
Match the full domain including subdomains (e.g., sub.example.com).

Example: Find URLs with domain example.com and all its subdomains.
```
▶ cat urls.txt | urlgrep domain "(^|\.)example\.com$"
https://sub.example.com/shop/electronics/laptops/?sort=new_price&rate=5
http://dev.sub.example.com/api/v1/users/?page=1
https://admin:pass@example.com/a/b/c/debug.php.bak?user=admin
```
### path 
Match the URL path (e.g., /user/profile, /shop/electronics).

Example: Find URLs with user in the path:
```
▶ cat urls.txt | urlgrep path "user"
https://sub.example.co.uk:1234/user/uid123/profile
http://dev.sub.example.com/api/v1/users/?page=1
```
### key
Match query parameter names (e.g., debug in ?debug=true).

Example: Find URLs with the query parameter called debug:
```
▶ cat urls.txt | urlgrep key "^debug$"
https://anotherexample.com/admin/dashboard/howdly.gif?debug=true#debugPanel
```
### value
Match query parameter values (e.g., true in ?debug=true).

Example: Find URLs where the query value contains price:
```
▶ cat urls.txt | urlgrep value "price"
https://sub.example.com/shop/electronics/laptops/?sort=new_price&rate=5
```
### keypairs
Match full query key=value pairs (e.g., debug=true).

Example: Find URLs with keypairs debug=true:
```
▶ cat urls.txt | urlgrep keypairs "debug=true"
https://anotherexample.com/admin/dashboard/howdly.gif?debug=true#debugPanel
```
### fragment
Match the URL fragment (the part after #, e.g., section1 in #section1).

Example: Find URLs with fragment debugPanel
```
▶ cat urls.txt | urlgrep fragment "debugPanel"
https://anotherexample.com/admin/dashboard/howdly.gif?debug=true#debugPanel
```
### ext
Match URLs by their file extension (e.g., php in .php)

Example: Find URLs with .bak extension
```
▶ cat urls.txt | urlgrep ext "^bak$"
https://admin:pass@example.com/a/b/c/debug.php.bak?user=admin
```
### opaque
Match the opaque component of a URL (rarely used; often empty unless the URL is in non-hierarchical format like mailto: or data:).

Example: Find data URLs
```
▶ echo "data:text/plain;base64,SGVsbG8gd29ybGQ=" | urlgrep opaque "base64"
data:text/plain;base64,SGVsbG8gd29ybGQ=
```
### userinfo
Match user credentials in the URL (e.g., admin:pass in admin:pass@example.com).

Example: Find URLs with embedded credentials
```
▶ cat urls.txt | urlgrep userinfo ".+"
https://admin:pass@example.com/a/b/c/debug.php.bak?user=admin
```

### subdomain
Match only the subdomain part of the domain (e.g., dev.sub in dev.sub.example.com).

Example: Find URLs with subdomain contains dev

```
▶ cat urls.txt | urlgrep subdomain "dev"
http://dev.sub.example.com/api/v1/users/?page=1
```
### apex
Match the apex (root) domain without subdomains (e.g., example.com).

Example: Find URLs from example.com or example.co.uk
```
▶ cat urls.txt | urlgrep apex "example\.co\.uk"
https://sub.example.co.uk:1234/user/uid123/profile
```
### tld
Match only the top-level domain (e.g., com, co.uk, org).

Example: Find URLs with .com
```
▶ cat urls.txt | urlgrep tld "^com$"
https://sub.example.com/shop/electronics/laptops/?sort=new_price&rate=5
http://dev.sub.example.com/api/v1/users/?page=1
https://admin:pass@example.com/a/b/c/debug.php.bak?user=admin
https://anotherexample.com/admin/dashboard/howdly.gif?debug=true#debugPanel
```
### port
Match the port number in the URL (e.g., 8080).

Example: Find URLs that has port number
```
▶ cat urls.txt | urlgrep port ".+"
https://sub.example.co.uk:1234/user/uid123/profile
```

## 🔄 Negating a Match (:not)
You can reverse the match for any mode using :not. This works like grep -v, meaning it excludes URLs that match your pattern.

Example: Exclude .gif file extensions
```
▶ cat urls | urlgrep ext:not "gif"
https://sub.example.co.uk:1234/user/uid123/profile
https://sub.example.com/shop/electronics/laptops/?sort=new_price&rate=5
http://dev.sub.example.com/api/v1/users/?page=1
https://admin:pass@example.com/a/b/c/debug.php.bak?user=admin
```
You can apply :not to any mode — path, value, keypairs, fragment, ext, domain, etc.

## 🧩 Combine Multiple Filters
You can chain multiple filters to refine your search even more.
```
urlgrep [<mode>[:not] <regex>] [<mode>[:not] <regex>] [<mode>[:not] <regex>] ...
```
Filters are applied together, so a URL must match all conditions

Example: Find URLs from example.com but extension is not png/css/gif.
```
▶ cat urls.txt | urlgrep domain "example\.com" ext:not "(png|css|gif)"
https://sub.example.com/shop/electronics/laptops/?sort=new_price&rate=5
http://dev.sub.example.com/api/v1/users/?page=1
https://admin:pass@example.com/a/b/c/debug.php.bak?user=admin
```

## 🧪 Sample URL Breakdown
Use this example to understand how each mode applies to URL components:
```
https://user:pass@sub.example.co.uk:8443/some/path/file.txt?foo=bar&baz=qux#section1

```
| 🧩 Mode       | Description                         | Matches from Sample                      |
|--------------|-------------------------------------|------------------------------------------|
| 🔗 `scheme`   | URL scheme                          | `https`                                  |
| 🔒 `userinfo` | Credentials                         | `user:pass`                              |
| 🌍 `domain`   | Full domain incl. subdomains        | `sub.example.co.uk`                      |
| 🧷 `subdomain`| Just the subdomain                  | `sub`                                    |
| 🏠 `apex`     | Apex/root domain                    | `example.co.uk`                          |
| 🏁 `tld`      | Top-level domain                    | `co.uk`                                  |
| 🔢 `port`     | Port number                         | `8443`                                   |
| 🛣️ `path`     | URL path                            | `/some/path/file.txt`                    |
| 🧮 `keypairs` | Full key=value pairs                | `foo=bar`, `baz=qux`                     |
| 🗝️ `key`      | Query parameter names               | `foo`, `baz`                             |
| 💬 `value`    | Query parameter values              | `bar`, `qux`                             |
| 🧾 `ext`      | File extension                      | `txt`                                    |
| 🔖 `fragment` | Fragment identifier                 | `section1`                               |
| 🧃 `opaque`   | Opaque section (e.g., data/mailto)  | *(empty in this case)*                   |


## Optional Global Flags ⚙️
```
urlgrep [flag] [<mode>[:not] <regex>]
```
### -v, -verbose

Enable verbose output for debugging.

### -r, -raw

By default, `key`, `value`, and `keypairs` are **URL-decoded** for consistency and simplicity.

Examples:

```
▶ cat urls.txt
https://example.com/search?q=hello%20world
https://example.com/search?q=hello+world
```
```
▶ cat urls.txt | urlgrep value "hello world"
https://example.com/search?q=hello%20world
https://example.com/search?q=hello+world
```
If you want **precise matching** of encoded values, use `-r`:
```
▶ cat urls.txt | urlgrep -r value "hello%20world"
https://example.com/search?q=hello%20world
```
```
▶ cat urls.txt | urlgrep -r value "hello\+world"
https://example.com/search?q=hello+world
```

#### ⚠️ -r, -raw only affects key, value, and keypairs modes.
**Paths are never decoded (e.g., /x%2fy/z stays as-is) to avoid breaking the URL structure—this is always the default behavior and cannot be changed..**
