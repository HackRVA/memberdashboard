CREATE EXTENSION "pgcrypto";

CREATE EXTENSION IF NOT EXISTS citext;

--
-- Name: membership; Type: SCHEMA; Schema: -; Owner: test
--

CREATE SCHEMA membership;

--
-- Name: member_resource; Type: TABLE; Schema: membership; Owner: test
--

CREATE TABLE membership.member_resource (
    id UUID DEFAULT gen_random_uuid(),
    member_id UUID NOT NULL,
    resource_id UUID NOT NULL
);

ALTER TABLE membership.member_resource
    ADD CONSTRAINT unique_relationship UNIQUE (member_id, resource_id)
    INCLUDE (member_id, resource_id);

CREATE TABLE membership.users
(
    email    text NOT NULL,
    password text NOT NULL,
    PRIMARY KEY (email)
);


--
-- Name: member_tiers; Type: TABLE; Schema: membership; Owner: test
--

CREATE TABLE membership.member_tiers (
    id integer NOT NULL,
    description text NOT NULL,
    price integer NOT NULL
);

--
-- Name: member_tiers_id_seq; Type: SEQUENCE; Schema: membership; Owner: test
--

CREATE SEQUENCE membership.member_tiers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: member_tiers_id_seq; Type: SEQUENCE OWNED BY; Schema: membership; Owner: test
--

ALTER SEQUENCE membership.member_tiers_id_seq OWNED BY membership.member_tiers.id;


--
-- Name: members; Type: TABLE; Schema: membership; Owner: test
--

CREATE TABLE membership.members (
    id UUID UNIQUE DEFAULT gen_random_uuid(),
    name text NOT NULL,
    email citext NOT NULL,
    rfid text,
    member_tier_id integer NOT NULL
);


ALTER TABLE membership.members
    ADD CONSTRAINT unique_email UNIQUE (email);

ALTER TABLE membership.members
    ADD CONSTRAINT unique_rfid UNIQUE (rfid);


CREATE TABLE membership.payments
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

--
-- Name: resources; Type: TABLE; Schema: membership; Owner: test
--

CREATE TABLE membership.resources (
    id UUID DEFAULT gen_random_uuid(),
    description text NOT NULL,
    device_identifier text NOT NULL,
    is_default boolean NOT NULL
);

--
-- Name: member_resource member_resource_pkey; Type: CONSTRAINT; Schema: membership; Owner: test
--

ALTER TABLE ONLY membership.member_resource
    ADD CONSTRAINT member_resource_pkey PRIMARY KEY (id);


--
-- Name: member_tiers member_tiers_pkey; Type: CONSTRAINT; Schema: membership; Owner: test
--

ALTER TABLE ONLY membership.member_tiers
    ADD CONSTRAINT member_tiers_pkey PRIMARY KEY (id);


--
-- Name: members members_pkey; Type: CONSTRAINT; Schema: membership; Owner: test
--

ALTER TABLE ONLY membership.members
    ADD CONSTRAINT members_pkey PRIMARY KEY (id);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: membership; Owner: test
--

ALTER TABLE ONLY membership.resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id);


--
-- Name: resources unique_device_id; Type: CONSTRAINT; Schema: membership; Owner: test
--

ALTER TABLE ONLY membership.resources
    ADD CONSTRAINT unique_device_id UNIQUE (device_identifier);


--
-- Name: member_resource member; Type: FK CONSTRAINT; Schema: membership; Owner: test
--

ALTER TABLE ONLY membership.member_resource
    ADD CONSTRAINT member FOREIGN KEY (member_id) REFERENCES membership.members(id);


--
-- Name: members member_tier_id; Type: FK CONSTRAINT; Schema: membership; Owner: test
--

ALTER TABLE ONLY membership.members
    ADD CONSTRAINT member_tier_id FOREIGN KEY (member_tier_id) REFERENCES membership.member_tiers(id) NOT VALID;


--
-- Name: member_resource resource; Type: FK CONSTRAINT; Schema: membership; Owner: test
--

ALTER TABLE ONLY membership.member_resource
    ADD CONSTRAINT resource FOREIGN KEY (resource_id) REFERENCES membership.resources(id) ON DELETE CASCADE;
