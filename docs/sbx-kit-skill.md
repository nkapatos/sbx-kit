# sbx-kit skill

This file is a **reference copy**. The canonical template lives in the sbx-kit CLI:

```bash
sbx-kit recipes skill              # print SKILL.md
sbx-kit recipes skill --cursor     # install for Cursor
sbx-kit recipes skill -o docs/sbx-kit-skill.md
```

Regenerate after changing `cli/internal/recipecreate/templates/skill.md.tmpl`.

The skill covers **two lands**: CLI overlay (`/etc/sbx-kit/context.md` in the box)
vs user catalog kits (optional `agentContext` may *point at* overlay docs).

See also: `sbx-kit recipes create <dir>` scaffolds a catalog bundle with `AGENTS.md`
for directory-specific **catalog** guidance (not box runtime docs).
