import { Pool } from 'pg';
import { BETTER_AUTH_SECRET, BETTER_AUTH_URL, DATABASE_URL, GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET } from '$env/static/private';
import { betterAuth } from 'better-auth';
import { jwt } from 'better-auth/plugins';

const pool = new Pool({
    connectionString: DATABASE_URL
});

export const auth = betterAuth({
    secret: BETTER_AUTH_SECRET,
    url: BETTER_AUTH_URL,
    database: pool,
    emailAndPassword: { enabled: true },

    plugins: [
        jwt(),
    ],

    socialProviders: {
        google: {
            clientId: GOOGLE_CLIENT_ID,
            clientSecret: GOOGLE_CLIENT_SECRET,
            prompt: "select_account"
        }
    }
})
