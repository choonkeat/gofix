ALTER TABLE departments
  DROP CONSTRAINT IF EXISTS departments_company_id;
ALTER TABLE departments
  DROP CONSTRAINT IF EXISTS departments_supervisor_id;
ALTER TABLE employees
  DROP CONSTRAINT IF EXISTS employees_department_id;

DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS departments;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments_employees;
