-- +goose Up
-- +goose StatementBegin

-- Insert initial plans into the subscription.plans table
INSERT INTO "subscription"."plans" (id)
VALUES
    ('PlanFree'),
    ('PlanPro'),
    ('PlanLifetime');

-- Insert initial plan features into the subscription.plan_features table
INSERT INTO "subscription"."plan_features" (plan_id, feature, "limit", days_to_reset)
VALUES
    ('PlanFree', 'FeatureGetMedia', 5, 1),
    ('PlanPro', 'FeatureGetMedia', 100, 1),
    ('PlanLifetime', 'FeatureGetMedia', 0, 0)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM "subscription"."plans"
WHERE id IN ('PlanFree', 'PlanPro', 'PlanLifetime');

-- +goose StatementEnd
