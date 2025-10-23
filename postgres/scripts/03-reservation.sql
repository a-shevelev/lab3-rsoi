\c reservations

CREATE TABLE IF NOT EXISTS reservation
(
    id              SERIAL PRIMARY KEY,
    reservation_uid UUID UNIQUE NOT NULL,
    username        VARCHAR(80) NOT NULL,
    book_uid        UUID NOT NULL,
    library_uid     UUID NOT NULL,
    status          VARCHAR(20) NOT NULL
    CHECK (status IN ('RENTED', 'RETURNED', 'EXPIRED')),
    start_date      TIMESTAMP NOT NULL,
    till_date       TIMESTAMP NOT NULL
    );
