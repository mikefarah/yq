# HCL

Encode and decode to and from [HashiCorp Configuration Language (HCL)](https://github.com/hashicorp/hcl).

HCL is commonly used in HashiCorp tools like Terraform for configuration files. The yq HCL encoder and decoder support:
- Blocks and attributes
- String interpolation and expressions (preserved without quotes)
- Comments (leading, head, and line comments)
- Nested structures (maps and lists)
- Syntax colorisation when enabled

