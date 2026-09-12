import { Input, Select, Switch } from "antd";
import { bailianModelKind, type BailianGenerationOptions } from "@/lib/bailian-generation";

export function BailianGenerationOptionsPanel({
    model,
    value = {},
    imageCount = 0,
    showBatchControl = false,
    onChange,
}: {
    model: string;
    value?: BailianGenerationOptions;
    imageCount?: number;
    showBatchControl?: boolean;
    onChange: (value: BailianGenerationOptions) => void;
}) {
    const kind = bailianModelKind(model);
    if (!kind) return null;
    const qwen = kind === "qwen";
    const extend = value.prompt_extend !== false;
    const patch = (next: Partial<BailianGenerationOptions>) => onChange({ ...value, ...next });
    return (
        <details className="min-w-0 max-w-full space-y-2" data-canvas-no-zoom onMouseDown={(event) => event.stopPropagation()}>
            <summary className="cursor-pointer text-sm">高级生成参数</summary>
            <label className="block text-sm">
                随机种子（留空随机）
                <Input aria-label="随机种子" value={value.seed ?? ""} placeholder="0–2147483647" onChange={(event) => patch({ seed: event.target.value })} />
            </label>
            {qwen || kind.startsWith("wan2.7") ? (
                <>
                    <label className="block text-sm">
                        负向提示词
                        <Input.TextArea aria-label="负向提示词" maxLength={500} showCount value={value.negative_prompt ?? ""} onChange={(event) => patch({ negative_prompt: event.target.value })} />
                    </label>
                    <label className="flex items-center justify-between gap-2 text-sm">
                        提示词扩写
                        <Switch checked={extend} onChange={(checked) => patch({ prompt_extend: checked, ...(!checked ? { prompt_extend_mode: undefined, enable_thinking: false } : {}) })} />
                    </label>
                </>
            ) : null}
            {qwen ? (
                <>
                    <label className="flex items-center justify-between gap-2 text-sm">
                        请求方式
                        <Select
                            aria-label="Qwen 请求方式"
                            value={value.request_mode ?? "sync"}
                            options={[
                                { value: "sync", label: "原生同步" },
                                { value: "async", label: "原生异步" },
                                { value: "openai", label: "OpenAI 兼容" },
                            ]}
                            onChange={(request_mode) => patch({ request_mode })}
                        />
                    </label>
                    {showBatchControl ? (
                        <>
                            <label className="flex items-center justify-between gap-2 text-sm">
                                单请求多图
                                <Switch checked={value.nativeBatch === true} onChange={(nativeBatch) => patch({ nativeBatch })} />
                            </label>
                            <p className="text-xs">单请求多图作为一次任务，按当前模型的按次价格计费；关闭则每张独立任务。</p>
                        </>
                    ) : null}
                    {extend ? (
                        <>
                            <label className="flex items-center justify-between gap-2 text-sm">
                                扩写方式
                                <Select
                                    aria-label="扩写方式"
                                    value={value.prompt_extend_mode ?? "direct"}
                                    options={[
                                        { value: "direct", label: "Direct" },
                                        { value: "agent", label: "Agent（仅文生图）", disabled: imageCount > 0 },
                                    ]}
                                    onChange={(mode) => patch({ prompt_extend_mode: mode })}
                                />
                            </label>
                            <label className="flex items-center justify-between gap-2 text-sm">
                                思考模式
                                <Switch checked={value.enable_thinking !== false} onChange={(checked) => patch({ enable_thinking: checked })} />
                            </label>
                        </>
                    ) : null}
                    {imageCount > 0 && value.prompt_extend_mode === "agent" ? <p role="alert">参考图编辑不支持 Agent 扩写，请选择 Direct。</p> : null}
                    <label className="flex items-center justify-between gap-2 text-sm">
                        图片水印
                        <Switch checked={value.watermark === true} onChange={(checked) => patch({ watermark: checked })} />
                    </label>
                </>
            ) : null}
        </details>
    );
}
