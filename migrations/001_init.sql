CREATE TABLE leads (
    id BIGSERIAL PRIMARY KEY,

    name TEXT NOT NULL,
    email TEXT NOT NULL,
    company TEXT NOT NULL DEFAULT '',
    message TEXT NOT NULL DEFAULT '',

    budget NUMERIC(12, 2) NOT NULL DEFAULT 0,
    timeline TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_leads_email ON leads(email);
CREATE INDEX idx_leads_created_at ON leads(created_at);