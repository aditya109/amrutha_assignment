create schema if not exists billing;

create table customers (
    id serial primary key,
    name varchar(200) not null,
    is_active boolean not null,
    address varchar(2000) not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    typ customer_type not null default 'DELINQUENT'
);

-- i am creating this for allowing further type loaning mechanism to be seamlessly adjusted 
create type interest_type add enum ('FIXED', 'VARIABLE', 'SIMPLE', 'COMPOUND') ;

create type customer_type add enum ('DELINQUENT', 'REGULAR');

create table loan_configs (
    id serial primary key,
    principal_amount decimal not null,
    max_span int not null,
    rate_of_interest decimal not null,
    type_of_loan interest_type not null default 'SIMPLE',
    is_active boolean not null
);

create table loans(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_config_id int not null,
    customer_id int not null,
    payment_completion_count int not null,
    missed_payment_count int not null,
    foreign key (loan_config_id) references loan_config(id),
    foreign key (customer_id) references customer(id)
);

create table payments(
    id serial primary key,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    loan_id int not null,
    customer_id int not null,
    amount decimal not null,
    is_accepted boolean not null default false,
    foreign key (loan_config_id) references loan_config(id),
    foreign key (customer_id) references customer(id)
);

create table billing_schedules(
    id serial primary key,
    is_paid boolean not null default false,
    loan_id int not null,
    payment_id int not null,
    is_late boolean not null true,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    foreign key (loan_id) references loan(id),
    foreign key (payment_id) references payment(id)
);

create view outstanding_payments_view as
select 
from 