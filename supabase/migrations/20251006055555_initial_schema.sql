create table "public"."conversations" (
    "id" uuid not null default uuid_generate_v4(),
    "user_id" text not null,
    "role" text not null,
    "content" text not null,
    "created_at" timestamp with time zone not null default now()
);


alter table "public"."conversations" enable row level security;

create table "public"."users" (
    "id" uuid not null,
    "line_id" text not null,
    "name" text,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


alter table "public"."users" enable row level security;

CREATE UNIQUE INDEX conversations_pkey ON public.conversations USING btree (id);

CREATE INDEX idx_conversations_created_at ON public.conversations USING btree (created_at DESC);

CREATE INDEX idx_conversations_user_id ON public.conversations USING btree (user_id);

CREATE INDEX idx_conversations_user_id_created_at ON public.conversations USING btree (user_id, created_at DESC);

CREATE INDEX idx_users_line_id ON public.users USING btree (line_id);

CREATE UNIQUE INDEX users_line_id_key ON public.users USING btree (line_id);

CREATE UNIQUE INDEX users_pkey ON public.users USING btree (id);

alter table "public"."conversations" add constraint "conversations_pkey" PRIMARY KEY using index "conversations_pkey";

alter table "public"."users" add constraint "users_pkey" PRIMARY KEY using index "users_pkey";

alter table "public"."conversations" add constraint "conversations_role_check" CHECK ((role = ANY (ARRAY['user'::text, 'assistant'::text]))) not valid;

alter table "public"."conversations" validate constraint "conversations_role_check";

alter table "public"."users" add constraint "users_id_fkey" FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE not valid;

alter table "public"."users" validate constraint "users_id_fkey";

alter table "public"."users" add constraint "users_line_id_key" UNIQUE using index "users_line_id_key";

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

grant delete on table "public"."conversations" to "anon";

grant insert on table "public"."conversations" to "anon";

grant references on table "public"."conversations" to "anon";

grant select on table "public"."conversations" to "anon";

grant trigger on table "public"."conversations" to "anon";

grant truncate on table "public"."conversations" to "anon";

grant update on table "public"."conversations" to "anon";

grant delete on table "public"."conversations" to "authenticated";

grant insert on table "public"."conversations" to "authenticated";

grant references on table "public"."conversations" to "authenticated";

grant select on table "public"."conversations" to "authenticated";

grant trigger on table "public"."conversations" to "authenticated";

grant truncate on table "public"."conversations" to "authenticated";

grant update on table "public"."conversations" to "authenticated";

grant delete on table "public"."conversations" to "service_role";

grant insert on table "public"."conversations" to "service_role";

grant references on table "public"."conversations" to "service_role";

grant select on table "public"."conversations" to "service_role";

grant trigger on table "public"."conversations" to "service_role";

grant truncate on table "public"."conversations" to "service_role";

grant update on table "public"."conversations" to "service_role";

grant delete on table "public"."users" to "anon";

grant insert on table "public"."users" to "anon";

grant references on table "public"."users" to "anon";

grant select on table "public"."users" to "anon";

grant trigger on table "public"."users" to "anon";

grant truncate on table "public"."users" to "anon";

grant update on table "public"."users" to "anon";

grant delete on table "public"."users" to "authenticated";

grant insert on table "public"."users" to "authenticated";

grant references on table "public"."users" to "authenticated";

grant select on table "public"."users" to "authenticated";

grant trigger on table "public"."users" to "authenticated";

grant truncate on table "public"."users" to "authenticated";

grant update on table "public"."users" to "authenticated";

grant delete on table "public"."users" to "service_role";

grant insert on table "public"."users" to "service_role";

grant references on table "public"."users" to "service_role";

grant select on table "public"."users" to "service_role";

grant trigger on table "public"."users" to "service_role";

grant truncate on table "public"."users" to "service_role";

grant update on table "public"."users" to "service_role";

create policy "Users can delete own conversations"
on "public"."conversations"
as permissive
for delete
to public
using ((user_id = ( SELECT users.line_id
   FROM users
  WHERE (users.id = auth.uid()))));


create policy "Users can insert own conversations"
on "public"."conversations"
as permissive
for insert
to public
with check ((user_id = ( SELECT users.line_id
   FROM users
  WHERE (users.id = auth.uid()))));


create policy "Users can view own conversations"
on "public"."conversations"
as permissive
for select
to public
using ((user_id = ( SELECT users.line_id
   FROM users
  WHERE (users.id = auth.uid()))));


create policy "Users can update own data"
on "public"."users"
as permissive
for update
to public
using ((auth.uid() = id));


create policy "Users can view own data"
on "public"."users"
as permissive
for select
to public
using ((auth.uid() = id));


CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


