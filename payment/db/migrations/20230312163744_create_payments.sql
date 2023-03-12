-- migrate:up
CREATE TABLE payments (
	id bigserial primary key,
	customer_id bigint,
	status text,
	order_id bigint,
	total_price decimal,
	created_at timestamp with time zone default now() not null
);

-- migrate:down
DROP TABLE payments;
