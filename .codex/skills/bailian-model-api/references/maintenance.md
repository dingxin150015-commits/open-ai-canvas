# Skill maintenance

## Source mirror

Default local source:

`<project-root>/百炼千问文档/百炼千问文档`

The mirror is ignored by Git and is not bundled into the skill. The generated index captures the current file paths and SHA-256 values.

## Refresh

After the official mirror changes:

```powershell
python scripts/build_source_index.py --output-dir references
```

Then review:

- added, removed and hash-changed documents;
- model IDs and dated snapshots;
- endpoint, auth and region changes;
- sync/async lifecycle changes;
- parameter limits and defaults;
- result URL retention and pricing/entitlement wording.

Do not update distilled references solely from filenames or changelog headlines. Read every changed exact API page and its create/query counterpart.

## Validation

- Run the skill validator.
- Run both helper scripts with UTF-8 output.
- Search at least one text, image, video, speech and vector model.
- Confirm every path linked from `SKILL.md` exists.
- Do not perform a real API call as a validation shortcut.

## Project synchronization

The versioned source of this skill is `<project-root>/.codex/skills/bailian-model-api`. The installed global copy is `<codex-home>/skills/bailian-model-api`.

When updating, modify and validate the project copy first, then replace the global copy as one reviewed operation. Do not edit the generated source catalog by hand.
