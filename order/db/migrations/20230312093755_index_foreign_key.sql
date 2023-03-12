-- migrate:up transaction:false
CREATE INDEX index_order_items_on_order_id ON order_items (order_id);

-- migrate:down transaction:false
DROP INDEX index_order_items_on_order_id;
