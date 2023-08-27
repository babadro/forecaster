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
    id SMALLINT NOT NULL,
    poll_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    is_actual_outcome BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (poll_id) REFERENCES forecaster.polls(id) ON DELETE CASCADE,
    PRIMARY KEY (poll_id, id)
);

CREATE TABLE forecaster.votes (
    poll_id INT NOT NULL,
    user_id BIGINT NOT NULL,
    option_id INT NOT NULL,
    epoch_unix_timestamp BIGINT NOT NULL,
    FOREIGN KEY (poll_id, option_id) REFERENCES forecaster.options(poll_id, id) ON DELETE CASCADE,
    PRIMARY KEY (poll_id, user_id, epoch_unix_timestamp)
);
