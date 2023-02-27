CREATE TABLE "students" (
  "id" bigserial PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "date_of_birth" date NOT NULL,
  "grade" int NOT NULL,
  "phone" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "enrollments" (
  "id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "course_id" bigint NOT NULL,
  "enrollment_date" date NOT NULL DEFAULT CURRENT_DATE,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "enrollments" ("student_id");

CREATE INDEX ON "enrollments" ("course_id");

CREATE INDEX ON "enrollments" ("student_id", "course_id");

CREATE INDEX ON "enrollments" ("course_id", "student_id");

COMMENT ON COLUMN "students"."grade" IS 'must be between 0 and 11';

ALTER TABLE "enrollments" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("id");
