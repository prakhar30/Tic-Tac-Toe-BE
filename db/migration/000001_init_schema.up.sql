CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL UNIQUE,
    "password_hash" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "games" (
    "id" bigserial PRIMARY KEY,
    "host_user_id" bigint,
    "status" varchar NOT NULL,
    "current_state" varchar,
    "next_turn_user_id" bigint
);

CREATE TABLE "game_participants" (
    "id" bigserial PRIMARY KEY,
    "game_id" bigint,
    "user_id" bigint
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "games" ADD FOREIGN KEY ("host_user_id") REFERENCES "users" ("id");
ALTER TABLE "games" ADD FOREIGN KEY ("next_turn_user_id") REFERENCES "users" ("id");
ALTER TABLE "game_participants" ADD FOREIGN KEY ("game_id") REFERENCES "games" ("id");
ALTER TABLE "game_participants" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
