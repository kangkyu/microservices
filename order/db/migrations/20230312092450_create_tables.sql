-- migrate:up
CREATE TABLE orders (
	id bigserial primary key,
	customer_id bigint,
	status text,
	created_at timestamp with time zone default now() not null
);
CREATE TABLE order_items (
	product_code text,
	unit_price decimal,
	quantity integer,
	order_id bigint not null
)

-- migrate:down
DROP TABLE order_items;
DROP TABLE orders;
