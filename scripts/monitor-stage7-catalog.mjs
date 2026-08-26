import { Database } from "bun:sqlite";

const [databasePath, rawTimeoutSeconds = "1800"] = process.argv.slice(2);
if (!databasePath) {
    console.error("usage: bun scripts/monitor-stage7-catalog.mjs <database> [timeout-seconds]");
    process.exit(2);
}

const timeoutSeconds = Math.max(30, Math.min(3600, Number(rawTimeoutSeconds) || 1800));
const deadline = Date.now() + timeoutSeconds * 1000;
let previous = "";

function readSnapshot() {
    const database = new Database(databasePath, { readonly: true, strict: true });
    try {
        const total = Number(database.query("SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL").get().value);
        const group = (column) => Object.fromEntries(database.query(`SELECT ${column} AS name, count(*) AS value FROM channel_models WHERE deleted_at IS NULL GROUP BY ${column} ORDER BY ${column}`).all().map((row) => [String(row.name || "(empty)"), Number(row.value)]));
        const scalar = (sql) => Number(database.query(sql).get().value);
        return {
            observedAt: new Date().toISOString(),
            total,
            supportStatus: group("support_status"),
            catalogSource: group("catalog_source"),
            capability: group("capability"),
            readyVideo: scalar("SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND support_status = 'ready' AND capability = 'video'"),
            plannedVideo: scalar("SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND support_status = 'planned' AND capability = 'video'"),
            plannedImage: scalar("SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND support_status = 'planned' AND capability = 'image'"),
            configuredCapabilities: scalar("SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND capability_config_json <> ''"),
            enabled: scalar("SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND enabled = 1"),
            priced: scalar("SELECT count(*) AS value FROM channel_models WHERE deleted_at IS NULL AND price_configured = 1"),
        };
    } finally {
        database.close();
    }
}

while (Date.now() < deadline) {
    const snapshot = readSnapshot();
    const comparable = JSON.stringify({ ...snapshot, observedAt: undefined });
    if (comparable !== previous) {
        console.log(JSON.stringify(snapshot));
        previous = comparable;
    }
    await Bun.sleep(1000);
}
