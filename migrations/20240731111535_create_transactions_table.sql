-- +goose Up
-- +goose StatementBegin
CREATE TABLE "transactions" (
  "transaction_id" INTEGER PRIMARY KEY,
  "user_id" INTEGER UNIQUE NOT NULL,
  "amount" FLOAT NOT NULL,
  "type" VARCHAR(255) NULL,
  "is_fraud" BOOLEAN NOT NULL,
  "timestamp" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "transactions";
-- +goose StatementEnd
