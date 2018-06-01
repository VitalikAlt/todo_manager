CREATE TABLE tasks (
  id serial UNIQUE NOT NULL,
  user_id int NOT NULL,
  "order" int NOT NULL,
  text varchar NOT NULL,
  priority varchar NOT NULL,
  due_date timestamp,
  created_at timestamp NOT NULL DEFAULT now(),
  finished_at timestamp
);

REVOKE ALL ON TABLE tasks FROM PUBLIC;
ALTER TABLE public.tasks OWNER TO SESSION_USER;
ALTER TABLE ONLY tasks ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


ALTER TABLE ONLY tasks ADD CONSTRAINT task_user_id_fkey FOREIGN KEY (user_id)
      REFERENCES public.users (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE;