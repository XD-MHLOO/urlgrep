package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	parser "github.com/Cgboal/DomainParser"
)

var extractor parser.Parser

func init() {
	extractor = parser.NewDomainParser()
}

type urlProc func(*url.URL, *regexp.Regexp) bool

type Condition struct {
	Mode   string
	Negate bool
	Regex  *regexp.Regexp
}

var verbose bool
var rawMode bool

func main() {

	flag.BoolVar(&verbose, "v", false, "Enable verbose output for debugging")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output for debugging")
	flag.BoolVar(&rawMode, "r", false, "Treat all URL components as raw (no decoding)")
	flag.BoolVar(&rawMode, "raw", false, "Treat all URL components as raw (no decoding)")

	flag.Parse()

	args := flag.Args()

	if len(args)%2 != 0 {
		fmt.Fprintln(os.Stderr, "Expected mode-regex pairs")
		os.Exit(1)
	}
	procMap := map[string]urlProc{

		"scheme":    scheme,
		"opaque":    opaque,
		"userinfo":  userinfo,
		"domain":    domain,
		"subdomain": subdomain,
		"apex":      apex,
		"tld":       tld,
		"port":      port,
		"path":      path,
		"keypairs":  keyPairs,
		"key":       key,
		"value":     value,
		"ext":       ext,
		"fragment":  fragment,
	}

	var conditions []Condition

	for i := 0; i < len(args); i += 2 {
		mode := args[i]
		regexStr := args[i+1]

		negate := false
		if strings.HasSuffix(mode, ":not") {
			negate = true
			mode = strings.TrimSuffix(mode, ":not")
		}

		_, ok := procMap[mode]
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown mode specified: %s\n", mode)
			os.Exit(1)
		}

		re, err := regexp.Compile(regexStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid regex for mode %s: %v\n", mode, err)
			os.Exit(1)
		}

		conditions = append(conditions, Condition{
			Mode:   mode,
			Negate: negate,
			Regex:  re,
		})
	}

	scanner := bufio.NewScanner(os.Stdin)
	seen := make(map[string]bool)
	bufWriter := bufio.NewWriter(os.Stdout) // <-- create bufio.Writer
	defer bufWriter.Flush()
	for scanner.Scan() {
		str := scanner.Text()
		u, err := parseURL(str)

		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "parse failure: %s\n", err)
			}
			continue
		}
		match := true

		for _, cond := range conditions {
			procFn, ok := procMap[cond.Mode]
			if !ok {
				fmt.Fprintf(os.Stderr, "unknown mode: %s\n", cond.Mode)
				match = false
				break
			}

			result := procFn(u, cond.Regex)
			if cond.Negate {
				result = !result
			}
			if !result {
				match = false
				break
			}
		}

		if match {

			if seen[str] {
				continue
			}

			fmt.Fprintln(bufWriter, str)

			seen[str] = true

		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Scanner error: %v\n", err)
	}
}

func key(u *url.URL, re *regexp.Regexp) bool {
	var q url.Values
	if rawMode {
		q = QueryV2(u)
	} else {
		q = u.Query()
	}
	for key, _ := range q {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}
func value(u *url.URL, re *regexp.Regexp) bool {
	var q url.Values
	if rawMode {
		q = QueryV2(u)
	} else {
		q = u.Query()
	}
	for _, vals := range q {
		for _, val := range vals {
			if re.MatchString(val) {
				return true
			}
		}
	}
	return false
}

func scheme(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(u.Scheme)

}

func opaque(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(u.Opaque)

}

func userinfo(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(u.User.String())

}

func domain(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(u.Hostname())

}

func apex(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(extractor.GetDomain(u.Hostname()) + "." + extractor.GetTld(u.Hostname()))

}

func path(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(u.EscapedPath())

}

func keyPairs(u *url.URL, re *regexp.Regexp) bool {
	var q url.Values
	if rawMode {
		q = QueryV2(u)
	} else {
		q = u.Query()
	}
	for key, vals := range q {
		for _, val := range vals {
			if re.MatchString(fmt.Sprintf("%s=%s", key, val)) {
				return true
			}
		}
	}
	return false
}

func subdomain(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(extractor.GetSubdomain(u.Hostname()))

}

func tld(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(extractor.GetTld(u.Hostname()))

}

func ext(u *url.URL, re *regexp.Regexp) bool {
	paths := strings.Split(u.Path, "/")
	if len(paths) == 0 || paths[len(paths)-1] == "" {
		return false
	}
	lastSegment := paths[len(paths)-1]
	parts := strings.Split(lastSegment, ".")
	if len(parts) > 1 {
		return re.MatchString(parts[len(parts)-1])
	}
	return false
}

func port(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(u.Port())

}

func fragment(u *url.URL, re *regexp.Regexp) bool {
	return re.MatchString(u.EscapedFragment())

}

func parseQueryV2(m url.Values, query string) (err error) {
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.Contains(key, ";") {
			err = fmt.Errorf("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		// key, err1 := url.QueryUnescape(key)
		// if err1 != nil {
		// 	if err == nil {
		// 		err = err1
		// 	}
		// 	continue
		// }
		// value, err1 = url.QueryUnescape(value)
		// if err1 != nil {
		// 	if err == nil {
		// 		err = err1
		// 	}
		// 	continue
		// }
		m[key] = append(m[key], value)
	}
	return err
}

func QueryV2(u *url.URL) url.Values {
	m := make(url.Values)
	_ = parseQueryV2(m, u.RawQuery)
	return m
}
func parseURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return url.Parse("http://" + raw)
	}

	return u, nil
}

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: urlgrep [OPTIONS] [<MODE>[:not] <REGEX>] ...

Filters by matching components of URLs from stdin using regular expressions.

Sample URL:
  https://user:pass@sub.example.co.uk:8443/some/path/file.txt?foo=bar&baz=qux#section1

Modes:
  scheme     : Match URL scheme (https)
  opaque     : Match opaque component (empty in sample URL)
  userinfo   : Match user credentials (user:pass)
  domain     : Match full domain (sub.example.co.uk)
  subdomain  : Match subdomain only (sub)
  apex       : Match apex domain (example.co.uk)
  tld        : Match top-level domain (co.uk)
  port       : Match port number (8443)
  path       : Match URL path (/some/path/file.txt)
  keypairs   : Match key=value query pairs (foo=bar, baz=qux)
  key        : Match query parameter keys (foo, baz)
  value      : Match query parameter values (bar, qux)
  ext        : Match file extension (txt)
  fragment   : Match fragment identifier (section1)

Negation:
  Add ":not" to any mode to invert the match (like grep -v).
  Example: urlgrep ext:not "css"    # Excludes URLs with the extension 'css'

Examples:
  cat urls.txt | urlgrep scheme "^https$"
  cat urls.txt | urlgrep domain "example\.co\.uk" path "^/some/path" # Combine multiple modes (AND logic)
  cat urls.txt | urlgrep ext:not "css"


Options:
`)
		flag.PrintDefaults()
	}
}
