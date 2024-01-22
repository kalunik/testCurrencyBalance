CREATE TABLE IF NOT EXISTS public.wallet(
    "id" SERIAL NOT NULL UNIQUE,
    "user_id" INT NOT NULL,
    "active_balance" DOUBLE PRECISION,
    "frozen_balance" DOUBLE PRECISION,
    "currency" CHAR(16) NOT NULL,

    CONSTRAINT "wallet_pk" PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS public.transactions (
    "id" SERIAL NOT NULL,
    "wallet_id" INTEGER NOT NULL,
    "status" CHAR(15) NOT NULL,
    "amount" DOUBLE PRECISION NOT NULL,
    "withdraw" BOOLEAN,
    "card_number" CHAR(127),
    "created_at" TIMESTAMP NOT NULL,

    CONSTRAINT "transactions_pk" PRIMARY KEY ("id")
);


ALTER TABLE "transactions" ADD CONSTRAINT "transactions_fk0" FOREIGN KEY ("wallet_id") REFERENCES "wallet"("id");
