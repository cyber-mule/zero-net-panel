package subscriptionutil

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

const (
	orderMetaSubscriptionID         = "subscription_id"
	orderMetaSubscriptionAction     = "subscription_action"
	orderMetaSubscriptionPlanName   = "subscription_plan_name"
	orderMetaSubscriptionExpiresAt  = "subscription_expires_at"
	orderMetaSubscriptionTemplateID = "subscription_template_id"
	orderMetaSubscriptionRefreshed  = "subscription_refreshed_at"
)

// ProvisionResult records subscription provisioning output for an order.
type ProvisionResult struct {
	Order        repository.Order
	Subscription repository.Subscription
	Action       string
}

// EnsureOrderSubscription creates or renews a subscription for a paid order.
func EnsureOrderSubscription(ctx context.Context, repos *repository.Repositories, order repository.Order, items []repository.OrderItem) (ProvisionResult, error) {
	var result ProvisionResult
	if repos == nil {
		return result, errors.New("subscriptionutil: repositories required")
	}
	if order.ID == 0 {
		return result, repository.ErrInvalidArgument
	}

	lockedOrder, err := repos.Order.GetForUpdate(ctx, order.ID)
	if err != nil {
		return result, err
	}

	if !strings.EqualFold(lockedOrder.Status, repository.OrderStatusPaid) {
		return result, repository.ErrInvalidArgument
	}

	paidAt := time.Now().UTC()
	if lockedOrder.PaidAt != nil && !lockedOrder.PaidAt.IsZero() {
		paidAt = lockedOrder.PaidAt.UTC()
	}

	if subID := metadataUint64(lockedOrder.Metadata, orderMetaSubscriptionID); subID != 0 {
		sub, err := repos.Subscription.Get(ctx, subID)
		if err == nil {
			result.Order = lockedOrder
			result.Subscription = sub
			result.Action = "existing"
			return result, nil
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return result, err
		}
	}

	info, err := buildPlanInfo(lockedOrder, items)
	if err != nil {
		return result, err
	}

	defaultTemplateID, availableTemplateIDs, err := loadTemplates(ctx, repos)
	if err != nil {
		return result, err
	}

	existing, found, err := findEligibleSubscription(ctx, repos, lockedOrder.UserID, info.PlanID, info.PlanName)
	if err != nil {
		return result, err
	}

	now := time.Now().UTC()
	var subscription repository.Subscription
	action := "created"
	if found {
		subscription, err = renewSubscription(ctx, repos, existing, info, paidAt, now, defaultTemplateID, availableTemplateIDs)
		if err != nil {
			return result, err
		}
		action = "renewed"
	} else {
		subscription, err = createSubscription(ctx, repos, lockedOrder.UserID, info, paidAt, now, defaultTemplateID, availableTemplateIDs)
		if err != nil {
			return result, err
		}
	}

	metadataPatch := map[string]any{
		orderMetaSubscriptionID:         subscription.ID,
		orderMetaSubscriptionAction:     action,
		orderMetaSubscriptionPlanName:   subscription.PlanName,
		orderMetaSubscriptionTemplateID: subscription.TemplateID,
		orderMetaSubscriptionRefreshed:  now.Unix(),
	}
	if !subscription.ExpiresAt.IsZero() {
		metadataPatch[orderMetaSubscriptionExpiresAt] = subscription.ExpiresAt.UTC().Unix()
	}

	lockedOrder.Metadata = mergeMetadata(lockedOrder.Metadata, metadataPatch)
	updatedOrder, err := repos.Order.Save(ctx, lockedOrder)
	if err != nil {
		return result, err
	}

	result.Order = updatedOrder
	result.Subscription = subscription
	result.Action = action
	return result, nil
}

type planInfo struct {
	PlanID            uint64
	PlanName          string
	Name              string
	DurationValue     int
	DurationUnit      string
	TrafficLimitBytes int64
	DevicesLimit      int
	Quantity          int
	PlanSnapshot      map[string]any
}

func buildPlanInfo(order repository.Order, items []repository.OrderItem) (planInfo, error) {
	info := planInfo{
		Quantity: 1,
	}

	var planItem *repository.OrderItem
	for i := range items {
		if strings.EqualFold(items[i].ItemType, "plan") {
			planItem = &items[i]
			break
		}
	}
	if planItem == nil && len(items) > 0 {
		planItem = &items[0]
	}

	if planItem != nil {
		if planItem.ItemID != 0 {
			info.PlanID = planItem.ItemID
		}
		if planItem.Quantity > 0 {
			info.Quantity = planItem.Quantity
		}
		info.PlanName = strings.TrimSpace(planItem.Name)
		info.Name = info.PlanName
		if unit := normalizeDurationUnit(stringFromMap(planItem.Metadata, "duration_unit")); unit != "" {
			info.DurationUnit = unit
		}
		if value, ok := intFromAny(planItem.Metadata["duration_value"]); ok {
			info.DurationValue = value
		}
		if info.DurationValue == 0 {
			if value, ok := intFromAny(planItem.Metadata["duration_days"]); ok {
				info.DurationValue = value
				if info.DurationUnit == "" {
					info.DurationUnit = repository.DurationUnitDay
				}
			}
		}
		if value, ok := int64FromAny(planItem.Metadata["traffic_limit_bytes"]); ok {
			info.TrafficLimitBytes = value
		}
		if value, ok := intFromAny(planItem.Metadata["devices_limit"]); ok {
			info.DevicesLimit = value
		}
	}

	if info.PlanName == "" {
		info.PlanName = stringFromMap(order.PlanSnapshot, "name", "plan_name", "plan")
		info.Name = info.PlanName
	}
	if info.PlanID == 0 && order.PlanID != nil {
		info.PlanID = *order.PlanID
	}
	if info.PlanName == "" && order.PlanID != nil {
		info.PlanName = fmt.Sprintf("plan-%d", *order.PlanID)
		info.Name = info.PlanName
	}

	if info.DurationUnit == "" {
		if unit := normalizeDurationUnit(stringFromMap(order.PlanSnapshot, "duration_unit")); unit != "" {
			info.DurationUnit = unit
		}
	}
	if info.DurationValue == 0 {
		if value, ok := intFromAny(order.PlanSnapshot["duration_value"]); ok {
			info.DurationValue = value
		}
	}
	if info.DurationValue == 0 {
		if value, ok := intFromAny(order.PlanSnapshot["duration_days"]); ok {
			info.DurationValue = value
			if info.DurationUnit == "" {
				info.DurationUnit = repository.DurationUnitDay
			}
		}
	}
	if info.TrafficLimitBytes == 0 {
		if value, ok := int64FromAny(order.PlanSnapshot["traffic_limit_bytes"]); ok {
			info.TrafficLimitBytes = value
		}
	}
	if info.DevicesLimit == 0 {
		if value, ok := intFromAny(order.PlanSnapshot["devices_limit"]); ok {
			info.DevicesLimit = value
		}
	}
	if info.PlanSnapshot == nil && len(order.PlanSnapshot) > 0 {
		info.PlanSnapshot = ClonePlanSnapshot(order.PlanSnapshot)
	}
	if info.PlanSnapshot == nil && info.PlanID != 0 {
		info.PlanSnapshot = map[string]any{
			"id":   info.PlanID,
			"name": info.PlanName,
		}
	}

	if info.PlanName == "" || info.PlanID == 0 {
		return planInfo{}, repository.ErrInvalidArgument
	}
	if info.Quantity <= 0 {
		info.Quantity = 1
	}
	if info.DurationValue > 0 && info.DurationUnit == "" {
		info.DurationUnit = repository.DurationUnitDay
	}
	if info.DevicesLimit <= 0 {
		info.DevicesLimit = 1
	}

	return info, nil
}

func loadTemplates(ctx context.Context, repos *repository.Repositories) (uint64, []uint64, error) {
	opts := repository.ListTemplatesOptions{PerPage: 100, IncludeDrafts: false}
	templates, _, err := repos.SubscriptionTemplate.List(ctx, opts)
	if err != nil {
		return 0, nil, err
	}
	if len(templates) == 0 {
		opts.IncludeDrafts = true
		templates, _, err = repos.SubscriptionTemplate.List(ctx, opts)
		if err != nil {
			return 0, nil, err
		}
	}
	if len(templates) == 0 {
		return 0, nil, repository.ErrInvalidArgument
	}

	available := make([]uint64, 0, len(templates))
	defaultID := uint64(0)
	for _, tpl := range templates {
		available = append(available, tpl.ID)
		if tpl.IsDefault && defaultID == 0 {
			defaultID = tpl.ID
		}
	}
	if defaultID == 0 {
		defaultID = templates[0].ID
	}

	return defaultID, available, nil
}

func findEligibleSubscription(ctx context.Context, repos *repository.Repositories, userID uint64, planID uint64, planName string) (repository.Subscription, bool, error) {
	subs, _, err := repos.Subscription.ListByUser(ctx, userID, repository.ListSubscriptionsOptions{
		PerPage: 100,
		Sort:    "updated_at",
	})
	if err != nil {
		return repository.Subscription{}, false, err
	}

	for i := range subs {
		if !isEligible(subs[i]) {
			continue
		}
		if planID > 0 && subs[i].PlanID == planID {
			return subs[i], true, nil
		}
		if planID == 0 && planName != "" && strings.EqualFold(subs[i].PlanName, planName) {
			return subs[i], true, nil
		}
	}

	for i := range subs {
		if !isEligible(subs[i]) {
			continue
		}
		return subs[i], true, nil
	}

	return repository.Subscription{}, false, nil
}

func isEligible(sub repository.Subscription) bool {
	return !strings.EqualFold(sub.Status, "disabled")
}

func normalizeDurationUnit(unit string) string {
	unit = strings.TrimSpace(strings.ToLower(unit))
	switch unit {
	case "hours":
		return repository.DurationUnitHour
	case "days":
		return repository.DurationUnitDay
	case "months":
		return repository.DurationUnitMonth
	case "years":
		return repository.DurationUnitYear
	default:
		return unit
	}
}

func addDuration(base time.Time, unit string, value int) (time.Time, error) {
	if value <= 0 {
		return base, nil
	}
	switch normalizeDurationUnit(unit) {
	case repository.DurationUnitHour:
		return base.Add(time.Duration(value) * time.Hour), nil
	case repository.DurationUnitDay:
		return base.Add(time.Duration(value) * 24 * time.Hour), nil
	case repository.DurationUnitMonth:
		return base.AddDate(0, value, 0), nil
	case repository.DurationUnitYear:
		return base.AddDate(value, 0, 0), nil
	default:
		return time.Time{}, repository.ErrInvalidArgument
	}
}

func renewSubscription(ctx context.Context, repos *repository.Repositories, sub repository.Subscription, info planInfo, paidAt, now time.Time, defaultTemplateID uint64, available []uint64) (repository.Subscription, error) {
	expiresAt := sub.ExpiresAt
	if info.DurationValue > 0 {
		base := paidAt
		if !sub.ExpiresAt.IsZero() && sub.ExpiresAt.After(paidAt) {
			base = sub.ExpiresAt
		}
		updated, err := addDuration(base, info.DurationUnit, info.DurationValue*info.Quantity)
		if err != nil {
			return repository.Subscription{}, err
		}
		expiresAt = updated
	}

	status := sub.Status
	if status == "" || strings.EqualFold(status, "expired") {
		status = "active"
	}
	if !expiresAt.IsZero() && expiresAt.After(now) {
		status = "active"
	}

	trafficTotal := info.TrafficLimitBytes * int64(info.Quantity)
	trafficUsed := sub.TrafficUsedBytes
	if !sub.ExpiresAt.IsZero() && sub.ExpiresAt.Before(paidAt) {
		trafficUsed = 0
	}
	if trafficTotal < trafficUsed {
		trafficTotal = trafficUsed
	}

	templateID := sub.TemplateID
	if templateID == 0 {
		templateID = defaultTemplateID
	}

	availableTemplates := sub.AvailableTemplateIDs
	if len(availableTemplates) == 0 {
		availableTemplates = append([]uint64(nil), available...)
	}
	if !containsUint64(availableTemplates, templateID) {
		availableTemplates = append(availableTemplates, templateID)
	}

	name := sub.Name
	if strings.TrimSpace(name) == "" || strings.EqualFold(strings.TrimSpace(name), strings.TrimSpace(sub.PlanName)) {
		name = info.Name
	}
	planName := sub.PlanName
	if info.PlanName != "" {
		planName = info.PlanName
	}
	planID := sub.PlanID
	if info.PlanID != 0 {
		planID = info.PlanID
	}

	planSnapshot := sub.PlanSnapshot
	if info.PlanSnapshot != nil {
		planSnapshot = ClonePlanSnapshot(info.PlanSnapshot)
	}

	devicesLimit := sub.DevicesLimit
	if info.DevicesLimit > 0 {
		devicesLimit = info.DevicesLimit
	}
	if devicesLimit <= 0 {
		devicesLimit = 1
	}

	input := repository.UpdateSubscriptionInput{
		Status:               &status,
		Name:                 &name,
		PlanName:             &planName,
		PlanID:               &planID,
		PlanSnapshot:         &planSnapshot,
		TemplateID:           &templateID,
		AvailableTemplateIDs: &availableTemplates,
		ExpiresAt:            &expiresAt,
		TrafficTotalBytes:    &trafficTotal,
		TrafficUsedBytes:     &trafficUsed,
		DevicesLimit:         &devicesLimit,
		LastRefreshedAt:      &now,
	}

	updated, err := repos.Subscription.Update(ctx, sub.ID, input)
	if err != nil {
		return repository.Subscription{}, err
	}
	if IsSubscriptionEffective(updated, now) {
		if err := repos.Subscription.DisableOtherActive(ctx, updated.UserID, updated.ID); err != nil {
			return repository.Subscription{}, err
		}
	}
	return updated, nil
}

func createSubscription(ctx context.Context, repos *repository.Repositories, userID uint64, info planInfo, paidAt, now time.Time, defaultTemplateID uint64, available []uint64) (repository.Subscription, error) {
	token, err := generateToken()
	if err != nil {
		return repository.Subscription{}, err
	}

	expiresAt := time.Time{}
	if info.DurationValue > 0 {
		updated, err := addDuration(paidAt, info.DurationUnit, info.DurationValue*info.Quantity)
		if err != nil {
			return repository.Subscription{}, err
		}
		expiresAt = updated
	}

	trafficTotal := info.TrafficLimitBytes * int64(info.Quantity)
	devicesLimit := info.DevicesLimit
	if devicesLimit <= 0 {
		devicesLimit = 1
	}

	subscription := repository.Subscription{
		UserID:               userID,
		Name:                 strings.TrimSpace(info.Name),
		PlanName:             strings.TrimSpace(info.PlanName),
		PlanID:               info.PlanID,
		PlanSnapshot:         ClonePlanSnapshot(info.PlanSnapshot),
		Status:               "active",
		TemplateID:           defaultTemplateID,
		AvailableTemplateIDs: append([]uint64(nil), available...),
		Token:                token,
		ExpiresAt:            expiresAt,
		TrafficTotalBytes:    trafficTotal,
		TrafficUsedBytes:     0,
		DevicesLimit:         devicesLimit,
		LastRefreshedAt:      now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if subscription.Name == "" {
		subscription.Name = subscription.PlanName
	}
	if subscription.PlanName == "" {
		subscription.PlanName = subscription.Name
	}

	created, err := repos.Subscription.Create(ctx, subscription)
	if err != nil {
		return repository.Subscription{}, err
	}
	if IsSubscriptionEffective(created, now) {
		if err := repos.Subscription.DisableOtherActive(ctx, created.UserID, created.ID); err != nil {
			return repository.Subscription{}, err
		}
	}
	return created, nil
}

func generateToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func mergeMetadata(base map[string]any, patch map[string]any) map[string]any {
	if base == nil {
		base = make(map[string]any, len(patch))
	}
	for k, v := range patch {
		base[k] = v
	}
	return base
}

func containsUint64(list []uint64, target uint64) bool {
	for _, value := range list {
		if value == target {
			return true
		}
	}
	return false
}

func metadataUint64(meta map[string]any, key string) uint64 {
	if meta == nil {
		return 0
	}
	value, ok := meta[key]
	if !ok {
		return 0
	}
	result, _ := uint64FromAny(value)
	return result
}

func stringFromMap(data map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := data[key]; ok {
			if str, ok := value.(string); ok && strings.TrimSpace(str) != "" {
				return strings.TrimSpace(str)
			}
		}
	}
	return ""
}

func intFromAny(value any) (int, bool) {
	if value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case int32:
		return int(v), true
	case float64:
		return int(v), true
	case float32:
		return int(v), true
	case json.Number:
		parsed, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return 0, false
		}
		return int(parsed), true
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func int64FromAny(value any) (int64, bool) {
	if value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case int32:
		return int64(v), true
	case float64:
		return int64(v), true
	case float32:
		return int64(v), true
	case json.Number:
		parsed, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	case string:
		parsed, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func uint64FromAny(value any) (uint64, bool) {
	if value == nil {
		return 0, false
	}
	switch v := value.(type) {
	case uint64:
		return v, true
	case uint:
		return uint64(v), true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case string:
		parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	case json.Number:
		parsed, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
