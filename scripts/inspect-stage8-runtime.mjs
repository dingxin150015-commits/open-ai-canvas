import { Database } from "bun:sqlite";

const [databasePath] = process.argv.slice(2);
if (!databasePath) {
  console.error("usage: bun scripts/inspect-stage8-runtime.mjs <database>");
  process.exit(2);
}

const database = new Database(databasePath, { readonly: true, strict: true });
const scalar = (sql, params = {}) => Number(database.query(sql).get(params)?.value || 0);
const latest = (table, columns) => database.query(`SELECT ${columns} FROM ${table} ORDER BY created_at DESC LIMIT 1`).get() || null;

const target = database
  .query(`
    SELECT id, model_key, support_status, enabled, price_configured, billing_mode,
           unit_price_microcredits, updated_at
    FROM channel_models
    WHERE model_key = 'wan3.0-video' AND deleted_at IS NULL
    LIMIT 1
  `)
  .get();
const targetTiers = target
  ? database
      .query(`
        SELECT id, selector_json, billing_mode, unit_price_microcredits,
               price_configured, enabled, price_version
        FROM channel_model_price_tiers
        WHERE channel_model_id = $modelID AND deleted_at IS NULL
        ORDER BY created_at
      `)
      .all({ modelID: target.id })
  : [];

const task = latest(
  "tasks",
  `id, type, status, operation, provider, model, attempts,
   json_extract(input_json, '$.config.size') AS requested_size,
   json_extract(input_json, '$.config.vquality') AS requested_resolution,
   json_extract(input_json, '$.config.videoSeconds') AS requested_seconds,
   json_extract(input_json, '$.config.videoGenerateAudio') AS requested_audio,
   json_extract(input_json, '$.config.videoWatermark') AS requested_watermark,
   json_extract(input_json, '$.metadata.seed') AS requested_seed,
   length(input_json) AS input_bytes, length(result_json) AS result_bytes,
   provider_request_id, error, created_at, updated_at`,
);
const billing = latest(
  "billing_orders",
  "id, task_id, model, capability, billing_mode, quantity, amount_microcredits, reserved_amount_microcredits, actual_amount_microcredits, refunded_amount_microcredits, status, provider_request_id, created_at, updated_at",
);
const resource = latest(
  "resources",
  "id, kind, status, provider, endpoint, mime_type, size, duration_ms, error, created_at, updated_at",
);
const asset = latest("assets", "id, kind, category, status, title, payload_json, created_at, updated_at");
const apiCall = latest(
  "api_call_logs",
  `id, task_id, model, capability, operation, request_kind, billable, status, status_code,
   provider_status, usage_available, video_seconds, provider_request_id, error_code,
   json_extract(request_body, '$.parameters.ratio') AS request_ratio,
   json_extract(request_body, '$.parameters.resolution') AS request_resolution,
   json_extract(request_body, '$.parameters.duration') AS request_duration,
   json_extract(request_body, '$.parameters.audio') AS request_audio,
   json_extract(request_body, '$.parameters.watermark') AS request_watermark,
   json_extract(request_body, '$.parameters.seed') AS request_seed,
   created_at`,
);
const credit = database.query("SELECT available_microcredits, reserved_microcredits, version, updated_at FROM credit_accounts LIMIT 1").get() || null;
const counts = {
  tasks: scalar("SELECT count(*) AS value FROM tasks"),
  resources: scalar("SELECT count(*) AS value FROM resources"),
  assets: scalar("SELECT count(*) AS value FROM assets"),
  billingOrders: scalar("SELECT count(*) AS value FROM billing_orders"),
  apiCallLogs: scalar("SELECT count(*) AS value FROM api_call_logs"),
};
database.close();

function safeTask(value) {
  if (!value) return null;
  const { error, provider_request_id: providerRequestID, ...rest } = value;
  return { ...rest, providerRequestID: String(providerRequestID || ""), hasError: Boolean(String(error || "").trim()) };
}

function safeResource(value) {
  if (!value) return null;
  const { endpoint, error, ...rest } = value;
  let endpointHost = "";
  try {
    endpointHost = new URL(String(endpoint || "")).host;
  } catch {}
  return { ...rest, endpointHost, hasError: Boolean(String(error || "").trim()) };
}

function safeAsset(value) {
  if (!value) return null;
  const { payload_json: payloadJSON, ...rest } = value;
  let storageKeyKind = "";
  try {
    const payload = JSON.parse(String(payloadJSON || "{}"));
    const storageKey = String(payload?.data?.storageKey || payload?.storageKey || "");
    storageKeyKind = storageKey.split(":", 1)[0] || "";
  } catch {}
  return { ...rest, payloadBytes: String(payloadJSON || "").length, storageKeyKind };
}

console.log(
  JSON.stringify(
    {
      counts,
      target: target
        ? {
            ...target,
            enabled: Boolean(target.enabled),
            price_configured: Boolean(target.price_configured),
            priceTiers: targetTiers.map((tier) => ({ ...tier, enabled: Boolean(tier.enabled), price_configured: Boolean(tier.price_configured) })),
          }
        : null,
      credit,
      latestTask: safeTask(task),
      latestBilling: billing,
      latestResource: safeResource(resource),
      latestAsset: safeAsset(asset),
      latestAPICall: apiCall,
    },
    null,
    2,
  ),
);
