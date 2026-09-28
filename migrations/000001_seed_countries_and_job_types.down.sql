
DELETE FROM locations;

DELETE FROM countries 
WHERE country_name IN ('Australia', 'Canada', 'India', 'United Kingdom', 'United States');

DELETE FROM job_types 
WHERE type_name IN ('Full-Time', 'Freelance', 'Contract', 'TeleCom');