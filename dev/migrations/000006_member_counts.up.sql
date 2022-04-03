CREATE TABLE IF NOT EXISTS membership.member_counts
(
    month date NOT NULL,
    classic integer NOT NULL,
    standard integer NOT NULL,
    premium integer NOT NULL,
    credited integer NOT NULL,
    PRIMARY KEY (month)
);
