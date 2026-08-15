# DeepSeek credentials (mixin)

Minimal kit for the **Hub + kits** proof: network allowlist + proxy-managed
`DEEPSEEK_API_KEY`. No agent install.

```bash
# Store on the host (service id must match the kit: deepseek)
sbx secret set deepseek

# Example recipe (stock shell, no local image):
sbx-kit run shell --yes
```

Inside the box the env var should be a sentinel; the proxy injects the real
key on outbound calls to `api.deepseek.com`.
