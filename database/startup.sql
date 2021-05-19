CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.questions
(
    uuid uuid NOT NULL DEFAULT uuid_generate_v1(),
    key text COLLATE pg_catalog."default",
    title text COLLATE pg_catalog."default",
    question text COLLATE pg_catalog."default",
    CONSTRAINT questions_pkey PRIMARY KEY (uuid)
)

TABLESPACE pg_default;

ALTER TABLE public.questions
    OWNER to postgres;

INSERT INTO public.questions(
    key, title, question)
VALUES  ('hello-world', 'Hello World', 'Please output Hello World'),
    ('two-sum', 'Two Sum', 'Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.You may assume that each input would have exactly one solution, and you may not use the same element twice. You can return the answer in any order.'),
    ('add-two-numbers', 'Add Two Numbers','You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list. You may assume the two numbers do not contain any leading zero, except the number 0 itself.'),
    ('longest-substring-without-repeating-characters', 'Longest Substring Without Repeating Characters','Given a string s, find the length of the longest substring without repeating characters.'),
    ('median-of-two-sorted-arrays', 'Median of Two Sorted Arrays','Given two sorted arrays nums1 and nums2 of size m and n respectively, return the median of the two sorted arrays. The overall run time complexity should be O(log (m+n)).');