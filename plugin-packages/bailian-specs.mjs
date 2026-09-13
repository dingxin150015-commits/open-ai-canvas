// Model-specific contracts are data. The host evaluates only generic operators.
const ref = (x) => ({ $ref: x });
const v = (x) => ref(`vars.${x}`);
const req = (x) => ref(`request.${x}`);
const op = (name, ...args) => ({ [name]: args });
const eq = (a, b) => op("$eq", a, b),
  ne = (a, b) => op("$ne", a, b),
  gt = (a, b) => op("$gt", a, b),
  le = (a, b) => op("$lte", a, b),
  ge = (a, b) => op("$gte", a, b);
const and = (...x) => op("$and", ...x),
  or = (...x) => op("$or", ...x),
  not = (x) => ({ $not: x });
const len = (x) => ({ $len: x }),
  omit = (x) => ({ $omitEmpty: x }),
  co = (...x) => op("$coalesce", ...x),
  inside = (x, values) => op("$in", x, values);
const iff = (condition, then, otherwise = null) => ({
  $if: { condition, then, else: otherwise },
});
const choice = (cases, otherwise = null) => ({
  $switch: {
    cases: cases.map(([when, then]) => ({ when, then })),
    default: otherwise,
  },
});
const map = (from, as, body) => ({ $map: { from, as, in: body } });
const filter = (from, as, where) => ({ $filter: { from, as, where } });
const cat = (...xs) => op("$concatArrays", ...xs);
const lower = (x) => ({ $lower: { $trim: x } }),
  upper = (x) => ({ $upper: { $trim: x } });
const integer = (value, def, min, max, allow = []) => ({
  $strictInt: { value, default: def, min, max, allow },
});
const boolean = (value, def) => ({ $strictBool: { value, default: def } });
const variable = (name, value) => ({ name, value });
const rule = (assert, message) => ({ assert, message });
const implies = (when, then) => or(not(when), then);
const first = (x) => ({ $first: x });
const media = (items, type) => map(items, "m", { type, url: ref("m.value") });
const policy = (items, params) => ({ $validateMedia: { items, ...params } });
const commonPath = { basePathSuffixes: ["/compatible-mode/v1", "/api/v1"] };
const create = (path, body, async = false) => ({
  ...commonPath,
  method: "POST",
  path,
  contentType: "application/json",
  body,
  ...(async ? { headers: { "X-DashScope-Async": "enable" } } : {}),
});
const auth = { type: "bearer", field: "apiKey" };
const MB = 1024 * 1024;
const imageMimes = ["image/jpeg", "image/jpg", "image/png", "image/webp"];
const audioMimes = ["audio/wav", "audio/x-wav", "audio/mpeg", "audio/mp3"];
const videoMimes = ["video/mp4", "video/quicktime"];
const ratio5 = ["16:9", "9:16", "1:1", "4:3", "3:4"];
const ratio9 = [...ratio5, "4:5", "5:4", "9:21", "21:9"];
const varsBase = (id) => [
  variable("model", { $replace: [lower(req("model")), "models/", ""] }),
  variable("meta", {
    $merge: [req("extra.generationMetadata"), req(`providerOptions.${id}`)],
  }),
  variable("raw", req("extra.rawOptions")),
  variable("prompt", { $trim: req("prompt") }),
  variable(
    "images",
    id === "dashscope-image"
      ? map(req("images"), "m", {
          $merge: [ref("m"), { value: co(ref("m.dataUrl"), ref("m.url")) }],
        })
      : req("images"),
  ),
  variable("videos", req("videos")),
  variable("audios", req("audios")),
];
const params = [
  ["model", "string", true, "model", "精确模型 ID"],
  ["prompt", "string", false, "input", "模型输入提示词"],
  ["images", "media[]", false, "input", "图片素材"],
  ["videos", "media[]", false, "input", "视频素材"],
  ["audios", "media[]", false, "input", "音频素材"],
  ["duration", "integer", false, "parameters.duration", "时长"],
  ["aspectRatio", "string", false, "parameters", "画幅或尺寸"],
];
const referenceIntegrity = rule(
  and(
    ...["image", "video", "audio"].map((kind) =>
      eq(
        len(v(`${kind}s`)),
        co(req(`extra.referenceCounts.${kind}`), len(v(`${kind}s`))),
      ),
    ),
  ),
  "参考素材缺少有效地址，不能静默忽略",
);
const qwenImageModels = ["qwen-image-3.0-pro"];
export const qwenPixelPolicy = {
  presets: {
    "1:1": "1024x1024",
    "2:3": "1024x1536",
    "3:2": "1536x1024",
    "3:4": "768x1024",
    "4:3": "1024x768",
    "9:16": "720x1280",
    "16:9": "1280x720",
  },
  targetPixels: 1048576,
  ratioMaxSide: 2048,
  minPixels: 262144,
  maxPixels: 4194304,
  minRatio: 0.125,
  maxRatio: 8,
  separator: "*",
};
const imageSpec = {
  apiVersion: "yingce.plugin/v3",
  version: "3.1.0",
  id: "dashscope-image",
  providerId: "dashscope-image",
  name: "DashScope Qwen Image",
  vendor: "Alibaba Cloud",
  capability: "image",
  baseUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1",
  auth,
  params,
  variables: [
    ...varsBase("dashscope-image").filter((item) => item.name !== "prompt"),
    variable("prompt", {
      $trim: iff(
        ne(req("instructions"), ""),
        op("$concat", req("instructions"), "\n\n", req("prompt")),
        req("prompt"),
      ),
    }),
    variable("count", integer(v("raw.count"), 1, 1, 6)),
    variable("size", {
      $pixelSize: { value: req("aspectRatio"), ...qwenPixelPolicy },
    }),
    variable("extend", boolean(v("meta.prompt_extend"), true)),
    variable("thinking", boolean(v("meta.enable_thinking"), v("extend"))),
    variable("extendMode", co(lower(v("meta.prompt_extend_mode")), "direct")),
    variable("watermark", boolean(v("meta.watermark"), false)),
    variable("seed", integer(v("meta.seed"), null, 0, 2147483647)),
    variable("negative", co(v("meta.negative_prompt"), "")),
  ],
  validations: [
    referenceIntegrity,
    rule(
      inside(v("model"), qwenImageModels),
      "该图片插件仅支持已验证的 Qwen Image 3.0 Pro",
    ),
    rule(gt(len(v("prompt")), 0), "必须提供提示词"),
    rule(
      and(
        le(len(v("images")), 3),
        eq(len(v("videos")), 0),
        eq(len(v("audios")), 0),
      ),
      "Qwen 最多支持 3 张参考图且不支持音视频输入",
    ),
    rule(
      eq(len(filter(v("images"), "m", eq(ref("m.role"), "mask"))), 0),
      "Qwen 不支持蒙版",
    ),
    rule(inside(lower(req("quality")), ["", "auto"]), "Qwen 不支持质量参数"),
    rule(
      not(boolean(v("raw.transparentBackground"), false)),
      "Qwen 不支持透明背景",
    ),
    rule(
      inside(v("extendMode"), ["direct", "agent"]),
      "改写模式必须为 direct 或 agent",
    ),
    rule(
      implies(
        not(v("extend")),
        and(eq(v("meta.prompt_extend_mode"), null), not(v("thinking"))),
      ),
      "关闭改写时不能设置改写模式或开启思考",
    ),
    rule(
      implies(gt(len(v("images")), 0), ne(v("extendMode"), "agent")),
      "图片编辑不支持 agent 改写",
    ),
    rule(le(len(v("negative")), 500), "负面提示词不能超过 500 字符"),
    rule(
      policy(v("images"), {
        maxBytes: 10 * MB,
        mimes: [
          ...imageMimes,
          "image/bmp",
          "image/x-ms-bmp",
          "image/tiff",
          "image/gif",
        ],
        allowData: true,
      }),
      "参考图片不符合限制",
    ),
  ],
  create: create("/api/v1/services/aigc/multimodal-generation/generation", {
    model: v("model"),
    input: {
      messages: [
        {
          role: "user",
          content: cat(
            map(v("images"), "m", {
              image: co(ref("m.dataUrl"), ref("m.url")),
            }),
            [{ text: v("prompt") }],
          ),
        },
      ],
    },
    parameters: {
      n: v("count"),
      size: v("size"),
      prompt_extend: v("extend"),
      enable_thinking: v("thinking"),
      prompt_extend_mode: iff(v("extend"), v("extendMode")),
      watermark: v("watermark"),
      negative_prompt: omit(v("negative")),
      seed: v("seed"),
    },
  }),
  response: {
    strict: true,
    status: "succeeded",
    images: {
      $flatten: map(
        reqUnusedResponseChoices(),
        "c",
        map(ref("c.message.content"), "m", ref("m.image")),
      ),
    },
    usage: ref("response.usage"),
    errorPaths: ["code", "output.code"],
    messagePaths: ["message", "output.message"],
    resultEphemeral: true,
    resultAnonymous: true,
  },
  notes:
    "仅支持 qwen-image-3.0-pro。同步文生图与 1–3 图编辑；不支持蒙版或质量参数。auto 省略尺寸；内部比例和像素在插件中规范化。保留账号区域主机和路径前缀。宿主负责 SSRF 防护、凭证、计费及持久化。",
};
function reqUnusedResponseChoices() {
  return ref("response.output.choices");
}
const t2v = ["wan2.7-t2v", "wan2.7-t2v-2026-06-12", "wan2.7-t2v-2026-04-25"];
const i2v = ["wan2.7-i2v", "wan2.7-i2v-2026-04-25"];
const r2v = ["wan2.7-r2v", "wan2.7-r2v-2026-06-12"];
const horse = [
  "happyhorse-1.1-t2v",
  "happyhorse-1.1-i2v",
  "happyhorse-1.1-r2v",
];
const is = (kind) => eq(v("kind"), kind),
  family = (f) => eq(v("family"), f);
const imageCount = len(v("images")),
  videoCount = len(v("videos")),
  audioCount = len(v("audios"));
const visual = op("$add", imageCount, videoCount),
  allMedia = op("$add", imageCount, videoCount, audioCount);
const frameMode = or(
  is("i2v"),
  and(family("wan30"), eq(v("operation"), "image_to_video")),
);
const audioDuration = {
  $sum: map(v("audios"), "m", ref("m.metadata.durationMs")),
};
const videoDuration = {
  $sum: map(v("videos"), "m", ref("m.metadata.durationMs")),
};
const orderedImages = iff(
  and(frameMode, or(ne(v("first"), ""), ne(v("last"), ""))),
  { $selectByIds: { from: v("images"), ids: [v("first"), v("last")] } },
  v("images"),
);
const videoSpec = {
  apiVersion: "yingce.plugin/v3",
  version: "3.1.0",
  id: "dashscope-video",
  providerId: "dashscope-video",
  name: "DashScope Wan / HappyHorse Video",
  vendor: "Alibaba Cloud",
  capability: "video",
  baseUrl: "https://dashscope.aliyuncs.com/compatible-mode/v1",
  auth,
  params,
  requiresPublicMediaUrls: true,
  pollIntervalMs: 15000,
  variables: [
    ...varsBase("dashscope-video"),
    variable(
      "family",
      choice(
        [
          [inside(v("model"), [...t2v, ...i2v, ...r2v]), "wan27"],
          [inside(v("model"), horse), "horse"],
          [eq(v("model"), "wan3.0-video"), "wan30"],
        ],
        "unknown",
      ),
    ),
    variable(
      "kind",
      choice(
        [
          [inside(v("model"), [...t2v, horse[0]]), "t2v"],
          [inside(v("model"), [...i2v, horse[1]]), "i2v"],
          [inside(v("model"), [...r2v, horse[2]]), "r2v"],
        ],
        "all",
      ),
    ),
    variable("first", { $trim: v("raw.firstFrameID") }),
    variable("last", { $trim: v("raw.lastFrameID") }),
    variable(
      "operation",
      co(
        lower(req("operation")),
        choice(
          [
            [or(ne(v("first"), ""), ne(v("last"), "")), "image_to_video"],
            [
              and(
                family("wan30"),
                or(gt(videoCount, 0), gt(audioCount, 0), gt(imageCount, 2)),
              ),
              "reference_to_video",
            ],
            [and(family("wan30"), gt(imageCount, 0)), "image_to_video"],
            [is("r2v"), "reference_to_video"],
            [is("i2v"), "image_to_video"],
          ],
          "text_to_video",
        ),
      ),
    ),
    variable(
      "duration",
      integer(
        co(v("raw.duration"), req("extra.videoSeconds")),
        5,
        iff(family("horse"), 3, 2),
        iff(family("wan30"), 30, 15),
        iff(family("wan30"), [-1], []),
      ),
    ),
    variable(
      "resolution",
      choice(
        [
          [
            inside(upper(req("resolution")), ["480", "720", "1080"]),
            op("$concat", upper(req("resolution")), "P"),
          ],
        ],
        co(upper(req("resolution")), "1080P"),
      ),
    ),
    variable(
      "ratio",
      co(lower(req("aspectRatio")), iff(family("wan30"), "adaptive", "16:9")),
    ),
    variable("watermark", boolean(v("raw.watermark"), family("horse"))),
    variable("audio", boolean(v("raw.generateAudio"), true)),
    variable("extend", boolean(v("meta.prompt_extend"), true)),
    variable("seed", integer(v("meta.seed"), null, 0, 2147483647)),
    variable("negative", co(v("meta.negative_prompt"), "")),
    variable("orderedImages", orderedImages),
  ],
  validations: [
    referenceIntegrity,
    rule(ne(v("family"), "unknown"), "模型尚未由此视频插件支持"),
    rule(
      inside(
        v("resolution"),
        iff(family("wan27"), ["720P", "1080P"], ["480P", "720P", "1080P"]),
      ),
      "分辨率不受模型支持",
    ),
    rule(
      implies(
        not(or(is("i2v"), and(family("wan27"), is("r2v"), ne(v("first"), "")))),
        inside(
          v("ratio"),
          iff(
            family("horse"),
            ratio9,
            iff(family("wan30"), ["adaptive", ...ratio5], ratio5),
          ),
        ),
      ),
      "比例不受模型支持",
    ),
    rule(
      le(len(v("prompt")), iff(family("wan30"), 20000, 5000)),
      "提示词超过模型限制",
    ),
    rule(
      implies(
        family("horse"),
        le({ $weightedLength: { value: v("prompt"), hanWeight: 2 } }, 5000),
      ),
      "HappyHorse 提示词最多 2500 中文字符或 5000 非中文字符",
    ),
    rule(
      implies(
        and(not(is("i2v")), not(family("wan30"))),
        gt(len(v("prompt")), 0),
      ),
      "必须提供提示词",
    ),
    rule(
      implies(family("wan30"), gt(op("$add", len(v("prompt")), allMedia), 0)),
      "提示词或参考素材至少提供一项",
    ),
    rule(le(len(v("negative")), 500), "负向提示词超过 500 字符"),
    rule(
      implies(
        is("t2v"),
        and(
          eq(imageCount, 0),
          eq(videoCount, 0),
          le(audioCount, iff(family("horse"), 0, 1)),
          inside(
            v("operation"),
            iff(
              family("horse"),
              ["text_to_video"],
              ["text_to_video", "audio_to_video"],
            ),
          ),
        ),
      ),
      "文生视频输入或操作不受支持",
    ),
    rule(
      implies(
        and(family("horse"), is("i2v")),
        and(
          eq(imageCount, 1),
          eq(videoCount, 0),
          eq(audioCount, 0),
          eq(v("last"), ""),
          eq(v("operation"), "image_to_video"),
        ),
      ),
      "HappyHorse 图生视频仅支持一张首帧",
    ),
    rule(
      implies(
        and(family("horse"), is("r2v")),
        and(
          ge(imageCount, 1),
          le(imageCount, 9),
          eq(videoCount, 0),
          eq(audioCount, 0),
          eq(v("first"), ""),
          eq(v("last"), ""),
          inside(v("operation"), ["reference_to_video", "image_to_video"]),
        ),
      ),
      "HappyHorse 参考生视频仅支持 1–9 张普通参考图",
    ),
    rule(
      implies(
        and(family("wan27"), is("i2v")),
        and(
          le(imageCount, 2),
          le(videoCount, 1),
          le(audioCount, 1),
          ge(visual, 1),
          inside(v("operation"), [
            "image_to_video",
            "extend",
            "reference_to_video",
            "audio_to_video",
          ]),
        ),
      ),
      "Wan 2.7 图生视频输入不受支持",
    ),
    rule(
      implies(
        and(family("wan27"), is("i2v"), gt(videoCount, 0)),
        and(
          eq(audioCount, 0),
          le(imageCount, 1),
          eq(v("first"), ""),
          gt(op("$multiply", v("duration"), 1000), videoDuration),
        ),
      ),
      "视频续写仅允许续写视频及可选尾帧，输出须长于输入",
    ),
    rule(
      implies(
        and(family("wan27"), is("r2v")),
        and(
          ge(visual, 1),
          le(visual, 5),
          le(audioCount, 1),
          eq(v("last"), ""),
          inside(v("operation"), [
            "reference_to_video",
            "image_to_video",
            "audio_to_video",
          ]),
        ),
      ),
      "Wan 2.7 参考生视频需要 1–5 个视觉素材和最多 1 个音频",
    ),
    rule(
      implies(
        and(family("wan27"), is("r2v"), gt(videoCount, 0)),
        le(v("duration"), 10),
      ),
      "含参考视频时输出最多 10 秒",
    ),
    rule(
      implies(and(ne(v("last"), ""), eq(videoCount, 0)), ne(v("first"), "")),
      "配置尾帧必须同时配置首帧",
    ),
    rule(
      implies(
        and(family("wan27"), is("r2v"), ne(v("first"), "")),
        eq(len(filter(v("images"), "m", eq(ref("m.id"), v("first")))), 1),
      ),
      "指定首帧未包含在参考图中",
    ),
    rule(
      implies(
        family("wan30"),
        and(
          le(imageCount, 9),
          le(videoCount, 3),
          le(audioCount, 1),
          inside(v("operation"), [
            "text_to_video",
            "image_to_video",
            "reference_to_video",
            "audio_to_video",
          ]),
        ),
      ),
      "Wan 3.0 参考数量或操作超出限制",
    ),
    rule(
      implies(
        and(family("wan30"), or(ne(v("first"), ""), ne(v("last"), ""))),
        and(
          eq(v("operation"), "image_to_video"),
          eq(videoCount, 0),
          eq(audioCount, 0),
        ),
      ),
      "Wan 3.0 帧图不能与参考音视频混用",
    ),
    rule(
      implies(
        and(family("wan30"), eq(v("operation"), "image_to_video")),
        and(
          ge(imageCount, 1),
          le(imageCount, 2),
          eq(videoCount, 0),
          eq(audioCount, 0),
        ),
      ),
      "图生视频仅支持 1–2 张帧图",
    ),
    rule(
      implies(
        and(family("wan30"), eq(v("operation"), "text_to_video")),
        eq(allMedia, 0),
      ),
      "文生视频不能携带素材",
    ),
    rule(
      implies(
        and(
          family("wan30"),
          inside(v("operation"), ["reference_to_video", "audio_to_video"]),
        ),
        gt(allMedia, 0),
      ),
      "参考模式至少需要一个素材",
    ),
    rule(
      implies(
        and(family("wan30"), gt(v("duration"), 0)),
        le(
          op("$add", videoDuration, op("$multiply", v("duration"), 1000)),
          30000,
        ),
      ),
      "参考视频与输出时长之和不能超过 30 秒",
    ),
    rule(
      policy(v("images"), {
        maxBytes: iff(family("wan30"), 30 * MB, 20 * MB),
        mimes: iff(family("horse"), imageMimes, [...imageMimes, "image/bmp"]),
        minSide: iff(
          family("horse"),
          iff(is("i2v"), 300, 400),
          iff(family("wan27"), 240, 0),
        ),
        maxSide: iff(family("wan27"), 8000, 0),
        maxAspect: iff(
          and(family("horse"), is("i2v")),
          2.5,
          iff(family("wan27"), 8, 0),
        ),
        allowData: not(family("wan27")),
        allowedOrigins: iff(family("wan30"), ["oss://dashscope-instant"], []),
      }),
      "图片不符合规格",
    ),
    rule(
      policy(v("videos"), {
        maxBytes: iff(family("wan30"), 50 * MB, 100 * MB),
        mimes: videoMimes,
        minDurationMs: iff(is("i2v"), 2000, 1000),
        maxDurationMs: iff(is("i2v"), 10000, 30000),
        minSide: iff(family("wan27"), 240, 0),
        maxSide: iff(family("wan27"), 4096, 0),
        maxAspect: iff(family("wan27"), 8, 0),
        allowData: family("wan30"),
        allowedOrigins: iff(family("wan30"), ["oss://dashscope-instant"], []),
      }),
      "视频不符合规格",
    ),
    rule(
      policy(v("audios"), {
        maxBytes: 15 * MB,
        mimes: audioMimes,
        minDurationMs: iff(is("r2v"), 1000, 2000),
        maxDurationMs: iff(and(family("wan27"), is("r2v")), 10000, 30000),
        allowData: family("wan30"),
        allowedOrigins: iff(family("wan30"), ["oss://dashscope-instant"], []),
      }),
      "音频不符合规格",
    ),
  ],
  create: create(
    "/api/v1/services/aigc/video-generation/video-synthesis",
    {
      model: v("model"),
      input: {
        prompt: omit(v("prompt")),
        negative_prompt: iff(family("wan27"), omit(v("negative"))),
        audio_url: iff(
          and(family("wan27"), is("t2v")),
          ref("vars.audios.0.value"),
        ),
        reference_voice: iff(
          and(family("wan27"), is("r2v")),
          ref("vars.audios.0.value"),
        ),
        media: omit(
          cat(
            iff(
              and(family("wan27"), is("i2v"), gt(videoCount, 0)),
              media(v("videos"), "first_clip"),
              [],
            ),
            map(v("orderedImages"), "m", {
              type: choice(
                [
                  [
                    and(family("wan27"), is("i2v"), gt(videoCount, 0)),
                    "last_frame",
                  ],
                  [
                    frameMode,
                    iff(eq(ref("mIndex"), 0), "first_frame", "last_frame"),
                  ],
                  [
                    and(
                      family("wan27"),
                      is("r2v"),
                      eq(ref("m.id"), v("first")),
                      ne(v("first"), ""),
                    ),
                    "first_frame",
                  ],
                ],
                "reference_image",
              ),
              url: ref("m.value"),
            }),
            iff(
              and(family("wan27"), is("i2v")),
              [],
              media(v("videos"), "reference_video"),
            ),
            iff(
              family("wan30"),
              media(v("audios"), "reference_audio"),
              iff(
                and(family("wan27"), is("i2v")),
                media(v("audios"), "driving_audio"),
                [],
              ),
            ),
          ),
        ),
      },
      parameters: {
        duration: v("duration"),
        resolution: v("resolution"),
        ratio: iff(
          or(is("i2v"), and(family("wan27"), is("r2v"), ne(v("first"), ""))),
          null,
          v("ratio"),
        ),
        watermark: v("watermark"),
        seed: v("seed"),
        prompt_extend: iff(family("wan27"), v("extend")),
        audio: iff(family("wan30"), v("audio")),
      },
    },
    true,
  ),
  poll: { ...commonPath, method: "GET", path: "/api/v1/tasks/{{taskId}}" },
  response: {
    strict: true,
    taskId: co(ref("response.output.task_id"), ref("taskId")),
    status: choice(
      [[eq(upper(ref("response.output.task_status")), "UNKNOWN"), "failed"]],
      ref("response.output.task_status"),
    ),
    videos: ref("response.output.video_url"),
    usage: ref("response.usage"),
    errorPaths: ["code", "output.code"],
    messagePaths: ["output.message", "message"],
    resultEphemeral: true,
    resultAnonymous: true,
  },
  notes:
    "精确支持 Wan 2.7 七个主线/快照 ID、HappyHorse 1.1 三个 ID、wan3.0-video。未登记模型拒绝。保留首尾帧、续写、参考媒体、驱动音频、智能时长及各自输入限制。异步轮询间隔 15 秒；UNKNOWN 视为终态失败；结果 URL 临时有效，宿主立即物化。",
};
// All Qwen transports reuse the native semantic contract; only wire encoding,
// endpoint and asynchronous lifecycle differ.
const nativeImageBody = imageSpec.create.body;
imageSpec.variables.push(variable("requestMode", co(v("meta.request_mode"), "sync")));
imageSpec.validations.push(rule(inside(v("requestMode"), ["sync", "async", "openai"]), "未知图片请求方式"));
imageSpec.create.pathTemplate = choice([
  [eq(v("requestMode"), "async"), "/api/v1/services/aigc/image-generation/generation"],
  [eq(v("requestMode"), "openai"), "/compatible-mode/v1/images/generations"],
], imageSpec.create.path);
imageSpec.create.headers = { "X-DashScope-Async": iff(eq(v("requestMode"), "async"), "enable") };
imageSpec.create.body = iff(eq(v("requestMode"), "openai"), {
  model: v("model"), prompt: v("prompt"), image: omit(map(v("images"), "m", co(ref("m.dataUrl"), ref("m.url")))),
  ...nativeImageBody.parameters, size: op("$replace", co(v("size"), "auto"), "*", "x"),
}, nativeImageBody);
imageSpec.poll = { ...commonPath, method: "GET", path: "/api/v1/tasks/{{taskId}}" };
imageSpec.pollIntervalMs = 10000;
imageSpec.response.taskId = co(ref("response.output.task_id"), ref("taskId"));
imageSpec.response.status = choice([[eq(upper(ref("response.output.task_status")), "UNKNOWN"), "failed"]], co(ref("response.output.task_status"), "succeeded"));
imageSpec.response.images = cat(imageSpec.response.images, map(ref("response.output.results"), "m", ref("m.url")), map(ref("response.data"), "m", ref("m.url")));
imageSpec.response.errorPaths.push("error.code", "error.type");
imageSpec.response.messagePaths.push("error.message");
imageSpec.notes += " 支持 request_mode=sync/async/openai，默认 sync。异步使用专属 image-generation 端点与 10 秒轮询；兼容接口使用字母 x 和 image 数组。";
imageSpec.modelContract = { mode: "exact", models: qwenImageModels };
videoSpec.modelContract = { mode: "exact", models: [...t2v, ...i2v, ...r2v, ...horse, "wan3.0-video"] };
export default [imageSpec, videoSpec];
