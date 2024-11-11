# HTTP-1.1-Lib

A Go library implementing HTTP/1.1 protocol from scratch, designed for learning and understanding the protocol's inner workings.

## 📋 Overview

This educational project aims to:
- Implement HTTP/1.1 protocol from the ground up
- Provide insights into the protocol's complexities
- Serve as a learning resource for HTTP protocol internals.
- Document the journey on step by step implementation of how one can implement the HTTP/1.1 protocol for their own learning.

## 🚀 Getting Started

### Prerequisites
- Go 1.x or higher
- Git

### Directory Structure
```
├── README.md
├── cookie/
│   ├── cookie.go
├── common/
│   ├── common.go       # Shared constants and types
├── request.go          # HTTP request parsing and handling
├── response.go         # HTTP response generation
└── http-proto.go      # Contains the core functions to handle the flow of an http request.
```
### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/http-1.1-lib.git
cd http-1.1-lib
```

2. Install dependencies
```bash
go mod tidy
```

3. Run the application
```bash
go run main.go
```

## ⚙️ Configuration

The library supports configuration through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| HTTP_DOMAIN | Custom domain for the server | localhost |

## 📘 HTTP Status Codes Reference

### 1xx Informational
| Code | Constant | Description |
|------|----------|-------------|
| 100 | CONTINUE | Continue with the request |
| 101 | SWITCHING_PROTOCOLS | Protocol switching in progress |
| 102 | PROCESSING | Request is being processed |
| 103 | EARLY_HINTS | Early hints about the response |

### 2xx Success
| Code | Constant | Description |
|------|----------|-------------|
| 200 | OK | Request succeeded |
| 201 | CREATED | Resource created successfully |
| 202 | ACCEPTED | Request accepted for processing |
| 203 | NON_AUTHORITATIVE | Modified server response |
| 204 | NO_CONTENT | No content to send |
| 205 | RESET_CONTENT | Reset document view |
| 206 | PARTIAL_CONTENT | Partial resource returned |
| 207 | MULTI_STATUS | Multiple status operations |
| 208 | ALREADY_REPORTED | Resource previously reported |
| 226 | IM_USED | Instance manipulation used |

### 3xx Redirection
| Code | Constant | Description |
|------|----------|-------------|
| 300 | MULTIPLE_CHOICES | Multiple options available |
| 301 | MOVED_PERMANENTLY | Resource moved permanently |
| 302 | FOUND | Resource found elsewhere |
| 303 | SEE_OTHER | See other resource |
| 304 | NOT_MODIFIED | Resource not modified |
| 305 | USE_PROXY | Use proxy for request |
| 307 | TEMPORARY_REDIRECT | Temporary redirect |
| 308 | PERMANENT_REDIRECT | Permanent redirect |

### 4xx Client Errors
| Code | Constant | Description |
|------|----------|-------------|
| 400 | BAD_REQUEST | Bad request syntax |
| 401 | UNAUTHORIZED | Authentication required |
| 402 | PAYMENT_REQUIRED | Payment required |
| 403 | FORBIDDEN | Request forbidden |
| 404 | NOT_FOUND | Resource not found |
| 405 | METHOD_NOT_ALLOWED | Method not allowed |
| 406 | NOT_ACCEPTABLE | Not acceptable response |
| 407 | PROXY_AUTH_REQUIRED | Proxy authentication required |
| 408 | REQUEST_TIMEOUT | Request timeout |
| 409 | CONFLICT | Request conflict |
| 410 | GONE | Resource gone |
| 411 | LENGTH_REQUIRED | Length required |
| 412 | PRECONDITION_FAILED | Precondition failed |
| 413 | PAYLOAD_TOO_LARGE | Payload too large |
| 414 | URI_TOO_LONG | URI too long |
| 415 | UNSUPPORTED_MEDIA_TYPE | Unsupported media type |
| 416 | RANGE_NOT_SATISFIABLE | Range not satisfiable |
| 417 | EXPECTATION_FAILED | Expectation failed |
| 418 | IM_A_TEAPOT | I'm a teapot |
| 421 | MISDIRECTED_REQUEST | Misdirected request |
| 422 | UNPROCESSABLE_ENTITY | Unprocessable entity |
| 423 | LOCKED | Resource locked |
| 424 | FAILED_DEPENDENCY | Failed dependency |
| 425 | TOO_EARLY | Too early |
| 426 | UPGRADE_REQUIRED | Upgrade required |
| 428 | PRECONDITION_REQUIRED | Precondition required |
| 429 | TOO_MANY_REQUESTS | Too many requests |
| 431 | HEADERS_TOO_LARGE | Headers too large |
| 451 | LEGAL_REASONS | Unavailable for legal reasons |

### 5xx Server Errors
| Code | Constant | Description |
|------|----------|-------------|
| 500 | INTERNAL_SERVER_ERROR | Internal server error |
| 501 | NOT_IMPLEMENTED | Not implemented |
| 502 | BAD_GATEWAY | Bad gateway |
| 503 | SERVICE_UNAVAILABLE | Service unavailable |
| 504 | GATEWAY_TIMEOUT | Gateway timeout |
| 505 | HTTP_VERSION_NOT_SUPPORTED | HTTP version not supported |
| 506 | VARIANT_ALSO_NEGOTIATES | Variant also negotiates |
| 507 | INSUFFICIENT_STORAGE | Insufficient storage |
| 508 | LOOP_DETECTED | Loop detected |
| 510 | NOT_EXTENDED | Not extended |
| 511 | NETWORK_AUTH_REQUIRED | Network authentication required |

## 🤝 Contributing

Contributions are welcome! Feel free to:
- Open issues for bug reports or feature requests
- Submit pull requests for improvements
- Provide feedback on code organization
- Suggest documentation improvements

## ✨ Acknowledgments

- Vedant Gandhi

If you contribute feel free to add your name here.