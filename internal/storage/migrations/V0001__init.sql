-- Создание функции для обновления updated_at
CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Таблица Hotel
CREATE TABLE hotels
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trigger_hotel_update
    BEFORE UPDATE
    ON hotels
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Таблица Room
CREATE TABLE rooms
(
    id         BIGSERIAL PRIMARY KEY,
    number     BIGINT       NOT NULL,
    type       VARCHAR(100) NOT NULL,
    hotel_id   BIGINT       NOT NULL REFERENCES hotels (id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trigger_room_update
    BEFORE UPDATE
    ON rooms
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Таблица Guest
CREATE TABLE guests
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trigger_guest_update
    BEFORE UPDATE
    ON guests
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Таблица Employee
CREATE TABLE employees
(
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    role       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trigger_employee_update
    BEFORE UPDATE
    ON employees
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Таблица Booking
CREATE TABLE bookings
(
    id             BIGSERIAL PRIMARY KEY,
    room_id        BIGINT NOT NULL REFERENCES rooms (id),
    start_date     DATE   NOT NULL,
    end_date       DATE   NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW(),
    comment        TEXT,
    status         VARCHAR(100) CHECK ( status in ('pending', 'confirmed', 'cancelled', 'checked-in', 'checked-out')),
    payment_status VARCHAR(100) CHECK ( payment_status in ('paid', 'unpaid', 'cancelled'))
);

CREATE TRIGGER trigger_booking_update
    BEFORE UPDATE
    ON bookings
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Таблица bookings_guests
CREATE TABLE IF NOT EXISTS bookings_guests
(
    booking_id BIGINT NOT NULL REFERENCES bookings (id),
    guests_id  BIGINT NOT NULL REFERENCES guests (id)
);

-- Таблица Review
CREATE TABLE reviews
(
    id         BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings (id),
    guest_id   BIGINT NOT NULL REFERENCES guests (id),
    rating     BIGINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment    TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trigger_review_update
    BEFORE UPDATE
    ON reviews
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Таблица HousekeepingRequest
CREATE TABLE house_keeping_requests
(
    id           BIGSERIAL PRIMARY KEY,
    room_id      BIGINT      NOT NULL REFERENCES rooms (id),
    request_time TIMESTAMP   NOT NULL,
    status       VARCHAR(50) NOT NULL,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER trigger_housekeeping_request_update
    BEFORE UPDATE
    ON house_keeping_requests
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
