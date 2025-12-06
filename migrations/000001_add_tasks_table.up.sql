CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       user_id BIGINT NOT NULL,
                       title TEXT NOT NULL,
                       notes TEXT,
                       due_at TIMESTAMP WITH TIME ZONE,
                       completed BOOLEAN DEFAULT false,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
