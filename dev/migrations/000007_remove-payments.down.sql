BEGIN;

CREATE TABLE IF NOT EXISTS membership.payments
(
    id uuid UNIQUE DEFAULT gen_random_uuid(),
    date date NOT NULL,
    amount numeric NOT NULL,
    member_id uuid NOT NULL,
    CONSTRAINT unique_payments PRIMARY KEY (date, amount, member_id),
    CONSTRAINT member_payment FOREIGN KEY (member_id)
        REFERENCES membership.members (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

COMMIT;
