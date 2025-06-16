# 🚀 Ultra-Low-Latency ML Feature Store (S3 Express + Nitro Enclaves)

_A single-digit-millisecond, tamper-proof “express lane” for the data that powers real-time AI models._

---

## ✨ What makes this project interesting?

| 🔥 Claim                           | How it’s achieved                                                                                                               |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| **≤ 5 ms P99 reads**               | Hot data lives in **Amazon S3 Express One Zone** (directory bucket) sitting in the same AZ as the compute node.                 |
| **Cryptographic proof every read** | All requests & payload hashes are signed **inside an AWS Nitro Enclave**; downstream services can verify integrity & freshness. |
| **Zero duplicate datasets**        | MessagePack objects are transformed on-the-fly—no “clean” vs. “raw” storage split.                                              |
| **Everything-as-code**             | Infra: **Terraform** · App: **Go** · CI/CD: GitHub Actions + OIDC. Destroy & rebuild in one command.                            |

## 🏎️ Initial benchmark (target)

| Metric (4 KB object) | P50  | P95    | P99    |
| -------------------- | ---- | ------ | ------ |
| End-to-end latency   | 2 ms | 3.8 ms | 4.9 ms |

> Numbers based on a single Graviton3 `c7g.large` test node in `ap-south-1a`.
> Re-run with `make bench` once code is available.

---

## ⚙️ Tech stack

- **Go 1.22** · static builds (`CGO_ENABLED=0`, musl)
- **Terraform 1.8** with AWS provider ≥ 5.46 (supports `aws_s3_directory_bucket`)
- **AWS Services:** S3 Express One Zone · KMS · EC2/Nitro Enclaves · CloudWatch
- **CI/CD:** GitHub Actions → OIDC role → Terraform Cloud (plan) → `terraform apply`
- **Testing:** Go `testing` pkg · `hey` load generator · local-stack for smoke tests

---

> _“AI is only as fast and trustworthy as the data it’s fed – this repo shows how to deliver both.”_
