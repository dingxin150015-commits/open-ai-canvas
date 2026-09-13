# DashScope Qwen Image 接口字段

## 协议身份

- 插件 ID：`dashscope-image`。
- Provider ID：`dashscope-image`。
- 能力：`image`。
- 默认 Base URL：`https://dashscope.aliyuncs.com/compatible-mode/v1`。
- 鉴权驱动：`bearer`。
- 创建：`POST /api/v1/services/aigc/multimodal-generation/generation`。
- 查询：`GET /api/v1/tasks/{{taskId}}`。

## 配置字段

| 字段 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `apiKey` | secret | 是 | API Key |

## 统一字段映射

| 统一字段 | 类型 | 必填 | 上游映射 | 说明 |
| --- | --- | --- | --- | --- |
| `model` | string | 是 | `model` | 精确模型 ID |
| `prompt` | string | 否 | `input` | 模型输入提示词 |
| `images` | media[] | 否 | `input` | 图片素材 |
| `videos` | media[] | 否 | `input` | 视频素材 |
| `audios` | media[] | 否 | `input` | 音频素材 |
| `duration` | integer | 否 | `parameters.duration` | 时长 |
| `aspectRatio` | string | 否 | `parameters` | 画幅或尺寸 |

## 上游请求模板逐字段清单

下表由插件请求模板生成，覆盖 body、query、headers 和 multipart 文件声明中的每个字段。

| 上游位置 | 值或转换表达式 |
| --- | --- |
| `create.method` | `"POST"` |
| `create.path` | `"/api/v1/services/aigc/multimodal-generation/generation"` |
| `create.contentType` | `"application/json"` |
| `create.body` | `{"$if":{"condition":{"$eq":[{"$ref":"vars.requestMode"},"openai"]},"then":{"model":{"$ref":"vars.model"},"prompt":{"$ref":"vars.prompt"},"image":{"$omitEmpty":{"$map":{"from":{"$ref":"vars.images"},"as":"m","in":{"$coalesce":[{"$ref":"m.dataUrl"},{"$ref":"m.url"}]}}}},"n":{"$ref":"vars.count"},"size":{"$replace":[{"$coalesce":[{"$ref":"vars.size"},"auto"]},"*","x"]},"prompt_extend":{"$ref":"vars.extend"},"enable_thinking":{"$ref":"vars.thinking"},"prompt_extend_mode":{"$if":{"condition":{"$ref":"vars.extend"},"then":{"$ref":"vars.extendMode"},"else":null}},"watermark":{"$ref":"vars.watermark"},"negative_prompt":{"$omitEmpty":{"$ref":"vars.negative"}},"seed":{"$ref":"vars.seed"}},"else":{"model":{"$ref":"vars.model"},"input":{"messages":[{"role":"user","content":{"$concatArrays":[{"$map":{"from":{"$ref":"vars.images"},"as":"m","in":{"image":{"$coalesce":[{"$ref":"m.dataUrl"},{"$ref":"m.url"}]}}}},[{"text":{"$ref":"vars.prompt"}}]]}}]},"parameters":{"n":{"$ref":"vars.count"},"size":{"$ref":"vars.size"},"prompt_extend":{"$ref":"vars.extend"},"enable_thinking":{"$ref":"vars.thinking"},"prompt_extend_mode":{"$if":{"condition":{"$ref":"vars.extend"},"then":{"$ref":"vars.extendMode"},"else":null}},"watermark":{"$ref":"vars.watermark"},"negative_prompt":{"$omitEmpty":{"$ref":"vars.negative"}},"seed":{"$ref":"vars.seed"}}}}}` |
| `create.headers.X-DashScope-Async` | `{"$if":{"condition":{"$eq":[{"$ref":"vars.requestMode"},"async"]},"then":"enable","else":null}}` |
| `poll.method` | `"GET"` |
| `poll.path` | `"/api/v1/tasks/{{taskId}}"` |
| `poll.contentType` | `"application/json"` |

## Provider 扩展键

- 无额外扩展键。

动态模型或工作流允许使用文档声明的完整 `parameters/input/extra_body` 对象；该对象是协议本身的开放 schema，不会被宿主裁剪。

## 响应映射逐字段清单

| 映射位置 | 上游路径或转换表达式 |
| --- | --- |
| `response.strict` | `true` |
| `response.status` | `{"$switch":{"cases":[{"when":{"$eq":[{"$upper":{"$trim":{"$ref":"response.output.task_status"}}},"UNKNOWN"]},"then":"failed"}],"default":{"$coalesce":[{"$ref":"response.output.task_status"},"succeeded"]}}}` |
| `response.images` | `{"$concatArrays":[{"$flatten":{"$map":{"from":{"$ref":"response.output.choices"},"as":"c","in":{"$map":{"from":{"$ref":"c.message.content"},"as":"m","in":{"$ref":"m.image"}}}}}},{"$map":{"from":{"$ref":"response.output.results"},"as":"m","in":{"$ref":"m.url"}}},{"$map":{"from":{"$ref":"response.data"},"as":"m","in":{"$ref":"m.url"}}}]}` |
| `response.usage` | `{"$ref":"response.usage"}` |
| `response.errorPaths[0]` | `"code"` |
| `response.errorPaths[1]` | `"output.code"` |
| `response.errorPaths[2]` | `"error.code"` |
| `response.errorPaths[3]` | `"error.type"` |
| `response.messagePaths[0]` | `"message"` |
| `response.messagePaths[1]` | `"output.message"` |
| `response.messagePaths[2]` | `"error.message"` |
| `response.resultEphemeral` | `true` |
| `response.resultAnonymous` | `true` |
| `response.taskId` | `{"$coalesce":[{"$ref":"response.output.task_id"},{"$ref":"taskId"}]}` |

## 响应与错误

插件把上游 task/status/text/media/usage 映射为统一结果。临时媒体 URL 标记为 ephemeral，由宿主立即下载持久化。HTTP 错误、业务 code 和 error object 保持失败语义，不包装成成功。

## 兼容边界

仅支持 qwen-image-3.0-pro。同步文生图与 1–3 图编辑；不支持蒙版或质量参数。auto 省略尺寸；内部比例和像素在插件中规范化。保留账号区域主机和路径前缀。宿主负责 SSRF 防护、凭证、计费及持久化。 支持 request_mode=sync/async/openai，默认 sync。异步使用专属 image-generation 端点与 10 秒轮询；兼容接口使用字母 x 和 image 数组。

<!-- YINGCE_MANIFEST_CONTRACT_START -->
## Manifest 完整接口定义

以下 JSON 与插件包内实际 `manifest.json` 逐字段一致，覆盖插件身份、权限、配置、鉴权、参数、校验、创建、Agent、查询、取消、结果下载、响应和 Agent 响应映射。`documentation` 字段的值就是当前完整文档；为避免文档在自身内部无限递归，JSON 中仅用等义占位文本表示正文。

```json
{
  "apiVersion": "yingce.plugin/v4",
  "id": "dashscope-image",
  "name": "DashScope Qwen Image",
  "version": "4.0.0",
  "author": "Alibaba Cloud / 影策",
  "description": "DashScope Qwen Image 独立请求协议插件。",
  "documentation": "<当前插件的完整 documentation，由 README.md 与 docs/interface.md 拼接而成；为避免 JSON 递归，此处不重复展开正文。>",
  "permissions": [
    "generation.run",
    "media.read"
  ],
  "configuration": {
    "fields": [
      {
        "name": "apiKey",
        "type": "secret",
        "label": "API Key",
        "required": true
      }
    ]
  },
  "contributes": {
    "providers": [
      {
        "id": "dashscope-image",
        "label": "DashScope Qwen Image",
        "capabilities": [
          "image"
        ],
        "scopes": [
          "admin.system-channel",
          "user.custom-channel",
          "canvas",
          "creation",
          "agent"
        ],
        "baseUrl": "https://dashscope.aliyuncs.com/compatible-mode/v1",
        "requiresPublicMediaUrls": false,
        "auth": {
          "type": "bearer",
          "field": "apiKey"
        },
        "parameters": [
          {
            "name": "model",
            "type": "string",
            "required": true,
            "mapping": "model",
            "description": "精确模型 ID"
          },
          {
            "name": "prompt",
            "type": "string",
            "required": false,
            "mapping": "input",
            "description": "模型输入提示词"
          },
          {
            "name": "images",
            "type": "media[]",
            "required": false,
            "mapping": "input",
            "description": "图片素材"
          },
          {
            "name": "videos",
            "type": "media[]",
            "required": false,
            "mapping": "input",
            "description": "视频素材"
          },
          {
            "name": "audios",
            "type": "media[]",
            "required": false,
            "mapping": "input",
            "description": "音频素材"
          },
          {
            "name": "duration",
            "type": "integer",
            "required": false,
            "mapping": "parameters.duration",
            "description": "时长"
          },
          {
            "name": "aspectRatio",
            "type": "string",
            "required": false,
            "mapping": "parameters",
            "description": "画幅或尺寸"
          }
        ],
        "modelContract": {
          "mode": "exact",
          "models": [
            "qwen-image-3.0-pro"
          ]
        },
        "modelCatalog": {
          "hosts": [
            "dashscope.aliyuncs.com",
            "dashscope-intl.aliyuncs.com",
            "dashscope-us.aliyuncs.com",
            "maas.aliyuncs.com"
          ],
          "entries": [
            {
              "id": "qwen-image-3.0-pro",
              "displayName": "通义千问图像 3.0 Pro",
              "supportStatus": "ready",
              "supportReason": "Qwen Image 3.0 同步文生图与 1-3 图编辑合同、严格序列化和回归测试已完成",
              "supportedOperations": [
                "text_to_image",
                "image_edit"
              ],
              "documentationPaths": [
                "api-reference/image-generation/qwen-text-to-image.md",
                "api-reference/image-generation/qwen-image-editing.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/multimodal-generation/generation"
            },
            {
              "id": "wan2.7-image-pro",
              "displayName": "万相图像 2.7 Pro",
              "supportStatus": "planned",
              "supportReason": "Wan 2.7 Image 模型专属能力合同尚未完成",
              "supportedOperations": [
                "text_to_image",
                "image_edit"
              ],
              "documentationPaths": [
                "api-reference/image-generation/wan27-image-gen-edit/create-task.md"
              ],
              "catalogSource": "bailian-official-docs",
              "catalogVersion": "2026-08-27",
              "apiPath": "/api/v1/services/aigc/multimodal-generation/generation"
            }
          ]
        },
        "variables": [
          {
            "name": "model",
            "value": {
              "$replace": [
                {
                  "$lower": {
                    "$trim": {
                      "$ref": "request.model"
                    }
                  }
                },
                "models/",
                ""
              ]
            }
          },
          {
            "name": "meta",
            "value": {
              "$merge": [
                {
                  "$ref": "request.extra.generationMetadata"
                },
                {
                  "$ref": "request.providerOptions.dashscope-image"
                }
              ]
            }
          },
          {
            "name": "raw",
            "value": {
              "$ref": "request.extra.rawOptions"
            }
          },
          {
            "name": "images",
            "value": {
              "$map": {
                "from": {
                  "$ref": "request.images"
                },
                "as": "m",
                "in": {
                  "$merge": [
                    {
                      "$ref": "m"
                    },
                    {
                      "value": {
                        "$coalesce": [
                          {
                            "$ref": "m.dataUrl"
                          },
                          {
                            "$ref": "m.url"
                          }
                        ]
                      }
                    }
                  ]
                }
              }
            }
          },
          {
            "name": "videos",
            "value": {
              "$ref": "request.videos"
            }
          },
          {
            "name": "audios",
            "value": {
              "$ref": "request.audios"
            }
          },
          {
            "name": "prompt",
            "value": {
              "$trim": {
                "$if": {
                  "condition": {
                    "$ne": [
                      {
                        "$ref": "request.instructions"
                      },
                      ""
                    ]
                  },
                  "then": {
                    "$concat": [
                      {
                        "$ref": "request.instructions"
                      },
                      "\n\n",
                      {
                        "$ref": "request.prompt"
                      }
                    ]
                  },
                  "else": {
                    "$ref": "request.prompt"
                  }
                }
              }
            }
          },
          {
            "name": "count",
            "value": {
              "$strictInt": {
                "value": {
                  "$ref": "vars.raw.count"
                },
                "default": 1,
                "min": 1,
                "max": 6,
                "allow": []
              }
            }
          },
          {
            "name": "size",
            "value": {
              "$pixelSize": {
                "value": {
                  "$ref": "request.aspectRatio"
                },
                "presets": {
                  "1:1": "1024x1024",
                  "2:3": "1024x1536",
                  "3:2": "1536x1024",
                  "3:4": "768x1024",
                  "4:3": "1024x768",
                  "9:16": "720x1280",
                  "16:9": "1280x720"
                },
                "targetPixels": 1048576,
                "ratioMaxSide": 2048,
                "minPixels": 262144,
                "maxPixels": 4194304,
                "minRatio": 0.125,
                "maxRatio": 8,
                "separator": "*"
              }
            }
          },
          {
            "name": "extend",
            "value": {
              "$strictBool": {
                "value": {
                  "$ref": "vars.meta.prompt_extend"
                },
                "default": true
              }
            }
          },
          {
            "name": "thinking",
            "value": {
              "$strictBool": {
                "value": {
                  "$ref": "vars.meta.enable_thinking"
                },
                "default": {
                  "$ref": "vars.extend"
                }
              }
            }
          },
          {
            "name": "extendMode",
            "value": {
              "$coalesce": [
                {
                  "$lower": {
                    "$trim": {
                      "$ref": "vars.meta.prompt_extend_mode"
                    }
                  }
                },
                "direct"
              ]
            }
          },
          {
            "name": "watermark",
            "value": {
              "$strictBool": {
                "value": {
                  "$ref": "vars.meta.watermark"
                },
                "default": false
              }
            }
          },
          {
            "name": "seed",
            "value": {
              "$strictInt": {
                "value": {
                  "$ref": "vars.meta.seed"
                },
                "default": null,
                "min": 0,
                "max": 2147483647,
                "allow": []
              }
            }
          },
          {
            "name": "negative",
            "value": {
              "$coalesce": [
                {
                  "$ref": "vars.meta.negative_prompt"
                },
                ""
              ]
            }
          },
          {
            "name": "requestMode",
            "value": {
              "$coalesce": [
                {
                  "$ref": "vars.meta.request_mode"
                },
                "sync"
              ]
            }
          }
        ],
        "pollIntervalMs": 10000,
        "validations": [
          {
            "assert": {
              "$and": [
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "vars.images"
                      }
                    },
                    {
                      "$coalesce": [
                        {
                          "$ref": "request.extra.referenceCounts.image"
                        },
                        {
                          "$len": {
                            "$ref": "vars.images"
                          }
                        }
                      ]
                    }
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "vars.videos"
                      }
                    },
                    {
                      "$coalesce": [
                        {
                          "$ref": "request.extra.referenceCounts.video"
                        },
                        {
                          "$len": {
                            "$ref": "vars.videos"
                          }
                        }
                      ]
                    }
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "vars.audios"
                      }
                    },
                    {
                      "$coalesce": [
                        {
                          "$ref": "request.extra.referenceCounts.audio"
                        },
                        {
                          "$len": {
                            "$ref": "vars.audios"
                          }
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "参考素材缺少有效地址，不能静默忽略"
          },
          {
            "assert": {
              "$in": [
                {
                  "$ref": "vars.model"
                },
                [
                  "qwen-image-3.0-pro"
                ]
              ]
            },
            "message": "该图片插件仅支持已验证的 Qwen Image 3.0 Pro"
          },
          {
            "assert": {
              "$gt": [
                {
                  "$len": {
                    "$ref": "vars.prompt"
                  }
                },
                0
              ]
            },
            "message": "必须提供提示词"
          },
          {
            "assert": {
              "$and": [
                {
                  "$lte": [
                    {
                      "$len": {
                        "$ref": "vars.images"
                      }
                    },
                    3
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "vars.videos"
                      }
                    },
                    0
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "vars.audios"
                      }
                    },
                    0
                  ]
                }
              ]
            },
            "message": "Qwen 最多支持 3 张参考图且不支持音视频输入"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$ref": "vars.images"
                      },
                      "as": "m",
                      "where": {
                        "$eq": [
                          {
                            "$ref": "m.role"
                          },
                          "mask"
                        ]
                      }
                    }
                  }
                },
                0
              ]
            },
            "message": "Qwen 不支持蒙版"
          },
          {
            "assert": {
              "$in": [
                {
                  "$lower": {
                    "$trim": {
                      "$ref": "request.quality"
                    }
                  }
                },
                [
                  "",
                  "auto"
                ]
              ]
            },
            "message": "Qwen 不支持质量参数"
          },
          {
            "assert": {
              "$not": {
                "$strictBool": {
                  "value": {
                    "$ref": "vars.raw.transparentBackground"
                  },
                  "default": false
                }
              }
            },
            "message": "Qwen 不支持透明背景"
          },
          {
            "assert": {
              "$in": [
                {
                  "$ref": "vars.extendMode"
                },
                [
                  "direct",
                  "agent"
                ]
              ]
            },
            "message": "改写模式必须为 direct 或 agent"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$not": {
                      "$ref": "vars.extend"
                    }
                  }
                },
                {
                  "$and": [
                    {
                      "$eq": [
                        {
                          "$ref": "vars.meta.prompt_extend_mode"
                        },
                        null
                      ]
                    },
                    {
                      "$not": {
                        "$ref": "vars.thinking"
                      }
                    }
                  ]
                }
              ]
            },
            "message": "关闭改写时不能设置改写模式或开启思考"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$gt": [
                      {
                        "$len": {
                          "$ref": "vars.images"
                        }
                      },
                      0
                    ]
                  }
                },
                {
                  "$ne": [
                    {
                      "$ref": "vars.extendMode"
                    },
                    "agent"
                  ]
                }
              ]
            },
            "message": "图片编辑不支持 agent 改写"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "vars.negative"
                  }
                },
                500
              ]
            },
            "message": "负面提示词不能超过 500 字符"
          },
          {
            "assert": {
              "$validateMedia": {
                "items": {
                  "$ref": "vars.images"
                },
                "maxBytes": 10485760,
                "mimes": [
                  "image/jpeg",
                  "image/jpg",
                  "image/png",
                  "image/webp",
                  "image/bmp",
                  "image/x-ms-bmp",
                  "image/tiff",
                  "image/gif"
                ],
                "allowData": true
              }
            },
            "message": "参考图片不符合限制"
          },
          {
            "assert": {
              "$in": [
                {
                  "$ref": "vars.requestMode"
                },
                [
                  "sync",
                  "async",
                  "openai"
                ]
              ]
            },
            "message": "未知图片请求方式"
          }
        ],
        "create": {
          "basePathSuffixes": [
            "/compatible-mode/v1",
            "/api/v1"
          ],
          "method": "POST",
          "path": "/api/v1/services/aigc/multimodal-generation/generation",
          "contentType": "application/json",
          "body": {
            "$if": {
              "condition": {
                "$eq": [
                  {
                    "$ref": "vars.requestMode"
                  },
                  "openai"
                ]
              },
              "then": {
                "model": {
                  "$ref": "vars.model"
                },
                "prompt": {
                  "$ref": "vars.prompt"
                },
                "image": {
                  "$omitEmpty": {
                    "$map": {
                      "from": {
                        "$ref": "vars.images"
                      },
                      "as": "m",
                      "in": {
                        "$coalesce": [
                          {
                            "$ref": "m.dataUrl"
                          },
                          {
                            "$ref": "m.url"
                          }
                        ]
                      }
                    }
                  }
                },
                "n": {
                  "$ref": "vars.count"
                },
                "size": {
                  "$replace": [
                    {
                      "$coalesce": [
                        {
                          "$ref": "vars.size"
                        },
                        "auto"
                      ]
                    },
                    "*",
                    "x"
                  ]
                },
                "prompt_extend": {
                  "$ref": "vars.extend"
                },
                "enable_thinking": {
                  "$ref": "vars.thinking"
                },
                "prompt_extend_mode": {
                  "$if": {
                    "condition": {
                      "$ref": "vars.extend"
                    },
                    "then": {
                      "$ref": "vars.extendMode"
                    },
                    "else": null
                  }
                },
                "watermark": {
                  "$ref": "vars.watermark"
                },
                "negative_prompt": {
                  "$omitEmpty": {
                    "$ref": "vars.negative"
                  }
                },
                "seed": {
                  "$ref": "vars.seed"
                }
              },
              "else": {
                "model": {
                  "$ref": "vars.model"
                },
                "input": {
                  "messages": [
                    {
                      "role": "user",
                      "content": {
                        "$concatArrays": [
                          {
                            "$map": {
                              "from": {
                                "$ref": "vars.images"
                              },
                              "as": "m",
                              "in": {
                                "image": {
                                  "$coalesce": [
                                    {
                                      "$ref": "m.dataUrl"
                                    },
                                    {
                                      "$ref": "m.url"
                                    }
                                  ]
                                }
                              }
                            }
                          },
                          [
                            {
                              "text": {
                                "$ref": "vars.prompt"
                              }
                            }
                          ]
                        ]
                      }
                    }
                  ]
                },
                "parameters": {
                  "n": {
                    "$ref": "vars.count"
                  },
                  "size": {
                    "$ref": "vars.size"
                  },
                  "prompt_extend": {
                    "$ref": "vars.extend"
                  },
                  "enable_thinking": {
                    "$ref": "vars.thinking"
                  },
                  "prompt_extend_mode": {
                    "$if": {
                      "condition": {
                        "$ref": "vars.extend"
                      },
                      "then": {
                        "$ref": "vars.extendMode"
                      },
                      "else": null
                    }
                  },
                  "watermark": {
                    "$ref": "vars.watermark"
                  },
                  "negative_prompt": {
                    "$omitEmpty": {
                      "$ref": "vars.negative"
                    }
                  },
                  "seed": {
                    "$ref": "vars.seed"
                  }
                }
              }
            }
          },
          "pathTemplate": {
            "$switch": {
              "cases": [
                {
                  "when": {
                    "$eq": [
                      {
                        "$ref": "vars.requestMode"
                      },
                      "async"
                    ]
                  },
                  "then": "/api/v1/services/aigc/image-generation/generation"
                },
                {
                  "when": {
                    "$eq": [
                      {
                        "$ref": "vars.requestMode"
                      },
                      "openai"
                    ]
                  },
                  "then": "/compatible-mode/v1/images/generations"
                }
              ],
              "default": "/api/v1/services/aigc/multimodal-generation/generation"
            }
          },
          "headers": {
            "X-DashScope-Async": {
              "$if": {
                "condition": {
                  "$eq": [
                    {
                      "$ref": "vars.requestMode"
                    },
                    "async"
                  ]
                },
                "then": "enable",
                "else": null
              }
            }
          }
        },
        "poll": {
          "basePathSuffixes": [
            "/compatible-mode/v1",
            "/api/v1"
          ],
          "method": "GET",
          "path": "/api/v1/tasks/{{taskId}}"
        },
        "response": {
          "strict": true,
          "status": {
            "$switch": {
              "cases": [
                {
                  "when": {
                    "$eq": [
                      {
                        "$upper": {
                          "$trim": {
                            "$ref": "response.output.task_status"
                          }
                        }
                      },
                      "UNKNOWN"
                    ]
                  },
                  "then": "failed"
                }
              ],
              "default": {
                "$coalesce": [
                  {
                    "$ref": "response.output.task_status"
                  },
                  "succeeded"
                ]
              }
            }
          },
          "images": {
            "$concatArrays": [
              {
                "$flatten": {
                  "$map": {
                    "from": {
                      "$ref": "response.output.choices"
                    },
                    "as": "c",
                    "in": {
                      "$map": {
                        "from": {
                          "$ref": "c.message.content"
                        },
                        "as": "m",
                        "in": {
                          "$ref": "m.image"
                        }
                      }
                    }
                  }
                }
              },
              {
                "$map": {
                  "from": {
                    "$ref": "response.output.results"
                  },
                  "as": "m",
                  "in": {
                    "$ref": "m.url"
                  }
                }
              },
              {
                "$map": {
                  "from": {
                    "$ref": "response.data"
                  },
                  "as": "m",
                  "in": {
                    "$ref": "m.url"
                  }
                }
              }
            ]
          },
          "usage": {
            "$ref": "response.usage"
          },
          "errorPaths": [
            "code",
            "output.code",
            "error.code",
            "error.type"
          ],
          "messagePaths": [
            "message",
            "output.message",
            "error.message"
          ],
          "resultEphemeral": true,
          "resultAnonymous": true,
          "taskId": {
            "$coalesce": [
              {
                "$ref": "response.output.task_id"
              },
              {
                "$ref": "taskId"
              }
            ]
          }
        }
      }
    ]
  }
}
```
<!-- YINGCE_MANIFEST_CONTRACT_END -->
