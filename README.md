# Http-1.1-Lib

This project aims for more personal growth to learn how the HTTP protocol works under the hood, understand the complexities of the protocol and gain more insights in a step by step way. I have started this as a side project and am enthusiastic for suggestions and improvements.

To start with the development - 
1. Clone the repo to local machine.
2. Just run the command `go run main.go`.
3. Currently the socket is hardcoded at "localhost:8811" but soon will shift to env/config.
4. Things are messy now so you are welcome to refactor the code to make it more readable.
5. Feel free to open issue for any suggestion,doubt or criticism.


#### I have added support for env files. Following is a list of supported env variables - 
1. **Domain** - Set the HTTP_DOMAIN key to set your custom domain.

# Supported HTTP Status Codes
*These are the constants you can use in the code.*

## 1xx Informational

- `CONTINUE` (100)
- `SWITCHING_PROTOCOLS` (101)
- `PROCESSING` (102)
- `EARLY_HINTS` (103)

## 2xx Success

- `OK` (200)
- `CREATED` (201)
- `ACCEPTED` (202)
- `NON_AUTHORITATIVE` (203)
- `NO_CONTENT` (204)
- `RESET_CONTENT` (205)
- `PARTIAL_CONTENT` (206)
- `MULTI_STATUS` (207)
- `ALREADY_REPORTED` (208)
- `IM_USED` (226)

## 3xx Redirection

- `MULTIPLE_CHOICES` (300)
- `MOVED_PERMANENTLY` (301)
- `FOUND` (302)
- `SEE_OTHER` (303)
- `NOT_MODIFIED` (304)
- `USE_PROXY` (305)
- `TEMPORARY_REDIRECT` (307)
- `PERMANENT_REDIRECT` (308)

## 4xx Client Errors

- `BAD_REQUEST` (400)
- `UNAUTHORIZED` (401)
- `PAYMENT_REQUIRED` (402)
- `FORBIDDEN` (403)
- `NOT_FOUND` (404)
- `METHOD_NOT_ALLOWED` (405)
- `NOT_ACCEPTABLE` (406)
- `PROXY_AUTH_REQUIRED` (407)
- `REQUEST_TIMEOUT` (408)
- `CONFLICT` (409)
- `GONE` (410)
- `LENGTH_REQUIRED` (411)
- `PRECONDITION_FAILED` (412)
- `PAYLOAD_TOO_LARGE` (413)
- `URI_TOO_LONG` (414)
- `UNSUPPORTED_MEDIA_TYPE` (415)
- `RANGE_NOT_SATISFIABLE` (416)
- `EXPECTATION_FAILED` (417)
- `IM_A_TEAPOT` (418)
- `MISDIRECTED_REQUEST` (421)
- `UNPROCESSABLE_ENTITY` (422)
- `LOCKED` (423)
- `FAILED_DEPENDENCY` (424)
- `TOO_EARLY` (425)
- `UPGRADE_REQUIRED` (426)
- `PRECONDITION_REQUIRED` (428)
- `TOO_MANY_REQUESTS` (429)
- `HEADERS_TOO_LARGE` (431)
- `LEGAL_REASONS` (451)

## 5xx Server Errors

- `INTERNAL_SERVER_ERROR` (500)
- `NOT_IMPLEMENTED` (501)
- `BAD_GATEWAY` (502)
- `SERVICE_UNAVAILABLE` (503)
- `GATEWAY_TIMEOUT` (504)
- `HTTP_VERSION_NOT_SUPPORTED` (505)
- `VARIANT_ALSO_NEGOTIATES` (506)
- `INSUFFICIENT_STORAGE` (507)
- `LOOP_DETECTED` (508)
- `NOT_EXTENDED` (510)
- `NETWORK_AUTH_REQUIRED` (511)