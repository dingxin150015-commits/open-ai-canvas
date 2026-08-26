import { Database } from "bun:sqlite";
import crypto from "node:crypto";
import fs from "node:fs";
import path from "node:path";

const [databasePath, keyPath] = process.argv.slice(2);
if (!databasePath || !keyPath) {
  console.error(
    "usage: bun scripts/verify-stage6-backup.mjs <database> <settings-key>",
  );
  process.exit(2);
}

const database = new Database(path.resolve(databasePath), {
  readonly: true,
  strict: true,
});
const key = fs.readFileSync(path.resolve(keyPath));
if (key.length !== 32) throw new Error("settings key length is invalid");

const integrityRows = database.query("PRAGMA integrity_check").all();
const integrityCheck = integrityRows
  .map((row) => String(Object.values(row)[0]))
  .join(",");
const foreignKeyViolations = database
  .query("PRAGMA foreign_key_check")
  .all().length;
const tableCount = Number(
  database
    .query(
      "SELECT count(*) AS value FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'",
    )
    .get().value,
);

const countTables = [
  "users",
  "auth_sessions",
  "system_settings",
  "model_channels",
  "channel_models",
  "logical_models",
  "resources",
  "tasks",
  "assets",
  "canvas_projects",
  "user_oss_settings",
];
const existingTables = new Set(
  database
    .query("SELECT name FROM sqlite_master WHERE type = 'table'")
    .all()
    .map((row) => String(row.name)),
);
const counts = Object.fromEntries(
  countTables
    .filter((table) => existingTables.has(table))
    .map((table) => [
      table,
      Number(
        database.query(`SELECT count(*) AS value FROM ${table}`).get().value,
      ),
    ]),
);

let channelModelMigration;
if (existingTables.has("channel_models")) {
  const columns = database
    .query("PRAGMA table_info(channel_models)")
    .all()
    .map((row) => String(row.name));
  const requiredColumns = [
    "provider_model_key",
    "support_status",
    "support_reason",
    "catalog_source",
    "catalog_version",
    "supported_operations_json",
    "documentation_paths_json",
  ];
  const hasCatalogColumns = requiredColumns.every((column) =>
    columns.includes(column),
  );
  channelModelMigration = { requiredColumnsPresent: hasCatalogColumns };
  if (hasCatalogColumns) {
    channelModelMigration.supportStatusCounts = Object.fromEntries(
      database
        .query(
          "SELECT support_status AS name, count(*) AS value FROM channel_models WHERE deleted_at IS NULL GROUP BY support_status ORDER BY support_status",
        )
        .all()
        .map((row) => [String(row.name), Number(row.value)]),
    );
    channelModelMigration.catalogSourceCounts = Object.fromEntries(
      database
        .query(
          "SELECT catalog_source AS name, count(*) AS value FROM channel_models WHERE deleted_at IS NULL GROUP BY catalog_source ORDER BY catalog_source",
        )
        .all()
        .map((row) => [String(row.name), Number(row.value)]),
    );
    channelModelMigration.missingProviderModelKeys = Number(
      database
        .query(
          "SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND (provider_model_key IS NULL OR provider_model_key = '')",
        )
        .get().value,
    );
    channelModelMigration.missingOperationsJSON = Number(
      database
        .query(
          "SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND (supported_operations_json IS NULL OR supported_operations_json = '')",
        )
        .get().value,
    );
    channelModelMigration.missingDocumentationJSON = Number(
      database
        .query(
          "SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND (documentation_paths_json IS NULL OR documentation_paths_json = '')",
        )
        .get().value,
    );
  }
}

function decryptSettingSecret(value) {
  if (typeof value !== "string" || value.length === 0)
    return { encrypted: false, decrypted: false };
  if (!value.startsWith("enc:v1:"))
    return { encrypted: false, decrypted: false };
  const payload = Buffer.from(value.slice("enc:v1:".length), "base64");
  if (payload.length < 12 + 16) return { encrypted: true, decrypted: false };
  const nonce = payload.subarray(0, 12);
  const ciphertext = payload.subarray(12, -16);
  const tag = payload.subarray(-16);
  const decipher = crypto.createDecipheriv("aes-256-gcm", key, nonce);
  decipher.setAuthTag(tag);
  const plaintext = Buffer.concat([
    decipher.update(ciphertext),
    decipher.final(),
  ]);
  const decrypted = plaintext.length > 0;
  plaintext.fill(0);
  return { encrypted: true, decrypted };
}

function collectSecrets(value, output = []) {
  if (!value || typeof value !== "object") return output;
  for (const [name, child] of Object.entries(value)) {
    if (
      name === "accessKeySecret" &&
      typeof child === "string" &&
      child.length > 0
    )
      output.push(child);
    else collectSecrets(child, output);
  }
  return output;
}

const ossPayloads = [];
if (existingTables.has("system_settings")) {
  const systemOSS = database
    .query("SELECT value_json FROM system_settings WHERE key = 'oss'")
    .get();
  if (systemOSS?.value_json) ossPayloads.push(String(systemOSS.value_json));
}
if (existingTables.has("user_oss_settings")) {
  for (const row of database
    .query("SELECT value_json FROM user_oss_settings")
    .all()) {
    if (row.value_json) ossPayloads.push(String(row.value_json));
  }
}

const secretChecks = [];
for (const payload of ossPayloads) {
  const parsed = JSON.parse(payload);
  for (const secret of collectSecrets(parsed))
    secretChecks.push(decryptSettingSecret(secret));
}

database.close();
key.fill(0);

const result = {
  integrityCheck,
  foreignKeyViolations,
  tableCount,
  counts,
  channelModelMigration,
  ossPayloadCount: ossPayloads.length,
  ossSecretCount: secretChecks.length,
  ossSecretsEncrypted:
    secretChecks.length > 0 && secretChecks.every((item) => item.encrypted),
  ossSecretsDecryptable:
    secretChecks.length > 0 && secretChecks.every((item) => item.decrypted),
};

console.log(JSON.stringify(result, null, 2));
if (
  integrityCheck !== "ok" ||
  foreignKeyViolations !== 0 ||
  !result.ossSecretsEncrypted ||
  !result.ossSecretsDecryptable
)
  process.exit(1);
