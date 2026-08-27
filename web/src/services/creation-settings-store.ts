import { localForageStorageForScope } from "@/lib/localforage-storage";
import { getActiveUserScope } from "@/lib/user-scope";

export const CREATION_SETTINGS_DRAFTS_KEY = "open_ai_canvas:creation_settings_drafts";

export type CreationSettingsDraftMode = "image" | "video";
export type CreationSettingsDraftIdentity = {
    mode: CreationSettingsDraftMode;
    model: string;
    operation: string;
};
export type CreationSettingsDraft = {
    ratio: string;
    seconds?: string;
    quality?: string;
    videoQuality?: string;
    count?: string;
    videoGenerateAudio?: boolean;
    videoWatermark?: boolean;
};

type CreationSettingsDraftDocument = {
    version: 1;
    drafts: Record<string, CreationSettingsDraft>;
};

const writeTails = new Map<string, Promise<void>>();

export function creationSettingsDraftKey(identity: CreationSettingsDraftIdentity) {
    return [identity.mode, identity.model.trim(), identity.operation.trim()].map((value) => encodeURIComponent(value)).join("|");
}

export async function loadCreationSettingsDraft(identity: CreationSettingsDraftIdentity, scope?: string) {
    const storage = localForageStorageForScope(scope);
    const document = parseCreationSettingsDraftDocument(await storage.getItem(CREATION_SETTINGS_DRAFTS_KEY));
    return document.drafts[creationSettingsDraftKey(identity)];
}

export function saveCreationSettingsDraft(identity: CreationSettingsDraftIdentity, draft: CreationSettingsDraft, scope?: string) {
    const resolvedScope = scope ?? getActiveUserScope();
    const previous = writeTails.get(resolvedScope) ?? Promise.resolve();
    const operation = previous
        .catch(() => undefined)
        .then(async () => {
            const storage = localForageStorageForScope(resolvedScope);
            const document = parseCreationSettingsDraftDocument(await storage.getItem(CREATION_SETTINGS_DRAFTS_KEY));
            document.drafts[creationSettingsDraftKey(identity)] = normalizeCreationSettingsDraft(draft);
            await storage.setItem(CREATION_SETTINGS_DRAFTS_KEY, JSON.stringify(document));
        });
    writeTails.set(resolvedScope, operation);
    void operation
        .finally(() => {
            if (writeTails.get(resolvedScope) === operation) writeTails.delete(resolvedScope);
        })
        .catch(() => undefined);
    return operation;
}

export function parseCreationSettingsDraftDocument(value: string | null): CreationSettingsDraftDocument {
    if (!value) return emptyDocument();
    try {
        const parsed = JSON.parse(value) as Partial<CreationSettingsDraftDocument>;
        if (parsed.version !== 1 || !parsed.drafts || typeof parsed.drafts !== "object" || Array.isArray(parsed.drafts)) return emptyDocument();
        const drafts: Record<string, CreationSettingsDraft> = {};
        for (const [key, draft] of Object.entries(parsed.drafts)) {
            if (!key || key.length > 2_048 || !draft || typeof draft !== "object" || Array.isArray(draft)) continue;
            const normalized = normalizeCreationSettingsDraft(draft as Partial<CreationSettingsDraft>);
            if (normalized.ratio) drafts[key] = normalized;
        }
        return { version: 1, drafts };
    } catch {
        return emptyDocument();
    }
}

function normalizeCreationSettingsDraft(draft: Partial<CreationSettingsDraft>): CreationSettingsDraft {
    return {
        ratio: boundedString(draft.ratio),
        ...(boundedString(draft.seconds) ? { seconds: boundedString(draft.seconds) } : {}),
        ...(boundedString(draft.quality) ? { quality: boundedString(draft.quality) } : {}),
        ...(boundedString(draft.videoQuality) ? { videoQuality: boundedString(draft.videoQuality) } : {}),
        ...(boundedString(draft.count) ? { count: boundedString(draft.count) } : {}),
        ...(typeof draft.videoGenerateAudio === "boolean" ? { videoGenerateAudio: draft.videoGenerateAudio } : {}),
        ...(typeof draft.videoWatermark === "boolean" ? { videoWatermark: draft.videoWatermark } : {}),
    };
}

function boundedString(value: unknown) {
    return typeof value === "string" ? value.trim().slice(0, 128) : "";
}

function emptyDocument(): CreationSettingsDraftDocument {
    return { version: 1, drafts: {} };
}
