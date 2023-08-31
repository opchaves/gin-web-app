BEGIN;

SET timezone TO 'GMT';

CREATE TABLE IF NOT EXISTS users(
  "id" UUID NOT NULL,
  "first_name" VARCHAR NOT NULL,
  "last_name" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "password" VARCHAR NOT NULL,
  "role" VARCHAR NOT NULL,
  "last_login" TIMESTAMP WITHOUT TIME ZONE NULL,
  "active" BOOLEAN NOT NULL DEFAULT true,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "uq_users_email" UNIQUE ("email"),
  CONSTRAINT "pk_users_id" PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS workspaces(
  "id" UUID NOT NULL,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR,
  "currency" VARCHAR NOT NULL,
  "language" VARCHAR NOT NULL,
  "user_id" UUID NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "pk_workspaces_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_workspaces_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS categories(
  "id" UUID NOT NULL,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR,
  "c_type" VARCHAR NOT NULL,
  "user_id" UUID NOT NULL,
  "workspace_id" UUID NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "pk_categories_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "fk_categories_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES "workspaces"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS accounts(
  "id" UUID NOT NULL,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR,
  "balance" NUMERIC(10,2),
  "financial_institution" VARCHAR,
  "account_type" VARCHAR,
  "user_id" UUID NOT NULL,
  "workspace_id" UUID NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "pk_accounts_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_accounts_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "fk_accounts_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES "workspaces"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS profiles(
  "id" UUID NOT NULL,
  "name" VARCHAR NOT NULL,
  "currency" VARCHAR NOT NULL,
  "language" VARCHAR NOT NULL,
  "user_id" UUID NOT NULL,
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "pk_profiles_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_profiles_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE IF NOT EXISTS transactions(
  "id" UUID NOT NULL,
  "title" VARCHAR NOT NULL,
  "note" VARCHAR,
  "currency" VARCHAR,
  "value" NUMERIC(10, 2) NOT NULL,
  "user_id" UUID NOT NULL,
  "workspace_id" UUID NOT NULL,
  "category_id" UUID NOT NULL,
  "account_id" UUID NOT NULL,
  "handled_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMP WITHOUT TIME ZONE,
  CONSTRAINT "pk_transactions_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_transactions_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "fk_transactions_category_id" FOREIGN KEY ("category_id") REFERENCES "categories"("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "fk_transactions_account_id" FOREIGN KEY ("account_id") REFERENCES "accounts"("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT "fk_transactions_workspace_id" FOREIGN KEY ("workspace_id") REFERENCES "workspaces"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

COMMIT;
