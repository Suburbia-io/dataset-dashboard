package migrations

// language=SQL
const TagJobsAndTasks = `
CREATE TABLE tag_jobs(
  tag_job_id  UUID NOT NULL PRIMARY KEY,
  dataset_id  UUID NOT NULL,
  fingerprint TEXT NOT NULL,
  tag_type_id UUID NOT NULL,
  tag_level   INT  NOT NULL,

  state  TEXT NOT NULL, -- 'open', 'waiting', 'complete', 'closed'

  tasks_completed   INT NOT NULL,
  tasks_outstanding INT NOT NULL,

  FOREIGN KEY(dataset_id, fingerprint)
    REFERENCES fingerprints(dataset_id, fingerprint),

  FOREIGN KEY(dataset_id, tag_type_id)
    REFERENCES tag_types(dataset_id, tag_type_id)
);

-- Find new or completed jobs.
CREATE INDEX ON tag_jobs(state);

-- List completed jobs for a given dataset (order by fp).
CREATE INDEX ON tag_jobs(dataset_id, state, fingerprint);

CREATE TABLE tag_tasks(
  tag_task_id  UUID NOT NULL UNIQUE,
  tag_job_id   UUID NOT NULL REFERENCES tag_jobs(tag_job_id),
  customer_id  UUID NOT NULL,
  user_id      UUID NOT NULL,
  custom_id    TEXT NOT NULL,

  state TEXT NOT NULL, -- 'open', 'closed', 'complete'

  created_at   TIMESTAMPTZ NOT NULL,
  completed_at TIMESTAMPTZ,

  tag_id       UUID,
  confidence   REAL NOT NULL, -- [0, 1]

  FOREIGN KEY(customer_id, user_id)
    REFERENCES users(customer_id, user_id),

  PRIMARY KEY(tag_job_id, customer_id, user_id, custom_id)
);

CREATE INDEX ON tag_tasks(state, created_at);

CREATE INDEX ON tag_tasks(customer_id, user_id, custom_id, state);
`
