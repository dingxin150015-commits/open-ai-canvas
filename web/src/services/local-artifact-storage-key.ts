export type LocalArtifactStore = "image" | "media";

const IMAGE_STORAGE_PREFIXES = ["image:", "generation-image:"];
const MEDIA_STORAGE_PREFIXES = ["video:", "audio:", "file:", "video-reference:", "audio-reference:", "generation-video:", "generation-audio:"];

export function localArtifactStoreForStorageKey(value: string): LocalArtifactStore | null {
    const key = value.trim();
    if (IMAGE_STORAGE_PREFIXES.some((prefix) => key.startsWith(prefix))) return "image";
    if (MEDIA_STORAGE_PREFIXES.some((prefix) => key.startsWith(prefix))) return "media";
    return null;
}
