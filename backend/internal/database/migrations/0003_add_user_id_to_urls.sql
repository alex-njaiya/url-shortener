-- migration_add_user_id_to_urls.sql

-- 1. Add the column (allow NULL initially)
ALTER TABLE urls 
ADD COLUMN user_id BIGINT;

-- 2. Add foreign key constraint
ALTER TABLE urls 
ADD CONSTRAINT fk_urls_user 
FOREIGN KEY (user_id) 
REFERENCES users(id) 
ON DELETE SET NULL;

-- 3. Create index
CREATE INDEX idx_urls_user_id ON urls(user_id);