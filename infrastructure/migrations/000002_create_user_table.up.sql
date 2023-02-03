CREATE EXTENSION citext;

CREATE OR REPLACE FUNCTION update_modified_column() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';


CREATE TABLE IF NOT EXISTS service.user(
    id uuid DEFAULT uuid_generate_v4(),
    email citext UNIQUE NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    CONSTRAINT pk_user PRIMARY KEY (id)
);

CREATE TRIGGER update_user_updated_at BEFORE UPDATE ON service.user FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();