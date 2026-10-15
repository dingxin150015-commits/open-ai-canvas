package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"yingce/backend/internal/model"
	"yingce/backend/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestComfyBridgeRequestSurvivesServiceRestart(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+newID()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ComfyBridge{}, &model.ComfyBridgeRequest{}); err != nil {
		t.Fatal(err)
	}
	repo := repository.New(db)
	bridge := model.ComfyBridge{ID: "bridge-1", UserID: "user-1", Name: "Local", TokenHash: "hash", Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := repo.Create(&bridge); err != nil {
		t.Fatal(err)
	}

	first := New(repo, t.TempDir())
	request, err := first.EnqueueComfyBridgeRequest(context.Background(), bridge.UserID, bridge.ID, "task-1", map[string]any{"prompt": "hello"})
	if err != nil {
		t.Fatal(err)
	}

	// 新 Service 没有旧进程的内存队列，仍应从数据库领取同一请求。
	second := New(repo, t.TempDir())
	claimed, err := second.PollComfyBridgeRequest(context.Background(), &bridge, 0)
	if err != nil {
		t.Fatal(err)
	}
	if claimed == nil || claimed.ID != request.ID {
		t.Fatalf("claimed request = %#v, want %q", claimed, request.ID)
	}
	if err := second.CompleteComfyBridgeRequest(bridge.ID, ComfyBridgeCompletion{RequestID: request.ID, Status: "succeeded", Result: map[string]any{"url": "https://example.test/result.png"}}); err != nil {
		t.Fatal(err)
	}

	// 完成结果先于任务 worker 恢复也不会丢失，恢复时复用原请求 ID 且不重复入队。
	third := New(repo, t.TempDir())
	resumed, err := third.enqueueComfyBridgeRequest(context.Background(), bridge.UserID, bridge.ID, "task-1", request.ID, map[string]any{"prompt": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if resumed.ID != request.ID {
		t.Fatalf("resumed request ID = %q, want %q", resumed.ID, request.ID)
	}
	completion, err := third.WaitComfyBridgeRequest(context.Background(), request.ID)
	if err != nil {
		t.Fatal(err)
	}
	if completion.Status != "succeeded" || completion.Result["url"] != "https://example.test/result.png" {
		t.Fatalf("completion = %#v", completion)
	}
	if duplicate, err := third.PollComfyBridgeRequest(context.Background(), &bridge, 0); err != nil || duplicate != nil {
		t.Fatalf("completed request was queued again: request=%#v err=%v", duplicate, err)
	}
}

func TestComfyBridgeTokenAndRevocationStayUserScoped(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+newID()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ComfyBridge{}, &model.ComfyBridgeRequest{}); err != nil {
		t.Fatal(err)
	}
	const validToken = "bridge-secret-with-at-least-thirty-two-bytes"
	bridge := model.ComfyBridge{ID: "bridge-1", UserID: "owner", Name: "Local", TokenHash: hashToken(validToken), Enabled: true}
	if err := db.Create(&bridge).Error; err != nil {
		t.Fatal(err)
	}
	svc := New(repository.New(db), t.TempDir())
	if _, err := svc.AuthenticateComfyBridge("wrong-token-with-at-least-thirty-two-bytes"); err == nil {
		t.Fatal("wrong Bridge token was accepted")
	}
	if _, err := svc.AuthenticateComfyBridge(validToken); err != nil {
		t.Fatalf("valid Bridge token rejected: %v", err)
	}
	if err := svc.RevokeComfyBridge("other-user", bridge.ID); err == nil {
		t.Fatal("another user revoked the Bridge")
	}
	if err := svc.RevokeComfyBridge("owner", bridge.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AuthenticateComfyBridge(validToken); err == nil {
		t.Fatal("revoked Bridge token remained valid")
	}
}

func TestFetchRunningHubAppInfoKeepsAPIKeyOutOfURL(t *testing.T) {
	t.Setenv("CANVAS_ALLOWED_PRIVATE_UPSTREAM_HOSTS", "127.0.0.1")
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", request.Method)
		}
		if request.URL.RawQuery != "" {
			t.Errorf("API Key leaked into URL query: %s", request.URL.RawQuery)
		}
		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body["apiKey"] != "secret-key" || body["webappId"] != "app-1" {
			t.Errorf("body = %#v", body)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"code":0,"data":{"nodeInfoList":[]}}`))
	}))
	defer server.Close()

	result, err := (&Service{}).FetchRunningHubAppInfo(context.Background(), RunningHubWorkflowFetchRequest{BaseURL: server.URL, APIKey: "secret-key", WebappID: "app-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result["webappId"] != "app-1" {
		t.Fatalf("result = %#v", result)
	}
}
