import { Database } from "bun:sqlite";
import crypto from "node:crypto";

const [databasePath] = process.argv.slice(2);
if (!databasePath) {
    console.error("usage: bun scripts/inspect-stage7-catalog.mjs <database>");
    process.exit(2);
}

const database = new Database(databasePath, { readonly: true, strict: true });
const rows = database
    .query(`
        SELECT id, model_key, provider_model_key, display_name, capability, protocol,
               support_status, support_reason, catalog_source, catalog_version,
               supported_operations_json, documentation_paths_json,
               capability_config_json, enabled, price_configured,
               updated_at
        FROM channel_models
        WHERE deleted_at IS NULL
        ORDER BY model_key, id
    `)
    .all();
database.close();

function groupBy(field) {
    const result = {};
    for (const row of rows) {
        const key = String(row[field] || "(empty)");
        result[key] = (result[key] || 0) + 1;
    }
    return Object.fromEntries(Object.entries(result).sort(([left], [right]) => left.localeCompare(right)));
}

function parseList(raw) {
    try {
        const value = JSON.parse(String(raw || "[]"));
        return Array.isArray(value) ? value : [];
    } catch {
        return [];
    }
}

const readyVideoIDs = rows.filter((row) => row.support_status === "ready" && row.capability === "video").map((row) => String(row.model_key)).sort();
const readyImageIDs = rows.filter((row) => row.support_status === "ready" && row.capability === "image").map((row) => String(row.model_key)).sort();
const qwenImage30 = rows.find((row) => row.model_key === "qwen-image-3.0-pro");
const officialRows = rows.filter((row) => String(row.catalog_source).includes("bailian"));
const invalid = {
    readyMissingCapabilityConfig: rows.filter((row) => row.support_status === "ready" && !String(row.capability_config_json || "").trim()).length,
    readyMissingProtocol: rows.filter((row) => row.support_status === "ready" && !String(row.protocol || "").trim()).length,
    readyMissingCapability: rows.filter((row) => row.support_status === "ready" && !String(row.capability || "").trim()).length,
    readyMissingOperations: rows.filter((row) => row.support_status === "ready" && parseList(row.supported_operations_json).length === 0).length,
    readyMissingDocumentation: rows.filter((row) => row.support_status === "ready" && parseList(row.documentation_paths_json).length === 0).length,
    nonReadyEnabled: rows.filter((row) => row.support_status !== "ready" && Number(row.enabled) === 1).length,
    nonReadyPriced: rows.filter((row) => row.support_status !== "ready" && Number(row.price_configured) === 1).length,
};

const upstreamExpectedReason = "上游目录未提供经过项目验证的执行器合同";
const upstreamMismatch = {
    providerModelKey: rows.filter((row) => row.catalog_source === "upstream" && row.provider_model_key !== row.model_key).length,
    displayName: rows.filter((row) => row.catalog_source === "upstream" && row.display_name !== row.model_key).length,
    capability: rows.filter((row) => row.catalog_source === "upstream" && String(row.capability || "").trim() !== "").length,
    protocol: rows.filter((row) => row.catalog_source === "upstream" && String(row.protocol || "").trim() !== "").length,
    supportStatus: rows.filter((row) => row.catalog_source === "upstream" && row.support_status !== "planned").length,
    supportReason: rows.filter((row) => row.catalog_source === "upstream" && row.support_reason !== upstreamExpectedReason).length,
    catalogVersion: rows.filter((row) => row.catalog_source === "upstream" && row.catalog_version !== "upstream").length,
    operationsJSON: rows.filter((row) => row.catalog_source === "upstream" && String(row.supported_operations_json) !== "[]").length,
    documentationJSON: rows.filter((row) => row.catalog_source === "upstream" && String(row.documentation_paths_json) !== "[]").length,
};

const canonicalRows = rows.map((row) => ({
    id: row.id,
    modelKey: row.model_key,
    providerModelKey: row.provider_model_key,
    displayName: row.display_name,
    capability: row.capability,
    protocol: row.protocol,
    supportStatus: row.support_status,
    supportReason: row.support_reason,
    catalogSource: row.catalog_source,
    catalogVersion: row.catalog_version,
    operations: parseList(row.supported_operations_json),
    documentation: parseList(row.documentation_paths_json),
    capabilityConfig: row.capability_config_json,
    enabled: Number(row.enabled),
    priceConfigured: Number(row.price_configured),
    updatedAt: row.updated_at,
}));
const digest = crypto.createHash("sha256").update(JSON.stringify(canonicalRows)).digest("hex").toUpperCase();

const result = {
    total: rows.length,
    supportStatus: groupBy("support_status"),
    catalogSource: groupBy("catalog_source"),
    capability: groupBy("capability"),
    protocol: groupBy("protocol"),
    officialRows: officialRows.length,
    officialImage: officialRows.filter((row) => row.capability === "image").length,
    officialVideo: officialRows.filter((row) => row.capability === "video").length,
    configuredCapabilities: rows.filter((row) => String(row.capability_config_json || "").trim()).length,
    enabled: rows.filter((row) => Number(row.enabled) === 1).length,
    priced: rows.filter((row) => Number(row.price_configured) === 1).length,
    readyVideoIDs,
    readyImageIDs,
    qwenImage30: qwenImage30 ? {
        supportStatus: qwenImage30.support_status,
        supportReason: qwenImage30.support_reason,
        catalogSource: qwenImage30.catalog_source,
        catalogVersion: qwenImage30.catalog_version,
        capability: qwenImage30.capability,
        protocol: qwenImage30.protocol,
        capabilityConfigBytes: String(qwenImage30.capability_config_json || "").length,
        enabled: Number(qwenImage30.enabled),
        priceConfigured: Number(qwenImage30.price_configured),
    } : null,
    invalid,
    upstreamMismatch,
    digest,
};

console.log(JSON.stringify(result, null, 2));
if (Object.values(invalid).some((count) => count !== 0)) process.exit(1);
