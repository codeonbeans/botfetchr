package tgbot

import (
	"botmediasaver/generated/sqlc"
	"botmediasaver/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrSubscriptionExpired  = errors.New("subscription has expired")
	ErrFeatureLimitExceeded = errors.New("feature limit exceeded")
)

func (b *DefaultBot) IsAccountAllow(ctx context.Context, accountID int64, feature model.Feature, pendingUsage int64) error {
	b.subscriptionMux.Lock() // Ensure that only one goroutine can access the subscription logic at a time
	defer b.subscriptionMux.Unlock()

	txStorage, err := b.storage.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Step 1: Check if the user has an active subscription
	subscriptions, err := txStorage.ListSubscriptions(ctx, sqlc.ListSubscriptionsParams{
		AccountID: pgtype.Int8{Int64: accountID, Valid: true},
		Status:    sqlc.NullSubscriptionStatuses{SubscriptionStatuses: sqlc.SubscriptionStatusesACTIVE, Valid: true},
		Limit:     1,
	})
	if err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		subscription, err := txStorage.CreateSubscription(ctx, sqlc.CreateSubscriptionParams{
			AccountID: accountID,
			PlanID:    model.PlanFree.String(),
			StartDate: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			Status:    sqlc.SubscriptionStatusesACTIVE,
		})
		if err != nil {
			return err
		}

		subscriptions = append(subscriptions, subscription)
	}

	subscription := subscriptions[0]
	if subscription.EndDate.Valid && subscription.EndDate.Time.Before(time.Now()) {
		// Subscription has expired
		return ErrSubscriptionExpired
	}

	// Step 2: Check if user has access to the feature and not exceeded the limit
	planFeatures, err := txStorage.ListPlanFeatures(ctx, sqlc.ListPlanFeaturesParams{
		PlanID:  pgtype.Text{String: subscription.PlanID, Valid: true},
		Feature: pgtype.Text{String: feature.String(), Valid: true},
		Limit:   1,
	})
	if err != nil {
		return err
	}
	if len(planFeatures) == 0 {
		return fmt.Errorf("feature %s is not available in plan %s", feature.String(), subscription.PlanID)
	}
	planFeature := planFeatures[0]

	usages, err := txStorage.ListAccountUsages(ctx, sqlc.ListAccountUsagesParams{
		AccountID: pgtype.Int8{Int64: accountID, Valid: true},
		Feature:   pgtype.Text{String: feature.String(), Valid: true},
		Limit:     1,
	})
	if err != nil {
		return err
	}
	if len(usages) == 0 {
		usage, err := txStorage.CreateAccountUsage(ctx, sqlc.CreateAccountUsageParams{
			AccountID: accountID,
			Feature:   feature.String(),
			Usage:     0,
		})
		if err != nil {
			return err
		}

		usages = append(usages, usage)
	}
	usage := usages[0]

	// reset usage if it meet the reset condition (days_to_reset)
	if planFeature.DaysToReset > 0 && usage.ResetAt.Time.AddDate(0, 0, int(planFeature.DaysToReset)).Before(time.Now()) {
		usage, err = txStorage.UpdateAccountUsage(ctx, sqlc.UpdateAccountUsageParams{
			ID:      usage.ID,
			Usage:   pgtype.Int8{Int64: 0, Valid: true},
			ResetAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}
	}

	if usage.Usage+pendingUsage > planFeature.Limit && planFeature.Limit > 0 {
		return ErrFeatureLimitExceeded
	}

	// Step 3: Update usage
	if _, err = txStorage.UpdateAccountUsage(ctx, sqlc.UpdateAccountUsageParams{
		ID:    usage.ID,
		Usage: pgtype.Int8{Int64: usage.Usage + pendingUsage, Valid: true},
	}); err != nil {
		return err
	}

	return txStorage.Commit(ctx)
}
