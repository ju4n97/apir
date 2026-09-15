server {
  host = "127.0.0.1"
  port = 8080
}

telemetry {
  service_name = "zero-dependency-api"
  logging {
    level  = "info"
    format = "text"
  }
}

route "GET /openapi.json" {
  step "spec" {
    format = "json"
  }
}

route "GET /docs" {
  step "docs" {
    renderer = "scalar"
  }
}

route "GET /api/v1/health" {
  summary = "System health check"
  tag     = "system"

  step "starlark" "sysinfo" {
    source = <<-STARLARK
      def execute(ctx):
          return {
              "status": "healthy",
              "engine": "esquema",
              "timestamp": ctx["timestamp"],
          }
    STARLARK
  }

  step "respond" {
    status = 200
    body   = steps.sysinfo.result
  }
}

route "POST /api/v1/sanitize" {
  summary = "String normalization and list deduplication"
  tag     = "utilities"

  request {
    body {
      field "prefix" {
        type    = "string"
        default = "tag"
      }
      field "tags" {
        type     = "array"
        required = true
      }
    }
  }

  step "starlark" "format_tags" {
    source = <<-STARLARK
      def execute(ctx):
          body = ctx["request"]["body"] or {}
          prefix = body.get("prefix", "tag")
          raw_tags = body.get("tags", [])
          cleaned = list(set([
              prefix + ":" + t.strip().lower()
              for t in raw_tags
              if len(t.strip()) > 0
          ]))
          return {
              "count": len(cleaned),
              "tags": cleaned,
          }
    STARLARK
  }

  step "respond" {
    status = 200
    body   = steps.format_tags.result
  }
}