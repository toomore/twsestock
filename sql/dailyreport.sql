CREATE TABLE IF NOT EXISTS dailyreport (
    no VARCHAR(10) NOT NULL,
    filter TINYINT UNSIGNED,
    timestamp TIMESTAMP NOT NULL,
    PRIMARY KEY(no, filter, timestamp)
    )
CHARACTER SET 'utf8';
