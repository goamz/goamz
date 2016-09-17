package aws

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

// SignerFunc represents a function that can sign a http.Request with the given Auth details.
type SignerFunc func(*http.Request, Auth, string) error

type V2Signer struct {
	auth     Auth
	endpoint string
}

var b64 = base64.StdEncoding

func NewV2Signer(auth Auth, endpoint string) (*V2Signer, error) {
	return &V2Signer{auth: auth, endpoint: endpoint}, nil
}

func (s *V2Signer) Sign(method, path string, params map[string]string) {
	req, _ := NewRequest(s.endpoint, method, path, params)
	SignV2(req, s.auth, "") // service is not used so just set blank
}

// SignV2 signs req using the AWS version 2 signature with the given credentials.
// serviceName is present to so that it can be used as a SignerFunc and is not used
func SignV2(req *http.Request, auth Auth, serviceName string) (err error) {
	vals := req.URL.Query()

	vals.Set("AWSAccessKeyId", auth.AccessKey)
	vals.Set("SignatureVersion", "2")
	vals.Set("SignatureMethod", "HmacSHA256")
	if auth.Token() != "" {
		vals.Set("SecurityToken", auth.Token())
	}

	// AWS specifies that the parameters in a signed request must
	// be provided in the natural order of the keys. This is distinct
	// from the natural order of the encoded value of key=value.
	// Percent and Equals affect the sorting order.
	var keys, sarray []string
	for k, _ := range vals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sarray = append(sarray, Encode(k)+"="+Encode(vals.Get(k)))
	}

	path := req.URL.Path
	if path == "" {
		path = "/"
	}

	joined := strings.Join(sarray, "&")
	payload := req.Method + "\n" + req.Host + "\n" + path + "\n" + joined
	hash := hmac.New(sha256.New, []byte(auth.SecretKey))
	hash.Write([]byte(payload))
	signature := make([]byte, b64.EncodedLen(hash.Size()))
	b64.Encode(signature, hash.Sum(nil))

	vals.Set("Signature", string(signature))

	req.URL.RawQuery = vals.Encode()

	return nil
}

// Common date formats for signing requests
const (
	ISO8601BasicFormat      = "20060102T150405Z"
	ISO8601BasicFormatShort = "20060102"
)

type Route53Signer struct {
	auth Auth
}

func NewRoute53Signer(auth Auth) *Route53Signer {
	return &Route53Signer{auth: auth}
}

// Creates the authorize signature based on the date stamp and secret key
func getHeaderAuthorize(date string, auth Auth) string {
	hmacSha256 := hmac.New(sha256.New, []byte(auth.SecretKey))
	hmacSha256.Write([]byte(date))

	return base64.StdEncoding.EncodeToString(hmacSha256.Sum(nil))
}

// Adds all the required headers for AWS Route53 API to the request
// including the authorization
func (s *Route53Signer) Sign(req *http.Request) {
	req.Header.Set("Host", req.Host)
	req.Header.Set("Content-Type", "application/xml")
	SignRoute53(req, s.auth, "")
}

// SignRoute53 signs req using the AWS Route53 signature with the given credentials.
// In contrast to Route53Signer.Sign this only sets the X-Amzn-Authorization and X-Amz-Date headers.
// serviceName is present to so that it can be used as a SignerFunc and is not used
func SignRoute53(req *http.Request, auth Auth, serviceName string) error {
	resp, err := http.Get("https://route53.amazonaws.com/date")
	if resp != nil {
		resp.Body.Close()
	}

	if err != nil {
		return err
	}

	date := resp.Header.Get("Date")
	authHeader := fmt.Sprintf("AWS3-HTTPS AWSAccessKeyId=%s,Algorithm=%s,Signature=%s",
		auth.AccessKey, "HmacSHA256", getHeaderAuthorize(date, auth))

	req.Header.Set("X-Amzn-Authorization", authHeader)
	req.Header.Set("X-Amz-Date", date)

	return nil
}

/*
The V4Signer encapsulates all of the functionality to sign a request with the AWS
Signature Version 4 Signing Process. (http://goo.gl/u1OWZz)
*/
type V4Signer struct {
	auth        Auth
	serviceName string
	regionName  string
}

// SignV4Region is an adapter to allow Signer to be used with SignV4 given a regionName.
func SignV4Region(regionName string) SignerFunc {
	return func(req *http.Request, auth Auth, serviceName string) error {
		return SignV4(req, auth, serviceName, regionName)
	}
}

// SignV4 signs the req with auth for regionName
func SignV4(req *http.Request, auth Auth, serviceName, regionName string) error {
	return NewV4Signer(auth, serviceName, regionName).Sign(req)
}

/*
Return a new instance of a V4Signer capable of signing AWS requests.
*/
func NewV4Signer(auth Auth, serviceName, regionName string) *V4Signer {
	return &V4Signer{auth: auth, serviceName: serviceName, regionName: regionName}
}

/*
Sign a request according to the AWS Signature Version 4 Signing Process. (http://goo.gl/u1OWZz)

The signed request will include an "x-amz-date" header with a current timestamp if a valid "x-amz-date"
or "date" header was not available in the original request. In addition, AWS Signature Version 4 requires
the "host" header to be a signed header, therefore the Sign method will manually set a "host" header from
the request.Host.

The signed request will include a new "Authorization" header indicating that the request has been signed.

Any changes to the request after signing the request will invalidate the signature.
*/
func (s *V4Signer) Sign(req *http.Request) error {
	req.Header.Set("Host", req.Host)     // host header must be included as a signed header
	t := s.requestTime(req)              // Get requst time
	creq, err := s.canonicalRequest(req) // Build canonical request
	if err != nil {
		return err
	}
	sts := s.stringToSign(t, creq)                    // Build string to sign
	signature := s.signature(t, sts)                  // Calculate the AWS Signature Version 4
	auth := s.authorization(req.Header, t, signature) // Create Authorization header value
	req.Header.Set("Authorization", auth)             // Add Authorization header to request
	return nil
}

/*
requestTime method will parse the time from the request "x-amz-date" or "date" headers.
If the "x-amz-date" header is present, that will take priority over the "date" header.
If neither header is defined or we are unable to parse either header as a valid date
then we will create a new "x-amz-date" header with the current time.
*/
func (s *V4Signer) requestTime(req *http.Request) time.Time {

	// Get "x-amz-date" header
	date := req.Header.Get("x-amz-date")

	// Attempt to parse as ISO8601BasicFormat
	t, err := time.Parse(ISO8601BasicFormat, date)
	if err == nil {
		return t
	}

	// Attempt to parse as http.TimeFormat
	t, err = time.Parse(http.TimeFormat, date)
	if err == nil {
		req.Header.Set("x-amz-date", t.Format(ISO8601BasicFormat))
		return t
	}

	// Get "date" header
	date = req.Header.Get("date")

	// Attempt to parse as http.TimeFormat
	t, err = time.Parse(http.TimeFormat, date)
	if err == nil {
		return t
	}

	// Create a current time header to be used
	t = time.Now().UTC()
	req.Header.Set("x-amz-date", t.Format(ISO8601BasicFormat))
	return t
}

/*
canonicalRequest method creates the canonical request according to Task 1 of the AWS Signature Version 4 Signing Process. (http://goo.gl/eUUZ3S)

    CanonicalRequest =
      HTTPRequestMethod + '\n' +
      CanonicalURI + '\n' +
      CanonicalQueryString + '\n' +
      CanonicalHeaders + '\n' +
      SignedHeaders + '\n' +
      HexEncode(Hash(Payload))
*/
func (s *V4Signer) canonicalRequest(req *http.Request) (string, error) {
	c := new(bytes.Buffer)

	// Precalculate hash as it can add a header needed by canonicalHeaders
	hash, err := s.payloadHash(req)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(c, "%s\n", req.Method)
	fmt.Fprintf(c, "%s\n", s.canonicalURI(req.URL))
	fmt.Fprintf(c, "%s\n", s.canonicalQueryString(req.URL))
	fmt.Fprintf(c, "%s\n\n", s.canonicalHeaders(req.Header))
	fmt.Fprintf(c, "%s\n", s.signedHeaders(req.Header))
	fmt.Fprintf(c, "%s", hash)
	return c.String(), nil
}

func (s *V4Signer) canonicalURI(u *url.URL) string {
	canonicalPath := u.RequestURI()
	if u.RawQuery != "" {
		canonicalPath = canonicalPath[:len(canonicalPath)-len(u.RawQuery)-1]
	}
	slash := strings.HasSuffix(canonicalPath, "/")
	canonicalPath = path.Clean(canonicalPath)
	if canonicalPath != "/" && slash {
		canonicalPath += "/"
	}
	return canonicalPath
}

func (s *V4Signer) canonicalQueryString(u *url.URL) string {
	return strings.Replace(u.Query().Encode(), "+", "%20", -1)
}

func (s *V4Signer) canonicalHeaders(h http.Header) string {
	i, a := 0, make([]string, len(h))
	for k, v := range h {
		for j, w := range v {
			v[j] = strings.Trim(w, " ")
		}
		sort.Strings(v)
		a[i] = strings.ToLower(k) + ":" + strings.Join(v, ",")
		i++
	}
	sort.Strings(a)
	return strings.Join(a, "\n")
}

func (s *V4Signer) signedHeaders(h http.Header) string {
	i, a := 0, make([]string, len(h))
	for k, _ := range h {
		a[i] = strings.ToLower(k)
		i++
	}
	sort.Strings(a)
	return strings.Join(a, ";")
}

// payloadHash returns the payload hash for the req.
// This will return the value of X-Amz-Content-Sha256 if present otherwise
// it will be calculated from the req details and the X-Amz-Content-Sha256 set.
//
// Due to the fact this can add a new header it must be called before any call
// to canonicalHeaders, otherwise the canonicalHeaders will be missing the added
// header.
func (s *V4Signer) payloadHash(req *http.Request) (string, error) {
	hash := req.Header.Get("X-Amz-Content-Sha256")
	if hash == "" {
		var b []byte
		if req.Body == nil {
			b = []byte("")
		} else {
			var err error
			b, err = ioutil.ReadAll(req.Body)
			if err != nil {
				return "", err
			}
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		hash = s.hash(string(b))
		req.Header.Add("X-Amz-Content-Sha256", hash)
	}

	return hash, nil
}

/*
stringToSign method creates the string to sign accorting to Task 2 of the AWS Signature Version 4 Signing Process. (http://goo.gl/es1PAu)

    StringToSign  =
      Algorithm + '\n' +
      RequestDate + '\n' +
      CredentialScope + '\n' +
      HexEncode(Hash(CanonicalRequest))
*/
func (s *V4Signer) stringToSign(t time.Time, creq string) string {
	w := new(bytes.Buffer)
	fmt.Fprint(w, "AWS4-HMAC-SHA256\n")
	fmt.Fprintf(w, "%s\n", t.Format(ISO8601BasicFormat))
	fmt.Fprintf(w, "%s\n", s.credentialScope(t))
	fmt.Fprintf(w, "%s", s.hash(creq))
	return w.String()
}

func (s *V4Signer) credentialScope(t time.Time) string {
	return fmt.Sprintf("%s/%s/%s/aws4_request", t.Format(ISO8601BasicFormatShort), s.regionName, s.serviceName)
}

/*
signature method calculates the AWS Signature Version 4 according to Task 3 of the AWS Signature Version 4 Signing Process. (http://goo.gl/j0Yqe1)

	signature = HexEncode(HMAC(derived-signing-key, string-to-sign))
*/
func (s *V4Signer) signature(t time.Time, sts string) string {
	h := s.hmac(s.derivedKey(t), []byte(sts))
	return fmt.Sprintf("%x", h)
}

/*
derivedKey method derives a signing key to be used for signing a request.

	kSecret = Your AWS Secret Access Key
    kDate = HMAC("AWS4" + kSecret, Date)
    kRegion = HMAC(kDate, Region)
    kService = HMAC(kRegion, Service)
    kSigning = HMAC(kService, "aws4_request")
*/
func (s *V4Signer) derivedKey(t time.Time) []byte {
	h := s.hmac([]byte("AWS4"+s.auth.SecretKey), []byte(t.Format(ISO8601BasicFormatShort)))
	h = s.hmac(h, []byte(s.regionName))
	h = s.hmac(h, []byte(s.serviceName))
	h = s.hmac(h, []byte("aws4_request"))
	return h
}

/*
authorization method generates the authorization header value.
*/
func (s *V4Signer) authorization(header http.Header, t time.Time, signature string) string {
	w := new(bytes.Buffer)
	fmt.Fprint(w, "AWS4-HMAC-SHA256 ")
	fmt.Fprintf(w, "Credential=%s/%s, ", s.auth.AccessKey, s.credentialScope(t))
	fmt.Fprintf(w, "SignedHeaders=%s, ", s.signedHeaders(header))
	fmt.Fprintf(w, "Signature=%s", signature)
	return w.String()
}

// hash method calculates the sha256 hash for a given string
func (s *V4Signer) hash(in string) string {
	h := sha256.New()
	fmt.Fprintf(h, "%s", in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// hmac method calculates the sha256 hmac for a given slice of bytes
func (s *V4Signer) hmac(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
