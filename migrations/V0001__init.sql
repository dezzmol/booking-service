-- Таблица Hotel
CREATE TABLE hotels
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Таблица Room
CREATE TABLE rooms
(
    id         BIGSERIAL PRIMARY KEY,
    number     BIGINT NOT NULL,
    type       INT8   NOT NULL,
    hotel_id   BIGINT NOT NULL REFERENCES hotels (id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Таблица Guest
CREATE TABLE guests
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Таблица Booking
CREATE TABLE bookings
(
    id                BIGSERIAL PRIMARY KEY,
    room_id           BIGINT  NOT NULL REFERENCES rooms (id),
    guest_id          BIGINT  NOT NULL REFERENCES guests (id),
    start_date        DATE    NOT NULL,
    end_date          DATE    NOT NULL,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW(),
    comment           TEXT,
    status            INT8    NOT NULL,
    is_paid           BOOLEAN NOT NULL
);

-- Таблица Review
CREATE TABLE reviews
(
    id         BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings (id),
    rating     BIGINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment    TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
