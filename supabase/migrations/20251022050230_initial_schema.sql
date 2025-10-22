create type "public"."role" as enum ('user', 'assistant');

create table "public"."conversation" (
    "id" uuid not null default uuid_generate_v4(),
    "user_id" uuid not null default gen_random_uuid(),
    "content" text not null,
    "created_at" timestamp with time zone not null default now(),
    "role" role not null
);


alter table "public"."conversation" enable row level security;

create table "public"."user" (
    "id" uuid not null,
    "line_id" text not null,
    "name" text not null default '""'::text,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


alter table "public"."user" enable row level security;

CREATE UNIQUE INDEX conversations_pkey ON public.conversation USING btree (id);

CREATE INDEX idx_conversations_created_at ON public.conversation USING btree (created_at DESC);

CREATE INDEX idx_conversations_user_id ON public.conversation USING btree (user_id);

CREATE INDEX idx_conversations_user_id_created_at ON public.conversation USING btree (user_id, created_at DESC);

CREATE INDEX idx_users_line_id ON public."user" USING btree (line_id);

CREATE UNIQUE INDEX users_line_id_key ON public."user" USING btree (line_id);

CREATE UNIQUE INDEX users_pkey ON public."user" USING btree (id);

alter table "public"."conversation" add constraint "conversations_pkey" PRIMARY KEY using index "conversations_pkey";

alter table "public"."user" add constraint "users_pkey" PRIMARY KEY using index "users_pkey";

alter table "public"."conversation" add constraint "conversations_user_id_fkey" FOREIGN KEY (user_id) REFERENCES "user"(id) not valid;

alter table "public"."conversation" validate constraint "conversations_user_id_fkey";

alter table "public"."user" add constraint "users_id_fkey" FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE not valid;

alter table "public"."user" validate constraint "users_id_fkey";

alter table "public"."user" add constraint "users_line_id_key" UNIQUE using index "users_line_id_key";

set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.update_updated_at_column()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$function$
;

grant delete on table "public"."conversation" to "anon";

grant insert on table "public"."conversation" to "anon";

grant references on table "public"."conversation" to "anon";

grant select on table "public"."conversation" to "anon";

grant trigger on table "public"."conversation" to "anon";

grant truncate on table "public"."conversation" to "anon";

grant update on table "public"."conversation" to "anon";

grant delete on table "public"."conversation" to "authenticated";

grant insert on table "public"."conversation" to "authenticated";

grant references on table "public"."conversation" to "authenticated";

grant select on table "public"."conversation" to "authenticated";

grant trigger on table "public"."conversation" to "authenticated";

grant truncate on table "public"."conversation" to "authenticated";

grant update on table "public"."conversation" to "authenticated";

grant delete on table "public"."conversation" to "service_role";

grant insert on table "public"."conversation" to "service_role";

grant references on table "public"."conversation" to "service_role";

grant select on table "public"."conversation" to "service_role";

grant trigger on table "public"."conversation" to "service_role";

grant truncate on table "public"."conversation" to "service_role";

grant update on table "public"."conversation" to "service_role";

grant delete on table "public"."user" to "anon";

grant insert on table "public"."user" to "anon";

grant references on table "public"."user" to "anon";

grant select on table "public"."user" to "anon";

grant trigger on table "public"."user" to "anon";

grant truncate on table "public"."user" to "anon";

grant update on table "public"."user" to "anon";

grant delete on table "public"."user" to "authenticated";

grant insert on table "public"."user" to "authenticated";

grant references on table "public"."user" to "authenticated";

grant select on table "public"."user" to "authenticated";

grant trigger on table "public"."user" to "authenticated";

grant truncate on table "public"."user" to "authenticated";

grant update on table "public"."user" to "authenticated";

grant delete on table "public"."user" to "service_role";

grant insert on table "public"."user" to "service_role";

grant references on table "public"."user" to "service_role";

grant select on table "public"."user" to "service_role";

grant trigger on table "public"."user" to "service_role";

grant truncate on table "public"."user" to "service_role";

grant update on table "public"."user" to "service_role";

create policy "Users can update own data"
on "public"."user"
as permissive
for update
to public
using ((auth.uid() = id));


create policy "Users can view own data"
on "public"."user"
as permissive
for select
to public
using ((auth.uid() = id));


CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public."user" FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


