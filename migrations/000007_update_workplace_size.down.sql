-- Revert workplace_size enum

DO $$ BEGIN
    CREATE TYPE workplace_size_old AS ENUM ('1-10', '11-50', '51-200', '201-500', '500+');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

ALTER TABLE public.workplaces 
    ALTER COLUMN size TYPE workplace_size_old 
    USING (
        CASE 
            WHEN size::text = '6-10' THEN '1-10'::workplace_size_old
            WHEN size::text = '21-50' THEN '11-50'::workplace_size_old
            ELSE '1-10'::workplace_size_old
        END
    );

DROP TYPE workplace_size;
ALTER TYPE workplace_size_old RENAME TO workplace_size;
