server {
  host = "127.0.0.1"
  port = 8080
}

telemetry {
  service_name = "user-validation-service"
  logging {
    level  = "info"
    redact = ["request.headers.x-api-key"]
  }
}

schema "UserCreate" {
  field "email" {
    type        = "string"
    format      = "email"
    required    = true
    description = "Primary user contact address"
  }
  field "username" {
    type        = "string"
    required    = true
    min_length  = 3
    max_length  = 20
    description = "Unique alphanumeric handle"
  }
  field "account_type" {
    type        = "string"
    required    = true
    enum        = ["individual", "business"]
    description = "Account billing classification"
  }
  field "age" {
    type        = "integer"
    min         = 18
    max         = 120
    description = "Legal age verification"
  }
  field "role" {
    type        = "string"
    default     = "member"
    enum        = ["admin", "member", "viewer"]
    description = "Authorization tier"
  }
  field "tags" {
    type        = "array"
    description = "User interest classifications"
  }
}

route "GET /docs" {
  step "docs" {
    renderer = "scalar"
  }
}

route "GET /openapi.json" {
  step "spec" {
    format = "json"
  }
}

route "POST /api/v1/users" {
  summary = "Create user account"
  tag     = "users"

  request {
    header "x-api-key" {
      type        = "string"
      format      = "uuid"
      required    = true
      description = "Client authorization key"
    }
    query "source" {
      type        = "string"
      default     = "direct"
      enum        = ["direct", "referral", "ad"]
      description = "User registration channel"
    }
    body = schema.UserCreate
  }

  step "respond" {
    status = 201
    body = {
      message    = "User validated and registered"
      user       = ctx.request.body
      api_key    = ctx.request.headers.x-api-key
      source     = ctx.request.query.source
      created_at = now()
    }
  }
}