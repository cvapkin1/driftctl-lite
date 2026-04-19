# driftctl-lite

Lightweight CLI to detect infrastructure drift between Terraform state and live cloud resources.

---

## Installation

```bash
go install github.com/your-org/driftctl-lite@latest
```

Or build from source:

```bash
git clone https://github.com/your-org/driftctl-lite.git
cd driftctl-lite
go build -o driftctl-lite .
```

---

## Usage

Point `driftctl-lite` at your Terraform state file and let it compare against your live cloud environment:

```bash
driftctl-lite scan --state terraform.tfstate --provider aws --region us-east-1
```

**Example output:**

```
[DRIFT DETECTED] aws_s3_bucket.my-bucket — tag "Environment" changed: "staging" → "production"
[DRIFT DETECTED] aws_security_group.web — ingress rule added outside Terraform
[OK]             aws_vpc.main — no drift detected

Summary: 2 drifted resources, 1 clean
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--state` | Path to Terraform state file | `terraform.tfstate` |
| `--provider` | Cloud provider (`aws`, `gcp`) | `aws` |
| `--region` | Cloud region to scan | `us-east-1` |
| `--output` | Output format (`text`, `json`) | `text` |

---

## Requirements

- Go 1.21+
- Valid cloud credentials (e.g. `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY`)

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE)