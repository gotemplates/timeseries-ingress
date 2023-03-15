DROP TABLE IF EXISTS access_log;

CREATE TABLE access_log (
                            customer_id TEXT NOT NULL,
                            start_time TIMESTAMPTZ DEFAULT now() NOT NULL,
                            duration_ms BIGINT NOT NULL,
                            duration_str TEXT,
                            traffic TEXT NOT NULL,

                            region TEXT,
                            zone TEXT,
                            sub_zone TEXT,
                            service TEXT NOT NULL,
                            instance_id TEXT NOT NULL,
                            route_name TEXT,

                            request_id TEXT,
                            url TEXT,
                            protocol TEXT,
                            method TEXT,
                            host TEXT,
                            path TEXT,

                            status_code INT,
                            bytes_sent BIGINT,
                            status_flags TEXT,

                            timeout INT,
                            rate_limit DOUBLE PRECISION,
                            rate_burst INT,
                            retry BOOLEAN,
                            retry_rate_limit DOUBLE PRECISION,
                            retry_rate_burst INT,
                            failover BOOLEAN
);

SELECT create_hypertable('access_log', 'start_time');
