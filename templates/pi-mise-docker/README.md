# Template: pi-mise-docker

**Tag:** `local/sbx-pi-mise-docker:latest` (sbx often shows `docker.io/local/…`)  
**Parent:** `local/sbx-shell-mise-docker:latest`  
**Adds:** Node (pinned) + `@earendil-works/pi-coding-agent` at `/usr/local/bin/pi`  
**Agent:** `pi` (sandbox kit)  
**CLI:** `sbx-kit run --agent pi`

## Why this is harder than plain Docker

Docker Sandboxes **templates** customize an existing agent family (cursor,
shell, …). They do **not** register a new agent runtime. A BYO harness like Pi
is a **sandbox kit** whose entrypoint is `pi`, backed by an image that already
contains that binary. Extending `shell-docker` means sbx may warn:

```text
template was built for the "shell" agent but you are using "pi"
```

That warning is expected. The hard failure `agent binary "pi" not found` means
the image loaded into sbx’s store does not have `/usr/local/bin/pi` — usually a
stale/partial import or parent image missing at build time.

## Load (order matters)

```bash
sbx-kit template load --engine docker shell-mise-docker
sbx-kit template load --engine docker pi-mise-docker
# load runs: docker run --rm --entrypoint which <tag> pi
sbx template ls | grep pi-mise

sbx rm -f sbx-tests   # if a previous create left a broken box
sbx-kit run --agent pi --yes
```

Host check without sbx:

```bash
docker run --rm --entrypoint which local/sbx-pi-mise-docker:latest pi
docker run --rm --entrypoint cat local/sbx-pi-mise-docker:latest /etc/sbx-kit-agent
```

Oh My Pi (`omp`) is out of scope for this example tree.
