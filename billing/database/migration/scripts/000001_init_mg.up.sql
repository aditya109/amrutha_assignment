create schema if not exists billing;

create type billing.customer_type as enum ('DELINQUENT', 'REGULAR');

create table if not exists billing.customers (
    id serial primary key,
    name varchar(200) not null,
    is_active boolean not null,
    address varchar(2000) not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    typ billing.customer_type not null default 'DELINQUENT'
);

-- i am creating this for allowing further type loaning mechanism to be seamlessly adjusted 
create type billing.interest_type as enum ('FIXED', 'VARIABLE', 'SIMPLE', 'COMPOUND') ;

create table if not exists billing.loan_configs (
    id serial primary key,
    principal_amount decimal not null,
    max_span int not null,
    rate_of_interest decimal not null,
    type_of_loan billing.interest_type not null default 'SIMPLE',
    is_active boolean not null
);

create type billing.loan_state as enum ('ACTIVE', 'PAID', 'INACTIVE');

create table if not exists billing.loans(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_config_id int not null,
    customer_id int not null,
    payment_completion_count int not null,
    missed_payment_count int not null,
    load_state billing.loan_state not null default 'INACTIVE',
    foreign key (loan_config_id) references billing.loan_configs(id),
    foreign key (customer_id) references billing.customers(id)
);

create table if not exists billing.payments(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_id int not null,
    customer_id int not null,
    amount decimal not null,
    is_accepted boolean not null default false,
    foreign key (loan_id) references billing.loans(id),
    foreign key (customer_id) references billing.customers(id)
);

create table if not exists billing.billing_schedules(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_id int not null,
    start_date timestamp not null,
    end_date timestamp not null,
    week_count int not null,
    foreign key (loan_id) references billing.loans(id)
);

create table if not exists billing.loan_account(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_id int not null,
    payable_principal_amount decimal not null,
    accrued_interest decimal not null,
    total_payable_amount decimal not null,
    total_paid_amount decimal not null,
    outstanding_amount decimal not null,
    foreign key (loan_id) references billing.loans(id)
);

create index if not exists idx_load_id on billing.loan_account(loan_id);
create index if not exists idx_load_id on billing.billing_schedules(loan_id);
create index if not exists idx_load_id on billing.payments(loan_id);
create index if not exists idx_customer_id on billing.payments(customer_id);
create index if not exists idx_customer_id on billing.loans(customer_id);


