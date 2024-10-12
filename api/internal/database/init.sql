CREATE TABLE IF NOT EXISTS face_biometry(
    id       SERIAL PRIMARY KEY,
    filename VARCHAR,
    photo    BYTEA 
);

CREATE TABLE IF NOT EXISTS voice_biometry(
    id       SERIAL PRIMARY KEY,
    filename VARCHAR,
    audio    BYTEA
);
