DROP TRIGGER IF EXISTS on_new_company_trigger ON public.companies;
DROP FUNCTION IF EXISTS public.user_toggle_company;
DROP TRIGGER IF EXISTS on_project_news_posted_trigger ON public.project_news;
DROP FUNCTION IF EXISTS public.notify_project_news;
DROP TRIGGER IF EXISTS on_project_news_updated_trigger ON public.project_news;
DROP FUNCTION IF EXISTS public.update_news_last_update_time;
DROP TRIGGER IF EXISTS on_project_question_changed ON public.project_questions;
DROP FUNCTION IF EXISTS public.notify_project_question;
DROP TRIGGER IF EXISTS on_project_question_answered ON public.project_question_answers;
DROP FUNCTION IF EXISTS public.notify_question_answered;

DROP TABLE IF EXISTS public.users CASCADE;
DROP TABLE IF EXISTS public.companies CASCADE;
DROP TABLE IF EXISTS public.user_socials CASCADE;
DROP TABLE IF EXISTS public.matching_rounds CASCADE;
DROP TABLE IF EXISTS public.projects CASCADE;
DROP TABLE IF EXISTS public.project_media CASCADE;
DROP TABLE IF EXISTS public.project_socials CASCADE;
DROP TABLE IF EXISTS public.funding_campaigns CASCADE;
DROP TABLE IF EXISTS public.project_news CASCADE;
DROP TABLE IF EXISTS public.project_questions CASCADE;
DROP TABLE IF EXISTS public.project_question_answers CASCADE;
DROP TABLE IF EXISTS public.contributions CASCADE;
DROP TABLE IF EXISTS public.algorand_contributions CASCADE;
DROP TABLE IF EXISTS public.matches CASCADE;
DROP TABLE IF EXISTS public.user_subscriptions CASCADE;
DROP TABLE IF EXISTS public.user_notifications CASCADE;
DROP TABLE IF EXISTS public.pending_projects CASCADE;
DROP TABLE IF EXISTS public.moderations CASCADE;

DROP TYPE IF EXISTS tx_status;
DROP TYPE IF EXISTS moderation_status;
