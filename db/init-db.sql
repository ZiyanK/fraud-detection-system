CREATE TABLE "transactions" (
  "transaction_id" INTEGER PRIMARY KEY,
  "user_id" INTEGER NOT NULL,
  "amount" FLOAT NOT NULL,
  "type" VARCHAR(255) NULL,
  "is_fraud" BOOLEAN NOT NULL,
  "source" VARCHAR(255) NULL,
  "timestamp" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);