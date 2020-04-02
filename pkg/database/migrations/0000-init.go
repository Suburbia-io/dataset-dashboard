package migrations

// language=SQL
const CreateAdminsTable = `
CREATE TABLE admins (
  admin_id       UUID        NOT NULL PRIMARY KEY,
  super          BOOL        NOT NULL DEFAULT FALSE,
  email          TEXT        NOT NULL UNIQUE,
  hash           TEXT        NOT NULL,
  name           TEXT        NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL,
  archived_at    TIMESTAMPTZ,
  last_active_at TIMESTAMPTZ
);
`

// language=SQL
const CreateAuthTokensTable = `
CREATE TABLE auth_tokens(
  token      TEXT        NOT NULL PRIMARY KEY,
  entity_id  UUID        NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL
);
`

// language=SQL
const CreateAdminSessionsTable = `
CREATE TABLE admin_session_tokens(
  token      TEXT        NOT NULL PRIMARY KEY,
  admin_id   UUID        NOT NULL REFERENCES admins (admin_id),
  expires_at TIMESTAMPTZ NOT NULL
);
`

// language=SQL
const CreateUsersTable = `
CREATE TABLE users (
  user_id      UUID        NOT NULL PRIMARY KEY,
  email        TEXT        NOT NULL,
  name         TEXT        NOT NULL,
  
  created_at     TIMESTAMPTZ NOT NULL,
  suspended_at   TIMESTAMPTZ,
  archived_at    TIMESTAMPTZ,
  last_active_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX users_lower_email ON users (lower(email));
CREATE INDEX users_lower_name         ON users (lower(name));
CREATE INDEX users_activity_index     ON users (archived_at,last_active_at,suspended_at);
CREATE INDEX users_created_index      ON users (archived_at,created_at,suspended_at);
`

// language=SQL
const CreateBrandsTable = `
CREATE TABLE brands(
  brand_id UUID PRIMARY KEY,
  
  slug        TEXT NOT NULL,
  label       TEXT NOT NULL,
  description TEXT NOT NULL,
  variant     TEXT NOT NULL DEFAULT '',
  
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  archived_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX brands_slug_variant_key ON brands (slug, variant);
`

// language=SQL
const CreateIXRulesTable = `
CREATE TABLE ixrules(
  rule_id UUID PRIMARY KEY,
  
  tag_group TEXT NOT NULL,
  tag_id    UUID NOT NULL,
  
  includes  TEXT[] NOT NULL DEFAULT '{}',
  excludes  TEXT[] NOT NULL DEFAULT '{}',
  
  created_at  TIMESTAMPTZ NOT NULL,
  updated_at  TIMESTAMPTZ NOT NULL,
  archived_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ixrules_unique_group_rule on ixrules(tag_group,includes,excludes);
CREATE INDEX ixrules_tags on ixrules(tag_group, tag_id);
`

// language=SQL
const CreateUserSessionsTable = `
CREATE TABLE user_session_tokens(
  token TEXT NOT NULL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (user_id),
  expires_at TIMESTAMPTZ NOT NULL
);
`

// language=SQL
const CreateTickersTable = `
CREATE TABLE tickers(
 ticker_id UUID PRIMARY KEY,

 exchange TEXT NOT NULL,
 code     TEXT NOT NULL,
 name     TEXT NOT NULL,

 created_at  TIMESTAMPTZ NOT NULL,
 updated_at  TIMESTAMPTZ NOT NULL,
 archived_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX tickers_symbol ON tickers(exchange, code);
CREATE INDEX tickers_name ON tickers(lower(name));
`

// language=SQL
const CreateCorporationsTable = `
CREATE TABLE corporations (
corporation_id UUID PRIMARY KEY,

slug text NOT NULL,
name text NOT NULL,

created_at  TIMESTAMPTZ NOT NULL,
updated_at  TIMESTAMPTZ NOT NULL,
archived_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX corporation_slug ON corporations(slug);
CREATE INDEX corporation_name ON corporations(lower(name));
`

// language=SQl
const CreateEventsTable = `
CREATE TABLE events (
  event_id  BIGSERIAL PRIMARY KEY,
  timestamp TIMESTAMPTZ NOT NULL,
  -- subject
  by_system BOOLEAN DEFAULT FALSE,
  by_admin  UUID REFERENCES admins (admin_id),
  by_user   UUID REFERENCES users (user_id),
  -- object
  on_admin       UUID REFERENCES admins (admin_id),
  on_user        UUID REFERENCES users (user_id),
  on_brand       UUID REFERENCES brands (brand_id),
  on_ixrule      UUID REFERENCES ixrules (rule_id),
  on_ticker      UUID REFERENCES tickers (ticker_id),
  on_corporation UUID REFERENCES corporations (corporation_id),
  -- verb
  name      TEXT NOT NULL,
  payload   JSON
);
CREATE INDEX events_name_index
  ON events (name);
CREATE INDEX events_by_admin_index
  ON events (by_admin,name);
CREATE INDEX events_by_user_index
  ON events (by_user,name);
CREATE INDEX events_on_admin_index
  ON events (on_admin,name);
CREATE INDEX events_on_user_index
  ON events (on_user,name);
CREATE INDEX events_on_brand_index
  ON events (on_brand,name);
CREATE INDEX events_on_ixrule_index
  ON events (on_ixrule,name);
CREATE INDEX events_on_ticker_index
  ON events (on_ticker,name);
CREATE INDEX events_on_corporation_index
  ON events (on_corporation,name);
`
