CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE EXTENSION IF NOT EXISTS citext;

CREATE SCHEMA IF NOT EXISTS membership;

CREATE TABLE IF NOT EXISTS membership.users
(
    email    text NOT NULL,
    password text NOT NULL,
    PRIMARY KEY (email)
);

CREATE TABLE IF NOT EXISTS membership.member_tiers (
    id integer NOT NULL PRIMARY KEY,
    description text NOT NULL,
    price integer NOT NULL
);

CREATE SEQUENCE IF NOT EXISTS membership.member_tiers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
    OWNED BY membership.member_tiers.id;

CREATE TABLE IF NOT EXISTS membership.members (
    id UUID UNIQUE DEFAULT gen_random_uuid() PRIMARY KEY,
    name text NOT NULL,
    email citext NOT NULL UNIQUE,
    rfid text UNIQUE,
    member_tier_id integer NOT NULL,
    CONSTRAINT member_tier_id
        FOREIGN KEY (member_tier_id)
            REFERENCES membership.member_tiers(id) NOT VALID
);

CREATE TABLE IF NOT EXISTS membership.member_credit
(
    member_id uuid NOT NULL,
    CONSTRAINT member_credit FOREIGN KEY (member_id)
        REFERENCES membership.members (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

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

CREATE TABLE IF NOT EXISTS membership.resources (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    description text NOT NULL,
    device_identifier text NOT NULL,
    is_default boolean NOT NULL,
    CONSTRAINT unique_device_id UNIQUE (device_identifier)
);

CREATE TABLE IF NOT EXISTS membership.member_resource (
    id UUID DEFAULT gen_random_uuid(),
    member_id UUID NOT NULL,
    resource_id UUID NOT NULL,
    CONSTRAINT unique_relationship
        UNIQUE (member_id, resource_id)
            INCLUDE (member_id, resource_id),
    CONSTRAINT member
        FOREIGN KEY (member_id)
            REFERENCES membership.members(id),
    CONSTRAINT resource
        FOREIGN KEY (resource_id)
            REFERENCES membership.resources(id)
            ON DELETE CASCADE
);
