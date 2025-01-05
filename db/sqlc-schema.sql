CREATE TABLE IF NOT EXISTS Entries (
    Date TEXT NOT NULL CONSTRAINT "PK_Entries" PRIMARY KEY,
    Content TEXT,
    Keyword TEXT,
    Mood TEXT,
    Remarks TEXT
);