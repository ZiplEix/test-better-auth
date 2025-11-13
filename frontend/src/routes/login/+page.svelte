<script lang="ts">
    import { authClient } from '$lib/auth-client';

    let email = $state('');
    let password = $state('');
    let error = $state('');
    let loading = $state(false);
    let socialLoading = $state(false);

    const handleSubmit = async () => {
        error = '';
        loading = true;

        try {
            await authClient.signIn.email(
                {
                    email,
                    password,
                    callbackURL: '/me'
                },
                {
                    onError: (ctx) => {
                        error = ctx.error?.message ?? 'Une erreur est survenue lors de la connexion.';
                    },
                    onResponse: () => {
                        loading = false;
                    },
                    onSuccess: () => {
                        // redirect in case of no callbackURL
                    },
                }
            );
        } catch (e: any) {
            console.error(e);
            error = "Erreur inatendu";
            loading = false;
        }
    }

    const handleGoogleSignIn = async () => {
        error = '';
        socialLoading = true;

        try {
            // This will redirect the user to Google's consent screen
            await authClient.signIn.social({
                provider: 'google',
                callbackURL: '/me',          // où revenir après login OK
                errorCallbackURL: '/login' // où revenir si erreur
            });
        } catch (e: any) {
            console.error(e);
            error = "Impossible de se connecter avec Google.";
            socialLoading = false;
        }
    };
</script>

<div class="min-h-screen flex items-center justify-center bg-linear-to-b from-slate-50 to-white px-4">
    <div class="w-full max-w-md">
        <div class="bg-white shadow-lg rounded-lg p-8">
            <h1 class="text-2xl font-semibold text-slate-800 mb-1">Se connecter</h1>
            <p class="text-sm text-slate-500 mb-6">Entrez votre email et mot de passe pour continuer.</p>

            <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Email</label>
                    <input
                        type="email"
                        bind:value={email}
                        required
                        class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="ex: vous@exemple.com"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Mot de passe</label>
                    <input
                        type="password"
                        bind:value={password}
                        required
                        class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="Votre mot de passe"
                    />
                </div>

                {#if error}
                    <p class="text-sm text-red-600">{error}</p>
                {/if}

                <button
                    type="submit"
                    class="w-full py-2 rounded-md bg-blue-600 text-white font-medium hover:bg-blue-700 transition disabled:opacity-60"
                    disabled={loading}
                >
                    {#if loading}
                        Connexion...
                    {:else}
                        Se connecter
                    {/if}
                </button>

                <div class="mt-4 border-t pt-4">
                    <p class="text-xs text-slate-400 text-center mb-3">
                        ou
                    </p>

                    <button
                        type="button"
                        class="w-full py-2 rounded-md border border-slate-300 bg-white text-slate-800 font-medium hover:bg-slate-50 transition disabled:opacity-60 flex items-center justify-center gap-2"
                        on:click={handleGoogleSignIn}
                        disabled={socialLoading}
                    >
                        {#if socialLoading}
                            Connexion avec Google...
                        {:else}
                            <!-- simple "G" badge, tu peux remplacer par un vrai logo -->
                            <span class="inline-flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-white text-xs font-bold">
                                G
                            </span>
                            <span>Continuer avec Google</span>
                        {/if}
                    </button>
                </div>

                <div class="text-center text-sm text-slate-600">
                    Pas encore de compte ? <a href="/signup" class="text-blue-600 underline">S'inscrire</a>
                </div>
            </form>
        </div>
    </div>
</div>
