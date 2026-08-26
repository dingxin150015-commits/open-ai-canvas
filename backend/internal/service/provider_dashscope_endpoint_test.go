package service

import "testing"

func TestDashScopeNativeEndpointPreservesHostAndReplacesProtocolPath(t *testing.T) {
	tests := []struct {
		name string
		base string
		path string
		want string
	}{
		{
			name: "workspace MaaS compatible endpoint",
			base: "https://llm-workspace.cn-beijing.maas.aliyuncs.com/compatible-mode/v1",
			path: "/api/v1/services/aigc/video-generation/video-synthesis",
			want: "https://llm-workspace.cn-beijing.maas.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis",
		},
		{
			name: "public compatible endpoint",
			base: "https://dashscope.aliyuncs.com/compatible-mode/v1/",
			path: "api/v1/tasks/task-1",
			want: "https://dashscope.aliyuncs.com/api/v1/tasks/task-1",
		},
		{
			name: "already native base",
			base: "https://token-plan.cn-beijing.maas.aliyuncs.com/api/v1",
			path: "/api/v1/tasks/task-1",
			want: "https://token-plan.cn-beijing.maas.aliyuncs.com/api/v1/tasks/task-1",
		},
		{
			name: "custom proxy prefix",
			base: "https://proxy.example.com/gateway/compatible-mode/v1",
			path: "/api/v1/tasks/task-1",
			want: "https://proxy.example.com/gateway/api/v1/tasks/task-1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := dashScopeNativeEndpoint(test.base, test.path)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("dashScopeNativeEndpoint() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestDashScopeNativeEndpointRejectsAmbiguousBaseURL(t *testing.T) {
	for _, base := range []string{"not-a-url", "https://dashscope.aliyuncs.com/compatible-mode/v1?token=secret", "https://dashscope.aliyuncs.com/#fragment"} {
		if _, err := dashScopeNativeEndpoint(base, "/api/v1/tasks/task-1"); err == nil {
			t.Fatalf("dashScopeNativeEndpoint(%q) error = nil", base)
		}
	}
}
