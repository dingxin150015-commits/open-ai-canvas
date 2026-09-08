// Exact dependencies of the standalone Director lab. Never enabled in normal
// development/production and never used as a blanket network-error allowlist.
export function directorE2EFixture(url, method = "GET") {
    if (method !== "GET") return undefined;
    const pathname = String(url || "").split("?", 1)[0];
    if (pathname === "/api/public/appearance")
        return {
            contentType: "application/json",
            body: JSON.stringify({ code: 0, msg: "ok", data: { brandName: "影策", activeSkin: "default" } }),
        };
    if (pathname === "/favicon.ico")
        return {
            contentType: "image/svg+xml",
            body: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path fill="black" d="M2 2h12v12H2z"/></svg>',
        };
    return undefined;
}
