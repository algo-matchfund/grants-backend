CREATE TABLE public.users
(
    id uuid NOT NULL,
    website character varying,
    github_user character varying,
    github_authorized boolean NOT NULL DEFAULT false,
    is_company boolean NOT NULL DEFAULT false,
    CONSTRAINT user_id_pk PRIMARY KEY (id)
);

CREATE TABLE public.companies
(
    id serial NOT NULL,
    user_id uuid NOT NULL,
    vat_number character varying NOT NULL,
    CONSTRAINT company_id_pk PRIMARY KEY (id),
    CONSTRAINT vat_number_unique_constraint UNIQUE (vat_number),
    FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.user_socials
(
    id serial NOT NULL,
    user_id uuid NOT NULL,
    github character varying,
    twitter character varying,
    linkedin character varying,
    email character varying,
    website character varying,
    other character varying[],
    CONSTRAINT user_socials_id_pk PRIMARY KEY (id),
    CONSTRAINT user_socials_user_id_fk FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.matching_rounds
(
    id serial NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    match_amount bigint NOT NULL,
    CONSTRAINT matching_round_id_pk PRIMARY KEY (id)

);

CREATE TABLE public.projects
(
    id serial NOT NULL,
    name character varying(60) NOT NULL,
    algorand_wallet character varying(58) NOT NULL,
    description character varying NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    created_by uuid NOT NULL,
    icon character varying,
    background character varying,
    image character varying,
    content character varying,
    app_id integer NOT NULL,
    CONSTRAINT project_id_pk PRIMARY KEY (id),
    CONSTRAINT project_name_unique_constraint UNIQUE (name),
    CONSTRAINT project_created_by_fk FOREIGN KEY (created_by)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.matches
(
    project_id integer NOT NULL,
    fund bigint NOT NULL,
    contributors integer NOT NULL,
    match double precision NOT NULL,
    factor double precision NOT NULL,
    percent double precision NOT NULL,
    CONSTRAINT matches_project_id_pk PRIMARY KEY (project_id),
    CONSTRAINT matches_project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.project_media
(
    id serial NOT NULL,
    project_id integer NOT NULL,
    name character varying(32) NOT NULL,
    media character varying NOT NULL,
    CONSTRAINT project_media_pk PRIMARY KEY (id),
    CONSTRAINT project_media_project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.project_socials
(
    id serial NOT NULL,
    project_id integer NOT NULL,
    github character varying,
    twitter character varying,
    email character varying,
    website character varying,
    other character varying[],
    CONSTRAINT project_social_id_pk PRIMARY KEY (id),
    CONSTRAINT project_socials_project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.funding_campaigns
(
    id serial NOT NULL,
    project_id integer NOT NULL,
    name character varying(32),
    description character varying,
    start_time timestamp NOT NULL,
    end_time timestamp,
    goal bigint,
    CONSTRAINT campaign_id_pk PRIMARY KEY (id),
    CONSTRAINT campaign_project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.project_news
(
    id bigserial NOT NULL,
    project_id integer NOT NULL,
    post_time timestamp NOT NULL DEFAULT NOW(),
    last_edit_time timestamp,
    title character varying NOT NULL,
    content character varying NOT NULL,
    CONSTRAINT news_id_pk PRIMARY KEY (id),
    CONSTRAINT news_project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.project_questions
(
    id bigserial NOT NULL,
    project_id integer NOT NULL,
    user_id uuid NOT NULL,
    ask_time timestamp NOT NULL DEFAULT NOW(),
    title character varying(32),
    content character varying NOT NULL,
    CONSTRAINT project_question_id_pk PRIMARY KEY (id),
    CONSTRAINT question_project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH FULL
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT question_user_id_fk FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH FULL
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.project_question_answers
(
    id bigserial NOT NULL,
    question_id bigint NOT NULL,
    answer_time timestamp NOT NULL DEFAULT NOW(),
    content character varying NOT NULL,
    CONSTRAINT project_answer_id_pk PRIMARY KEY (id),
    CONSTRAINT project_answer_question_id_fk FOREIGN KEY (question_id)
        REFERENCES public.project_questions (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.contributions
(
    id bigserial NOT NULL,
    project_id integer NOT NULL,
    matching_round_id integer NOT NULL,
    contribution_time timestamp NOT NULL DEFAULT now(),
    user_id uuid NOT NULL,
    amount bigint NOT NULL,
    CONSTRAINT contribution_id_pk PRIMARY KEY (id),
    CONSTRAINT project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH FULL
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT matching_round_id_fk FOREIGN KEY (matching_round_id)
        REFERENCES public.matching_rounds (id) MATCH FULL
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT contribution_user_id_fk FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH FULL
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TYPE tx_status AS ENUM ('pending', 'confirmed', 'error');

CREATE TABLE public.algorand_contributions
(
    id bigserial NOT NULL,
    txid character varying NOT NULL UNIQUE,
    sender_address character varying(58) NOT NULL,
    user_id uuid NOT NULL,
    receiver_address character varying(58) NOT NULL,
    project_id integer NOT NULL,
    signature bytea NOT NULL,
    fee bigint NOT NULL,
    amount bigint NOT NULL,
    status tx_status NOT NULL,
    confirmation_round integer,
    CONSTRAINT algorand_contribution_id_pk PRIMARY KEY (id),
    CONSTRAINT algorand_contribution_user_id_fk FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT algorand_contribution_project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.user_subscriptions
(
    id bigserial NOT NULL,
    user_id uuid NOT NULL,
    project_id integer NOT NULL,
    CONSTRAINT subscription_id_pk PRIMARY KEY (id),
    CONSTRAINT subscription_user_id_fk FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH FULL
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT subscription_project_id_fk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH FULL
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.user_notifications
(
    id bigserial NOT NULL,
    user_id uuid NOT NULL,
    "timestamp" timestamp NOT NULL DEFAULT NOW(),
    read boolean NOT NULL DEFAULT false,
    actor character varying NOT NULL,
    object character varying NOT NULL,
    type character varying NOT NULL,
    verb character varying NOT NULL,
    CONSTRAINT notification_id_pk PRIMARY KEY (id),
    CONSTRAINT notification_unique_constraint UNIQUE (actor, object, type),
    CONSTRAINT notification_user_id_fk FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE public.pending_projects
(
    id serial NOT NULL,
    project_id integer,
    name character varying(60),
    algorand_wallet character varying(58) NOT NULL,
    description character varying,
    created_at timestamp NOT NULL DEFAULT NOW(),
    created_by uuid NOT NULL,
    icon character varying,
    background character varying,
    image character varying,
    content character varying,
    github character varying,
    twitter character varying,
    email character varying,
    website character varying,
    other character varying[],
    CONSTRAINT pending_project_id_pk PRIMARY KEY (id),
    CONSTRAINT pending_project_name_unique_constraint UNIQUE (name),
    CONSTRAINT pending_project_project_id_pk FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH FULL
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT pending_project_project_created_by_fk FOREIGN KEY (created_by)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TYPE moderation_status as ENUM('approve', 'deny');

CREATE TABLE public.moderations
(
    id serial NOT NULL,
    moderator_id uuid,
    created_at timestamp NOT NULL DEFAULT NOW(),
    before_project_id integer,
    after_project_id integer NOT NULL,
    status moderation_status,
    comment character varying,
    CONSTRAINT moderation_id_pk PRIMARY KEY (id),
    CONSTRAINT moderations_moderator_id_fk FOREIGN KEY (moderator_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT moderations_before_project_id_fk FOREIGN KEY (before_project_id)
        REFERENCES public.pending_projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT moderations_after_project_id_fk FOREIGN KEY (after_project_id)
        REFERENCES public.pending_projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);


/*
This function sets user's is_company to true if a new row is inserted to companies table
with user_id existing in users table as id.
If a row is deleted, the user's is_company to false
*/
CREATE FUNCTION public.user_toggle_company()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    VOLATILE LEAKPROOF STRICT
AS $BODY$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE users
        SET is_company=true
        WHERE id=NEW.user_id;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE users
        SET is_company=false
        WHERE id=OLD.user_id;
    END IF;

    RETURN NULL;
END;
$BODY$;

CREATE TRIGGER on_new_company_trigger
    AFTER INSERT OR DELETE
    ON public.companies
    FOR EACH ROW
    EXECUTE FUNCTION public.user_toggle_company();


/*
This function inserts new rows about project's news items to user_notifications table.
Notifications are inserted for all users, who are subscribed to the project's updated as described by user_subscriptions table
*/
CREATE FUNCTION public.notify_project_news()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    VOLATILE LEAKPROOF STRICT
AS $BODY$
DECLARE
    subscriber uuid;
    verb character varying;
    news_item_id integer := NEW.id;
    project_id integer := NEW.project_id;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        -- if a news item was deleted, delete all notifications about this news item
        DELETE from public.user_notifications
            WHERE type = 'news_item'
            AND actor = CAST(OLD.project_id as character varying)
            AND object = CAST(OLD.id as character VARYING);

        RETURN NULL;
    END IF;

    IF (TG_OP = 'INSERT') THEN
        verb := 'new';
    ELSIF (TG_OP = 'UPDATE') THEN
        verb := 'updated';
    END IF;

    -- for every project subscriber
    for subscriber in SELECT user_id
        FROM public.user_subscriptions
        WHERE user_subscriptions.project_id=project_id
    loop
        -- in case of a new news item or an existing news item update,
        -- insert notifications for the subscribed users, or replace with
        -- updated verb if notification exists
        INSERT INTO public.user_notifications (user_id, actor, object, type, verb)
        VALUES (subscriber, project_id, news_item_id, 'news_item', verb)
        ON CONFLICT ON CONSTRAINT notification_unique_constraint DO
        UPDATE SET verb = EXCLUDED.verb
            WHERE user_notifications.user_id = subscriber
            AND user_notifications.actor = CAST(project_id as character varying)
            AND user_notifications.object = 'news_item';
    end loop;

    RETURN NULL;
END;
$BODY$;

/*
Insert a news item notification on every operation with project news
*/
CREATE TRIGGER on_project_news_posted_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON public.project_news
    FOR EACH ROW
    EXECUTE FUNCTION public.notify_project_news();

-- Updates news item last_edit_time to current time
CREATE FUNCTION public.update_news_last_update_time()
RETURNS trigger
    LANGUAGE 'plpgsql'
    STABLE LEAKPROOF STRICT
AS $BODY$
BEGIN
    -- edit the newly updated row directly instead of doing another UPDATE query to avoid recursive UPDATE trigger
    NEW.last_edit_time := now();

    RETURN NEW;
END;
$BODY$;

-- Updates news item last_edit_time on updating a news item
CREATE TRIGGER on_project_news_updated_trigger
    BEFORE UPDATE
    ON public.project_news
    FOR EACH ROW
    EXECUTE FUNCTION public.update_news_last_update_time();

CREATE FUNCTION public.notify_project_question()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    VOLATILE LEAKPROOF STRICT
AS $BODY$
DECLARE
    verb character varying;
    project_owner uuid;
    answer_id bigint;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        -- if a new question is inserted, notify the owner of the project
        SELECT created_by INTO project_owner FROM public.projects WHERE id = NEW.project_id;
        INSERT INTO public.user_notifications (user_id, actor, object, type, verb)
        VALUES (project_owner, NEW.project_id, NEW.id, 'project_question', 'new');

        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        -- find all answers to the deleted question and delete notifications about those answers
        for answer_id in SELECT id FROM public.project_question_answers WHERE question_id = OLD.id
        loop
            DELETE FROM public.user_notifications
            WHERE actor = CAST(OLD.project_id as character varying)
                AND type = 'question_answer'
                AND object = CAST(answer_id as character varying);
        end loop;
        -- if a question was deleted by either, all notifications about that question should be deleted
        DELETE FROM public.user_notifications
        WHERE actor = CAST(OLD.project_id as character varying)
            AND type = 'project_question'
            AND object = CAST(OLD.id as character varying);

        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$BODY$;

CREATE TRIGGER on_project_question_changed
    BEFORE INSERT OR DELETE
    ON public.project_questions
    FOR EACH ROW
    EXECUTE FUNCTION public.notify_project_question();

/*
Inserts a notification to a user, whose question's answer is affected,
or if an existing answer was changed
*/
CREATE FUNCTION public.notify_question_answered()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    VOLATILE LEAKPROOF STRICT
AS $BODY$
DECLARE
    verb character varying;
    question_author uuid;
    project_id integer;
    question_id bigint := NEW.question_id;
    answer_id bigint := NEW.id;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        verb := 'new';
    ELSIF (TG_OP = 'UPDATE') THEN
        verb := 'updated';
    END IF;

    -- find the author of the answered question
    SELECT project_questions.user_id, project_questions.project_id INTO question_author, project_id FROM public.project_questions WHERE id = question_id;
    IF (question_author = NULL OR project_id = NULL) THEN
        RETURN NULL;
    END IF;

    -- notify the question's author that the question was answered,
    -- or replace existing notification with updated verb
    INSERT INTO public.user_notifications (user_id, actor, object, type, verb)
    VALUES (question_author, project_id, answer_id, 'question_answer', verb)
    ON CONFLICT ON CONSTRAINT notification_unique_constraint DO
        UPDATE SET verb = EXCLUDED.verb
            WHERE user_notifications.user_id = question_author
            AND user_notifications.actor = CAST(project_id as character varying)
            AND user_notifications.object = CAST(answer_id as character varying);

    RETURN NULL;
END;
$BODY$;

CREATE TRIGGER on_project_question_answered
    AFTER INSERT OR UPDATE
    ON public.project_question_answers
    FOR EACH ROW
    EXECUTE FUNCTION public.notify_question_answered();


ALTER TABLE IF EXISTS public.users
    OWNER to grants;
ALTER TABLE IF EXISTS public.companies
    OWNER to grants;
ALTER TABLE IF EXISTS public.user_socials
    OWNER to grants;
ALTER TABLE IF EXISTS public.projects
    OWNER to grants;
ALTER TABLE IF EXISTS public.project_media
    OWNER to grants;
ALTER TABLE IF EXISTS public.project_socials
    OWNER to grants;
ALTER TABLE IF EXISTS public.funding_campaigns
    OWNER to grants;
ALTER TABLE IF EXISTS public.project_news
    OWNER to grants;
ALTER TABLE IF EXISTS public.project_question_answers
    OWNER to grants;
ALTER TABLE IF EXISTS public.contributions
    OWNER to grants;
ALTER TABLE IF EXISTS public.algorand_contributions
    OWNER to grants;
ALTER TABLE IF EXISTS public.user_subscriptions
    OWNER to grants;
ALTER TABLE IF EXISTS public.user_notifications
    OWNER to grants;



ALTER FUNCTION public.user_toggle_company()
    OWNER TO grants;
ALTER FUNCTION public.notify_project_news()
    OWNER TO grants;
ALTER FUNCTION public.notify_question_answered()
    OWNER TO grants;
ALTER FUNCTION public.notify_project_question()
    OWNER TO grants;
