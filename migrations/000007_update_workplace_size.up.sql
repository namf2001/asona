-- Update workplace_size enum to match UI requirements

-- 1. Create a temporary type with new values
DO $$ BEGIN
    CREATE TYPE workplace_size_new AS ENUM ('2-5', '6-10', '11-20', '21-50', '51-100', '101-250', '250-more');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- 2. Alter the table to use the new type
ALTER TABLE public.workplaces 
    ALTER COLUMN size TYPE workplace_size_new 
    USING (
        CASE 
            WHEN size::text = '1-10' THEN '6-10'::workplace_size_new
            WHEN size::text = '11-50' THEN '21-50'::workplace_size_new
            ELSE '2-5'::workplace_size_new -- Default fallback
        END
    );

-- 3. Drop the old type and rename the new type
DROP TYPE workplace_size;
ALTER TYPE workplace_size_new RENAME TO workplace_size;
