CREATE TYPE leave_type_enum AS ENUM ('annual', 'sick', 'unpaid');
CREATE TYPE leave_status_enum AS ENUM ('draft', 'waiting_approval', 'approved', 'rejected');

CREATE TABLE leave_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    leave_type leave_type_enum NOT NULL,
    status leave_status_enum NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);