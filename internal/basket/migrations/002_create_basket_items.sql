CREATE TABLE IF NOT EXISTS basket_items (
    basket_id CHAR(26) NOT NULL REFERENCES baskets(id),
    item_id CHAR(26) NOT NULL REFERENCES items(id),
    quantity INT NOT NULL,
    price_per_item BIGINT NOT NULL,
    PRIMARY KEY (basket_id, item_id)
);
