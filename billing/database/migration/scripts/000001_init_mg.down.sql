-- Down migration script

-- Drop foreign key constraints
alter table billing.payments drop constraint if exists payments_loan_id_fkey;
alter table billing.payments drop constraint if exists payments_customer_id_fkey;
alter table billing.loan_account drop constraint if exists loan_account_loan_id_fkey;
alter table billing.billing_schedules drop constraint if exists billing_schedules_loan_account_id_fkey;
alter table billing.loans drop constraint if exists loans_loan_config_id_fkey;
alter table billing.loans drop constraint if exists loans_customer_id_fkey;

-- Drop indexes
drop index if exists billing.idx_load_id;
drop index if exists billing.idx_loan_account_id;
drop index if exists billing.idx_customer_id;

-- Drop tables
drop table if exists billing.billing_schedules;
drop table if exists billing.loan_account;
drop table if exists billing.payments;
drop table if exists billing.loans;
drop table if exists billing.loan_configs;
drop table if exists billing.customers;

-- Drop types
drop type if exists billing.loan_state;
drop type if exists billing.interest_type;
drop type if exists billing.customer_type;

-- Drop schema
drop schema if exists billing cascade;
