SET TIME ZONE 'America/Sao_Paulo';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE simple_storage (
                                variable_value INTEGER NOT NULL,
                                last_synced_at TIMESTAMP DEFAULT NOW()
);
insert into simple_storage (variable_value) values (0);



