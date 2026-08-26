# Embeddings, reranking and platform APIs

## Vector APIs

Primary sources:

- `api-reference/text-embedding/`
- `api-reference/multimodal-embedding/`
- `api-reference/rerank/`
- `developer-guides/embeddings/`
- `resources/faq-embedding-reranking.md`

Do not mix dimensions, input batch limits, multimodal fields or reranking response shapes across models. Record exact model, dimension, normalization and truncation behavior in the integration contract.

## Platform APIs

`api-reference/platform-api/` covers models, files, batches, fine-tuning and related resources. Read the exact page before implementing CRUD or upload behavior.

- File IDs, purposes, retention and size limits are API-specific.
- Batch is not the same as client-side concurrency.
- Fine-tuning, datasets, evaluation and deployment have dedicated developer guides.

## Billing and entitlements

Billing, free quota, coupons, invoices and Token Plan are account/product surfaces rather than model payload contracts. Consult:

- `resources/billing-overview.md`
- `resources/free-quota.md`
- `resources/faq-billing.md`
- `token-plan/`

Never infer current balance, quota, regional availability or price from static documentation; verify with the authorized account before consequential calls.
