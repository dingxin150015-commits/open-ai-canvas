package service

import (
	"context"
	"testing"

	"infinite-canvas/backend/internal/model"
)

func TestV127BatchDeleteRejectsPlannedWithoutPartialDeletion(t *testing.T) {
	svc, db := newChannelModelTestService(t)
	admin := &model.User{ID: "admin", Role: model.UserRoleAdmin}
	channel := model.ModelChannel{ID: "merge-contract", UserID: admin.ID, Scope: model.ChannelScopeSystem, Name: "Test", Enabled: true}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	items := []model.ChannelModel{
		{ID: "merge-ready", ChannelID: channel.ID, ModelKey: "ready", SupportStatus: model.ChannelModelSupportReady},
		{ID: "merge-planned", ChannelID: channel.ID, ModelKey: "planned", SupportStatus: model.ChannelModelSupportPlanned},
	}
	if err := db.Create(&items).Error; err != nil {
		t.Fatal(err)
	}
	deleted, err := svc.DeleteAdminChannelModels(admin, channel.ID, []string{items[0].ID, items[1].ID})
	if err == nil || deleted != 0 {
		t.Fatalf("deleted=%d err=%v, want atomic rejection", deleted, err)
	}
	var count int64
	if err := db.Model(&model.ChannelModel{}).Where("channel_id = ?", channel.ID).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("remaining=%d, want2", count)
	}
}

func TestV127DashScopeHydrationRequiresPublicURL(t *testing.T) {
	policy := providerMediaHydrationPolicyFor(context.Background(), canvasGenerationInput{Mode: "video", Config: providerConfig{InterfaceType: "dashscope-video"}})
	if !policy.requireURL || !policy.preferURL {
		t.Fatalf("policy=%#v", policy)
	}
}
