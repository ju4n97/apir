server {
  host = "127.0.0.1"
  port = 8080
}

openapi {
  title       = "Acme Documentation Showcase"
  version     = "1.0.0"
  description = "Demonstration of multiple interactive documentation portals and raw OpenAPI 3.1 artifacts."

  server {
    url         = "/"
    description = "Current server origin"
  }

  tag {
    name        = "system"
    description = "Core runtime and health probes"
  }

  contact {
    name  = "API Architecture Team"
    email = "architecture@example.com"
    url   = "https://example.com/support"
  }

  license {
    name = "MIT"
    url  = "https://opensource.org/licenses/MIT"
  }
}

route "GET /openapi.json" {
  step "spec" {
    format = "json"
  }
}

route "GET /openapi.yaml" {
  step "spec" {
    format = "yaml"
  }
}

route "GET /docs" {
  step "docs" {
    renderer = "scalar"
  }
}

route "GET /docs/swagger" {
  step "docs" {
    renderer = "swagger"
  }
}

route "GET /docs/elements" {
  step "docs" {
    renderer = "elements"
  }
}

route "GET /docs/redoc" {
  step "docs" {
    renderer = "redoc"
  }
}

route "GET /docs/custom" {
  step "docs" {
    template = <<HTML
	<html>
	<head>
		<title>{{ .Title }}</title>
	</head>
	<body>
		<h1>Custom Template</h1>
		<p>
			Source Specification:
			<a href="{{ .SpecURL }}">{{ .SpecURL }}</a>
		</p>
	</body>
	</html>
	HTML
  }
}
route "GET /api/v1/ping" {
  summary = "Simple latency check"
  tag     = "system"

  step "respond" {
    status = 200
    body = {
      status = "pong"
    }
  }
}