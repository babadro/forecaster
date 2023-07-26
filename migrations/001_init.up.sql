CREATE TABLE forecaster.series (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE TABLE forecaster.polls (
    id SERIAL PRIMARY KEY,
    series_id INT NOT NULL DEFAULT 0,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start timestamp with time zone NOT NULL,
    finish timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    FOREIGN KEY (series_id) REFERENCES forecaster.series(id) ON DELETE CASCADE
);

CREATE TABLE forecaster.options (
    id SERIAL PRIMARY KEY,
    poll_id INT,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY (poll_id) REFERENCES forecaster.polls(id) ON DELETE CASCADE
);
