# Password reset

When authentication is enabled, you may want to reset it.
To do so, you have two option.

## If you know your current password

If you know your current password, just go under your profile's page and select the "Security/Password" tab
(the link is `/profile/settings`).
You will be able to reset your password using your current password.

## If you lost your password

In this case, you can still reset your password by running a special command with the binary.
You will need to have the same configuration as when your ZobTube instance is running.

First, start ZobTube with just the command and no other parameters to list users and their ID.

```
./zobtube password-reset
2026-05-24T11:42:57+02:00 ??? setting up configuration
2026-05-24T11:42:57+02:00 INF valid configuration found db-driver=sqlite media-path=data metadata-path=data metadata-type=filesystem server-bind=0.0.0.0:8069
2026-05-24T11:42:57+02:00 ??? initializing database connection
2026-05-24T11:42:57+02:00 ??? get user list
2026-05-24T11:42:57+02:00 ??? * ID: 0548f2c1-550c-4733-997a-2e731823106c (username: admin)
2026-05-24T11:42:57+02:00 ??? please now use the --user-id flag to select the user
```

Now you can start the command with the `--user-id` parameter.

```
./zobtube password-reset --user-id 0548f2c1-550c-4733-997a-2e731823106c
2026-05-24T11:44:30+02:00 ??? setting up configuration
2026-05-24T11:44:30+02:00 INF valid configuration found db-driver=sqlite media-path=data metadata-path=data metadata-type=filesystem server-bind=0.0.0.0:8069
2026-05-24T11:44:30+02:00 ??? initializing database connection
2026-05-24T11:44:30+02:00 ??? get selected user user-id=0548f2c1-550c-4733-997a-2e731823106c
2026-05-24T11:44:30+02:00 ??? new password for user admin will be D27EOGOICZAZBZ6CC7CEEBY2KB user-id=0548f2c1-550c-4733-997a-2e731823106c
2026-05-24T11:44:30+02:00 INF new password set successfully user-id=0548f2c1-550c-4733-997a-2e731823106c
```

Your new password will be in the second last log line.