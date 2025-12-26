SELECT emp_id,
       name,
       wallet_address,
       address,
       password_hash,
       password_salt
FROM public.customers
-- WHERE emp_id = '536c28c8-4afa-513e-ad48-9b5df13fe24f'
LIMIT 1000;