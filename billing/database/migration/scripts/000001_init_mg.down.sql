-- Migration down script

-- Drop the inserted row in billing.loan_configs if it exists
DELETE FROM billing.loan_configs WHERE type_of_loan = 'FLAT'::billing.interest_type;

-- Drop columns added in the alter table statements if they exist
ALTER TABLE billing.loan_accounts DROP COLUMN IF EXISTS display_id;
ALTER TABLE billing.loan_accounts DROP COLUMN IF EXISTS installment_amount;
ALTER TABLE billing.loans DROP COLUMN IF EXISTS display_id;
ALTER TABLE billing.customers DROP COLUMN IF EXISTS display_id;

-- Drop indexes if they exist
DROP INDEX IF EXISTS idx_customer_id;
DROP INDEX IF EXISTS idx_loan_id;
DROP INDEX IF EXISTS idx_loan_account_id;
DROP INDEX IF EXISTS idx_load_id;

-- Drop tables if they exist
DROP TABLE IF EXISTS billing.payments;
DROP TABLE IF EXISTS billing.billing_schedules;
DROP TABLE IF EXISTS billing.loan_accounts;
DROP TABLE IF EXISTS billing.loans;
DROP TABLE IF EXISTS billing.loan_configs;
DROP TABLE IF EXISTS billing.customers;

-- Drop types if they exist
DROP TYPE billing.interest_type;
DROP TYPE billing.customer_type;
DROP TYPE billing.loan_state;

-- Drop schema if it exists
DROP SCHEMA IF EXISTS billing;
