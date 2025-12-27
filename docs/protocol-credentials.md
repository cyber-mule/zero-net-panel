# User Credential Design

## Goals
- One credential per user (unique identity).
- Not coupled to subscription; only active subscriptions can be rendered.
- Manual rotation only (user-initiated or admin-initiated).
- No plaintext storage (encrypted seed + fingerprint for audit).

## Data Model
Table: `user_credentials`
- `user_id` (unique per active credential, historical versions retained)
- `version`, `status` (`active`/`deprecated`/`revoked`)
- `secret_ciphertext`, `secret_nonce`, `master_key_id`
- `fingerprint` (HMAC-SHA256)
- `issued_at`, `deprecated_at`, `revoked_at`, `last_seen_at`

## Identity Derivation
- A per-user encrypted seed is generated and stored.
- Account ID and password are derived from the seed + version.
- Derived fields exposed to templates:
  - `user_identity.account_id`
  - `user_identity.account`
  - `user_identity.password`
  - aliases: `user_identity.id`, `user_identity.uuid`, `user_identity.username`, `user_identity.secret`

## Rotation Flow
1. User or admin triggers rotate endpoint.
2. Current credential marked `deprecated`.
3. New seed created and stored as `active` with incremented `version`.

## Subscription Rendering
- If subscription status is not `active`, nodes and credentials are omitted.
- If active, `user_identity` is injected into the template context.
- `subscription.token` remains for compatibility, but should not be used for auth.

## Configuration
`Credentials.MasterKey` is required. A per-user encryption key is derived from this root to keep user secrets isolated.
