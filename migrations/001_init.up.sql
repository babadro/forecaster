CREATE TABLE forecaster.polls (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start timestamp with time zone NOT NULL,
    finish timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
)

CREATE TABLE forecaster.options {
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    poll_id INT REFERENCES Poll(id)
}