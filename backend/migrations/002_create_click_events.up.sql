CREATE TABLE click_events (
    id         BIGSERIAL PRIMARY KEY,
    url_id     BIGINT NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    referer    TEXT,
    country    VARCHAR(2)
);

CREATE INDEX idx_click_events_url_id ON click_events (url_id);
CREATE INDEX idx_click_events_url_id_clicked_at ON click_events (url_id, clicked_at);
