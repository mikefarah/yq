# HCL

Encode and decode to and from [HashiCorp Configuration Language (HCL)](https://github.com/hashicorp/hcl).

HCL is commonly used in HashiCorp tools like Terraform for configuration files. The yq HCL encoder and decoder support:
- Blocks and attributes
- String interpolation and expressions (preserved without quotes)
- Comments (leading, head, and line comments)
- Nested structures (maps and lists)
- Syntax colorisation when enabled


## Parse HCL
Given a sample.hcl file of:
```hcl
io_mode = "async"
```
then
```bash
yq -oy sample.hcl
```
will output
```yaml
io_mode: "async"
```

## Roundtrip: Sample Doc
Given a sample.hcl file of:
```hcl
service "cat" {
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }

  process "management" {
    command = ["/usr/local/bin/awesome-app", "management"]
  }
}

```
then
```bash
yq sample.hcl
```
will output
```hcl
service "cat" {
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }
  process "management" {
    command = ["/usr/local/bin/awesome-app", "management"]
  }
}
```

## Roundtrip: With an update
Given a sample.hcl file of:
```hcl
service "cat" {
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }

  process "management" {
    command = ["/usr/local/bin/awesome-app", "management"]
  }
}

```
then
```bash
yq '.service.cat.process.main.command += "meow"' sample.hcl
```
will output
```hcl
service "cat" {
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server", "meow"]
  }
  process "management" {
    command = ["/usr/local/bin/awesome-app", "management"]
  }
}
```

## Parse HCL: Sample Doc
Given a sample.hcl file of:
```hcl
service "cat" {
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }

  process "management" {
    command = ["/usr/local/bin/awesome-app", "management"]
  }
}

```
then
```bash
yq -oy sample.hcl
```
will output
```yaml
service:
  cat:
    process:
      main:
        command:
          - "/usr/local/bin/awesome-app"
          - "server"
      management:
        command:
          - "/usr/local/bin/awesome-app"
          - "management"
```

## Parse HCL: with comments
Given a sample.hcl file of:
```hcl
# Configuration
port = 8080 # server port
```
then
```bash
yq -oy sample.hcl
```
will output
```yaml
# Configuration
port: 8080 # server port
```

## Roundtrip: with comments
Given a sample.hcl file of:
```hcl
# Configuration
port = 8080
```
then
```bash
yq sample.hcl
```
will output
```hcl
# Configuration
port = 8080
```

## Roundtrip: With templates, functions and arithmetic
Given a sample.hcl file of:
```hcl
# Arithmetic with literals and application-provided variables
sum = 1 + addend

# String interpolation and templates
message = "Hello, ${name}!"

# Application-provided functions
shouty_message = upper(message)
```
then
```bash
yq sample.hcl
```
will output
```hcl
# Arithmetic with literals and application-provided variables
sum = 1 + addend
# String interpolation and templates
message = "Hello, ${name}!"
# Application-provided functions
shouty_message = upper(message)
```

## Roundtrip: Separate blocks with same name.
Given a sample.hcl file of:
```hcl
resource "aws_instance" "web" {
  ami = "ami-12345"
}
resource "aws_instance" "db" {
  ami = "ami-67890"
}
```
then
```bash
yq sample.hcl
```
will output
```hcl
resource "aws_instance" "web" {
  ami = "ami-12345"
}
resource "aws_instance" "db" {
  ami = "ami-67890"
}
```

