-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "account";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "resource";

-- CreateSchema
CREATE SCHEMA IF NOT EXISTS "subscription";

-- CreateEnum
CREATE TYPE "subscription"."plan_intervals" AS ENUM ('MONTHLY', 'YEARLY');

-- CreateEnum
CREATE TYPE "subscription"."statuses" AS ENUM ('ACTIVE', 'CANCELED', 'EXPIRED', 'TRIALING');

-- CreateEnum
CREATE TYPE "resource"."types" AS ENUM ('USER_AGENT', 'PROXY');

-- CreateTable
CREATE TABLE "account"."telegrams" (
    "id" BIGSERIAL NOT NULL,
    "telegram_id" BIGINT NOT NULL,
    "is_bot" BOOLEAN NOT NULL DEFAULT false,
    "first_name" VARCHAR(64) NOT NULL,
    "last_name" VARCHAR(64) NOT NULL,
    "username" VARCHAR(32),
    "language_code" VARCHAR(5) NOT NULL DEFAULT 'en',
    "photo_url" TEXT,
    "is_premium" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "telegrams_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "subscription"."plans" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "price" DOUBLE PRECISION NOT NULL,
    "interval" "subscription"."plan_intervals" NOT NULL,
    "description" TEXT,

    CONSTRAINT "plans_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "subscription"."subscriptions" (
    "id" TEXT NOT NULL,
    "account_id" BIGINT NOT NULL,
    "plan_id" TEXT NOT NULL,
    "status" "subscription"."statuses" NOT NULL,
    "start_date" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "end_date" TIMESTAMPTZ(3),
    "cancel_at" TIMESTAMPTZ(3),

    CONSTRAINT "subscriptions_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "subscription"."invoices" (
    "id" TEXT NOT NULL,
    "subscription_id" TEXT NOT NULL,
    "amount" DOUBLE PRECISION NOT NULL,
    "issued_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "paid" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "invoices_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "resource"."resources" (
    "id" BIGSERIAL NOT NULL,
    "resource" TEXT NOT NULL,
    "type" "resource"."types" NOT NULL,
    "attempts" INTEGER NOT NULL DEFAULT 0,
    "failed" INTEGER NOT NULL DEFAULT 0,
    "disabled" BOOLEAN NOT NULL DEFAULT false,
    "last_used" TIMESTAMPTZ(3),

    CONSTRAINT "resources_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "telegrams_telegram_id_key" ON "account"."telegrams"("telegram_id");

-- CreateIndex
CREATE INDEX "telegrams_username_idx" ON "account"."telegrams"("username");

-- CreateIndex
CREATE INDEX "telegrams_telegram_id_idx" ON "account"."telegrams"("telegram_id");

-- CreateIndex
CREATE INDEX "telegrams_first_name_last_name_idx" ON "account"."telegrams"("first_name", "last_name");

-- CreateIndex
CREATE INDEX "telegrams_is_premium_idx" ON "account"."telegrams"("is_premium");

-- CreateIndex
CREATE INDEX "telegrams_created_at_idx" ON "account"."telegrams"("created_at");

-- CreateIndex
CREATE UNIQUE INDEX "plans_name_key" ON "subscription"."plans"("name");

-- CreateIndex
CREATE INDEX "plans_name_idx" ON "subscription"."plans"("name");

-- CreateIndex
CREATE INDEX "plans_price_idx" ON "subscription"."plans"("price");

-- CreateIndex
CREATE INDEX "plans_interval_idx" ON "subscription"."plans"("interval");

-- CreateIndex
CREATE INDEX "plans_description_idx" ON "subscription"."plans"("description");

-- CreateIndex
CREATE INDEX "subscriptions_account_id_idx" ON "subscription"."subscriptions"("account_id");

-- CreateIndex
CREATE INDEX "subscriptions_plan_id_idx" ON "subscription"."subscriptions"("plan_id");

-- CreateIndex
CREATE INDEX "subscriptions_status_idx" ON "subscription"."subscriptions"("status");

-- CreateIndex
CREATE INDEX "subscriptions_start_date_idx" ON "subscription"."subscriptions"("start_date");

-- CreateIndex
CREATE INDEX "subscriptions_end_date_idx" ON "subscription"."subscriptions"("end_date");

-- CreateIndex
CREATE INDEX "subscriptions_cancel_at_idx" ON "subscription"."subscriptions"("cancel_at");

-- CreateIndex
CREATE INDEX "invoices_subscription_id_idx" ON "subscription"."invoices"("subscription_id");

-- CreateIndex
CREATE INDEX "invoices_issued_at_idx" ON "subscription"."invoices"("issued_at");

-- CreateIndex
CREATE INDEX "invoices_paid_idx" ON "subscription"."invoices"("paid");

-- CreateIndex
CREATE UNIQUE INDEX "resources_resource_key" ON "resource"."resources"("resource");

-- CreateIndex
CREATE INDEX "resources_type_idx" ON "resource"."resources"("type");

-- CreateIndex
CREATE INDEX "resources_attempts_idx" ON "resource"."resources"("attempts");

-- CreateIndex
CREATE INDEX "resources_failed_idx" ON "resource"."resources"("failed");

-- CreateIndex
CREATE INDEX "resources_last_used_idx" ON "resource"."resources"("last_used");

-- AddForeignKey
ALTER TABLE "subscription"."subscriptions" ADD CONSTRAINT "subscriptions_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "account"."telegrams"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "subscription"."subscriptions" ADD CONSTRAINT "subscriptions_plan_id_fkey" FOREIGN KEY ("plan_id") REFERENCES "subscription"."plans"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "subscription"."invoices" ADD CONSTRAINT "invoices_subscription_id_fkey" FOREIGN KEY ("subscription_id") REFERENCES "subscription"."subscriptions"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

