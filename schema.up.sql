CREATE TABLE companies (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  name text,
  address text,
  created_at timestamp with time zone NOT NULL);

CREATE TABLE departments (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  company_id bigint NOT NULL,
  name text,
  supervisor_id bigint,
  created_at timestamp with time zone NOT NULL);

CREATE TABLE employees (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  department_id bigint NOT NULL,
  name text NOT NULL,
  email text NOT NULL,
  created_at timestamp with time zone NOT NULL);

ALTER TABLE departments
  ADD CONSTRAINT departments_company_id FOREIGN KEY(company_id) REFERENCES companies(id);

ALTER TABLE departments
  ADD CONSTRAINT departments_supervisor_id FOREIGN KEY(supervisor_id) REFERENCES employees(id);

ALTER TABLE employees
  ADD CONSTRAINT employees_department_id FOREIGN KEY(department_id) REFERENCES departments(id);
