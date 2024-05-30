BEGIN;

-- Drop indexes
DROP INDEX IF EXISTS billing.idx_customer_id;
DROP INDEX IF EXISTS billing.idx_load_id;

-- Drop tables
DROP TABLE IF EXISTS billing.loan_account;
DROP TABLE IF EXISTS billing.billing_schedules;
DROP TABLE IF EXISTS billing.payments;
DROP TABLE IF EXISTS billing.loans;
DROP TABLE IF EXISTS billing.loan_configs;
DROP TABLE IF EXISTS billing.customers;

-- Drop types
DROP TYPE IF EXISTS billing.loan_state;
DROP TYPE IF EXISTS billing.interest_type;
DROP TYPE IF EXISTS billing.customer_type;

-- Drop schema
DROP SCHEMA IF EXISTS billing CASCADE;

COMMIT;
