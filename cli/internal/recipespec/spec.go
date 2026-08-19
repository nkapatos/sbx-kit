package recipespec

// Status describes the recipe spec situation for experimental tooling.
const Status = `Recipe spec: not finalized (parked).

Interim manifest: <catalog>/<dir>/recipes/agents.yaml
  defaults: resources, kits[]
  agents.<name>: sbx_agent, kits[], image_name, template_fallback, stub

Formal schemaVersion and recipes.yaml naming will be defined when the spec stabilizes.`
