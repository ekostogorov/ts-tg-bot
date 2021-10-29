-- students
CREATE SEQUENCE IF NOT EXISTS students_id_seq;

CREATE TABLE "public"."students" (
    "id" int4 NOT NULL DEFAULT nextval('students_id_seq'::regclass),
    "name" text NOT NULL DEFAULT ''''''::text,
    "login" text NOT NULL DEFAULT ''''''::text,
    "folder" text NOT NULL DEFAULT ''''''::text,
    "telegram_user_id" text,
    "is_activated" bool NOT NULL DEFAULT false,
    "activated_at" timestamptz,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);


-- lectures
CREATE SEQUENCE IF NOT EXISTS lectures_id_seq;

CREATE TABLE "public"."lectures" (
    "id" int4 NOT NULL DEFAULT nextval('lectures_id_seq'::regclass),
    "name" text NOT NULL DEFAULT ''''''::text,
    "file_path" text NOT NULL DEFAULT ''''''::text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);


-- homeworks
CREATE SEQUENCE IF NOT EXISTS homeworks_id_seq;

CREATE TABLE "public"."homeworks" (
    "id" int4 NOT NULL DEFAULT nextval('homeworks_id_seq'::regclass),
    "name" text NOT NULL DEFAULT ''''''::text,
    "lecture_id" int8 NOT NULL,
    "file_path" text NOT NULL DEFAULT ''''''::text,
    "is_active" bool NOT NULL DEFAULT false,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT "homeworks_lecture_id_fkey" FOREIGN KEY ("lecture_id") REFERENCES "public"."lectures"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);


-- student homeworks
CREATE TABLE "public"."student_homeworks" (
    "student_id" int8 NOT NULL,
    "homework_id" int8 NOT NULL,
    "passed_at" timestamptz NOT NULL DEFAULT now(),
    "file_path" text NOT NULL DEFAULT ''''''::text,
    "grade" int2 NOT NULL DEFAULT 0,
    CONSTRAINT "student_homeworks_homework_id_fkey" FOREIGN KEY ("homework_id") REFERENCES "public"."homeworks"("id") ON DELETE CASCADE,
    CONSTRAINT "student_homeworks_student_id_fkey" FOREIGN KEY ("student_id") REFERENCES "public"."students"("id") ON DELETE CASCADE
);
