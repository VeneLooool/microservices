create table if not exists warehouse_products_1
(
    productsId integer,
    amount integer
);

create table if not exists order_active_orders_1
(
    id integer,
    userId integer,
    paymentsMethod text,
    paymentsData text,
    productsId integer,
    amount integer
);