create table warehouse_products
(
    productsId integer,
    amount integer
) partition by hash (productsId);

create table order_active_orders
(
    id integer,
    userId integer,
    paymentsMethod text,
    paymentsData text,
    productsId integer,
    amount integer
) partition by hash (userId);
create extension postgres_fdw;