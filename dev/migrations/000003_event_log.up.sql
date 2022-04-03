
CREATE TABLE IF NOT EXISTS membership.event_log
(
    id UUID UNIQUE DEFAULT gen_random_uuid() PRIMARY KEY,
    type text NOT NULL,
    event_time timestamp without time zone NOT NULL,
    is_known boolean NOT NULL,
    username text NOT NULL,
    rfid text NOT NULL,
    door text NOT NULL
);
