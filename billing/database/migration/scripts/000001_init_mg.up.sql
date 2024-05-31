create schema if not exists billing;

do $$
begin
if not exists (select 1 from pg_type where typname = 'billing.customer_type') then
    create type billing.customer_type as enum ('DELINQUENT', 'REGULAR');
end if;

if not exists (select 1 from pg_type where typname = 'billing.interest_type') then
    -- i am creating this for allowing further type loaning mechanism to be seamlessly adjusted
    create type billing.interest_type as enum ('FIXED', 'VARIABLE', 'SIMPLE', 'COMPOUND') ;
end if;

if not exists (select 1 from pg_type where typname = 'billing.loan_state') then
    create type billing.loan_state as enum ('ACTIVE', 'PAID', 'INACTIVE');
end if;
end $$;

create table if not exists billing.customers (
    id serial primary key,
    name varchar(200) not null,
    is_active boolean not null,
    address varchar(2000) not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    typ billing.customer_type not null default 'DELINQUENT'
);


create table if not exists billing.loan_configs (
    id serial primary key,
    principal_amount varchar not null,
    max_span int not null,
    rate_of_interest varchar not null,
    type_of_loan billing.interest_type not null default 'SIMPLE',
    is_active boolean not null
);



create table if not exists billing.loans(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_config_id int not null,
    customer_id int not null,
    payment_completion_count int not null,
    missed_payment_count int not null,
    loan_state billing.loan_state not null default 'INACTIVE',
    foreign key (loan_config_id) references billing.loan_configs(id),
    foreign key (customer_id) references billing.customers(id)
);

create table if not exists billing.payments(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_account_id int not null,
    customer_id int not null,
    amount varchar not null,
    is_accepted boolean not null default false,
    foreign key (loan_account_id) references billing.loan_accounts(id),
    foreign key (customer_id) references billing.customers(id)
);

create table if not exists billing.loan_accounts(
   id serial primary key,
   created_at timestamp default now(),
   updated_at timestamp default now(),
   loan_id int not null,
   payable_principal_amount varchar not null,
   accrued_interest varchar not null,
   total_payable_amount varchar not null,
   total_paid_amount varchar not null,
   outstanding_amount varchar not null,
   installment_amount varchar not null,
   foreign key (loan_id) references billing.loans(id)
);

create table if not exists billing.billing_schedules(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_account_id int not null,
    installment_amount varchar not null,
    start_date timestamp not null,
    end_date timestamp not null,
    week_count int not null,
    foreign key (loan_account_id) references billing.loan_accounts(id)
);



create index if not exists idx_load_id on billing.loan_accounts(loan_id);
create index if not exists idx_loan_account_id on billing.billing_schedules(loan_account_id);
create index if not exists idx_load_id on billing.payments(loan_id);
create index if not exists idx_customer_id on billing.payments(customer_id);
create index if not exists idx_customer_id on billing.loans(customer_id);

alter table billing.customers
add column if not exists display_id varchar(20) not null;

alter table billing.loans
add column if not exists display_id varchar(20) not null;

alter table billing.loan_accounts
add column if not exists installment_amount varchar not null;

alter table billing.loan_accounts
add column if not exists display_id varchar(20) not null;


ALTER TYPE billing.interest_type ADD VALUE 'FLAT';

INSERT INTO billing.loan_configs
(id, principal_amount, max_span, rate_of_interest, type_of_loan, is_active)
VALUES(nextval('billing.loan_configs_id_seq'::regclass), '5000000.00', 50, '0.1', 'FLAT'::billing.interest_type, true);

