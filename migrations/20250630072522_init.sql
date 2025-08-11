-- +goose Up
-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "account";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "subscription";

-- CreateEnum
CREATE TYPE "subscription"."plan_intervals" AS ENUM ('MONTHLY', 'YEARLY');

-- CreateEnum
CREATE TYPE "subscription"."statuses" AS ENUM ('ACTIVE', 'CANCELED', 'EXPIRED', 'TRIALING');

-- CreateTable
CREATE TABLE "account"."telegrams"
(
  "id"            BIGSERIAL   NOT NULL,
  "telegram_id"   BIGINT      NOT NULL,
  "is_bot"        BOOLEAN     NOT NULL DEFAULT false,
  "first_name"    VARCHAR(64) NOT NULL,
  "last_name"     VARCHAR(64) NOT NULL,
  "username"      VARCHAR(32),
  "language_code" VARCHAR(5)  NOT NULL DEFAULT 'en',
  "photo_url"     TEXT,
  "is_premium"    BOOLEAN     NOT NULL DEFAULT false,
  "created_at"    TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT "telegrams_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account"."usage"
(
  "id"         BIGSERIAL NOT NULL,
  "account_id" BIGINT    NOT NULL,
  "feature"    TEXT      NOT NULL,
  "usage"      BIGINT    NOT NULL DEFAULT 0,
  "reset_at"   TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT "usage_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "subscription"."plans"
(
  "id" TEXT NOT NULL,

  CONSTRAINT "plans_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "subscription"."plan_prices"
(
  "id"       BIGSERIAL                       NOT NULL,
  "plan_id"  TEXT                            NOT NULL,
  "price"    DOUBLE PRECISION                NOT NULL,
  "interval" "subscription"."plan_intervals" NOT NULL,

  CONSTRAINT "plan_prices_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "subscription"."plan_features"
(
  "id"            BIGSERIAL NOT NULL,
  "plan_id"       TEXT      NOT NULL,
  "feature"       TEXT      NOT NULL,
  "limit"         BIGINT    NOT NULL DEFAULT 0,
  "days_to_reset" INTEGER   NOT NULL DEFAULT 0,

  CONSTRAINT "plan_features_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "subscription"."subscriptions"
(
  "id"         TEXT                      NOT NULL DEFAULT gen_random_uuid(),
  "account_id" BIGINT                    NOT NULL,
  "plan_id"    TEXT                      NOT NULL,
  "status"     "subscription"."statuses" NOT NULL,
  "start_date" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "end_date"   TIMESTAMPTZ(3),
  "cancel_at"  TIMESTAMPTZ(3),

  CONSTRAINT "subscriptions_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "subscription"."invoices"
(
  "id"              TEXT             NOT NULL DEFAULT gen_random_uuid(),
  "subscription_id" TEXT             NOT NULL,
  "amount"          DOUBLE PRECISION NOT NULL,
  "issued_at"       TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "paid"            BOOLEAN          NOT NULL DEFAULT false,

  CONSTRAINT "invoices_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "telegrams_telegram_id_key" ON "account"."telegrams" ("telegram_id");

-- CreateIndex
CREATE INDEX "telegrams_username_idx" ON "account"."telegrams" ("username");

-- CreateIndex
CREATE INDEX "telegrams_telegram_id_idx" ON "account"."telegrams" ("telegram_id");

-- CreateIndex
CREATE INDEX "telegrams_first_name_last_name_idx" ON "account"."telegrams" ("first_name", "last_name");

-- CreateIndex
CREATE INDEX "telegrams_is_premium_idx" ON "account"."telegrams" ("is_premium");

-- CreateIndex
CREATE INDEX "telegrams_created_at_idx" ON "account"."telegrams" ("created_at");

-- CreateIndex
CREATE INDEX "usage_account_id_idx" ON "account"."usage" ("account_id");

-- CreateIndex
CREATE INDEX "usage_feature_idx" ON "account"."usage" ("feature");

-- CreateIndex
CREATE INDEX "usage_reset_at_idx" ON "account"."usage" ("reset_at");

-- CreateIndex
CREATE UNIQUE INDEX "usage_account_id_feature_key" ON "account"."usage" ("account_id", "feature");

-- CreateIndex
CREATE INDEX "plan_features_plan_id_idx" ON "subscription"."plan_features" ("plan_id");

-- CreateIndex
CREATE UNIQUE INDEX "plan_features_plan_id_feature_key" ON "subscription"."plan_features" ("plan_id", "feature");

-- CreateIndex
CREATE INDEX "subscriptions_account_id_idx" ON "subscription"."subscriptions" ("account_id");

-- CreateIndex
CREATE INDEX "subscriptions_plan_id_idx" ON "subscription"."subscriptions" ("plan_id");

-- CreateIndex
CREATE INDEX "subscriptions_status_idx" ON "subscription"."subscriptions" ("status");

-- CreateIndex
CREATE INDEX "subscriptions_start_date_idx" ON "subscription"."subscriptions" ("start_date");

-- CreateIndex
CREATE INDEX "subscriptions_end_date_idx" ON "subscription"."subscriptions" ("end_date");

-- CreateIndex
CREATE INDEX "subscriptions_cancel_at_idx" ON "subscription"."subscriptions" ("cancel_at");

-- CreateIndex
CREATE INDEX "invoices_subscription_id_idx" ON "subscription"."invoices" ("subscription_id");

-- CreateIndex
CREATE INDEX "invoices_issued_at_idx" ON "subscription"."invoices" ("issued_at");

-- CreateIndex
CREATE INDEX "invoices_paid_idx" ON "subscription"."invoices" ("paid");

-- AddForeignKey
ALTER TABLE "account"."usage"
  ADD CONSTRAINT "usage_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"."telegrams" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "subscription"."plan_prices"
  ADD CONSTRAINT "plan_prices_plan_id_fkey" FOREIGN KEY ("plan_id") REFERENCES "subscription"."plans" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "subscription"."plan_features"
  ADD CONSTRAINT "plan_features_plan_id_fkey" FOREIGN KEY ("plan_id") REFERENCES "subscription"."plans" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "subscription"."subscriptions"
  ADD CONSTRAINT "subscriptions_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"."telegrams" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "subscription"."subscriptions"
  ADD CONSTRAINT "subscriptions_plan_id_fkey" FOREIGN KEY ("plan_id") REFERENCES "subscription"."plans" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "subscription"."invoices"
  ADD CONSTRAINT "invoices_subscription_id_fkey" FOREIGN KEY ("subscription_id") REFERENCES "subscription"."subscriptions" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose Down
-- RemoveForeignKey
ALTER TABLE "subscription"."invoices" DROP CONSTRAINT "invoices_subscription_id_fkey";

-- RemoveForeignKey
ALTER TABLE "subscription"."subscriptions" DROP CONSTRAINT "subscriptions_plan_id_fkey";

-- RemoveForeignKey
ALTER TABLE "subscription"."subscriptions" DROP CONSTRAINT "subscriptions_account_id_fkey";

-- RemoveForeignKey
ALTER TABLE "subscription"."plan_features" DROP CONSTRAINT "plan_features_plan_id_fkey";

-- RemoveForeignKey
ALTER TABLE "subscription"."plan_prices" DROP CONSTRAINT "plan_prices_plan_id_fkey";

-- RemoveForeignKey
ALTER TABLE "account"."usage" DROP CONSTRAINT "usage_account_id_fkey";

-- DropIndex
DROP INDEX "subscription"."invoices_paid_idx";

-- DropIndex
DROP INDEX "subscription"."invoices_issued_at_idx";

-- DropIndex
DROP INDEX "subscription"."invoices_subscription_id_idx";

-- DropIndex
DROP INDEX "subscription"."subscriptions_cancel_at_idx";

-- DropIndex
DROP INDEX "subscription"."subscriptions_end_date_idx";

-- DropIndex
DROP INDEX "subscription"."subscriptions_start_date_idx";

-- DropIndex
DROP INDEX "subscription"."subscriptions_status_idx";

-- DropIndex
DROP INDEX "subscription"."subscriptions_plan_id_idx";

-- DropIndex
DROP INDEX "subscription"."subscriptions_account_id_idx";

-- DropIndex
DROP INDEX "subscription"."plan_features_plan_id_feature_key";

-- DropIndex
DROP INDEX "subscription"."plan_features_plan_id_idx";

-- DropIndex
DROP INDEX "account"."usage_account_id_feature_key";

-- DropIndex
DROP INDEX "account"."usage_reset_at_idx";

-- DropIndex
DROP INDEX "account"."usage_feature_idx";

-- DropIndex
DROP INDEX "account"."usage_account_id_idx";

-- DropIndex
DROP INDEX "account"."telegrams_created_at_idx";

-- DropIndex
DROP INDEX "account"."telegrams_is_premium_idx";

-- DropIndex
DROP INDEX "account"."telegrams_first_name_last_name_idx";

-- DropIndex
DROP INDEX "account"."telegrams_telegram_id_idx";

-- DropIndex
DROP INDEX "account"."telegrams_username_idx";

-- DropIndex
DROP INDEX "account"."telegrams_telegram_id_key";

-- DropTable
DROP TABLE "subscription"."invoices";

-- DropTable
DROP TABLE "subscription"."subscriptions";

-- DropTable
DROP TABLE "subscription"."plan_features";

-- DropTable
DROP TABLE "subscription"."plan_prices";

-- DropTable
DROP TABLE "subscription"."plans";

-- DropTable
DROP TABLE "account"."usage";

-- DropTable
DROP TABLE "account"."telegrams";

-- DropEnum
DROP TYPE "subscription"."statuses";

-- DropEnum
DROP TYPE "subscription"."plan_intervals";

-- DropSchema
DROP SCHEMA "subscription";

-- DropSchema
DROP SCHEMA "account";
