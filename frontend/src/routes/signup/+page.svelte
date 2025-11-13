<script lang="ts">
    import { authClient } from '$lib/auth-client';

    let name = $state('');
    let email = $state('');
    let password = $state('');
    let confirmPassword = $state('');
    let error = $state('');
    let loading = $state(false);

    const handleSubmit = async () => {
        error = '';

        if (password !== confirmPassword) {
            error = 'Les mots de passe ne correspondent pas.';
            return;
        }

        loading = true;

        try {
            await authClient.signUp.email(
                {
                    email,
                    password,
                    name,
                    callbackURL: '/' // ou '/login' si tu préfères
                },
                {
                    onError: (ctx) => {
                        error = ctx.error?.message ?? "Une erreur est survenue lors de l'inscription.";
                    },
                    onResponse: () => {
                        loading = false;
                    },
                    onSuccess: () => {
                        // si pas de callbackURL, tu peux rediriger ici
                    }
                }
            );
        } catch (e: any) {
            console.error(e);
            error = 'Erreur inattendue';
            loading = false;
        }
    };
</script>

<div class="min-h-screen flex items-center justify-center bg-linear-to-b from-slate-50 to-white px-4">
    <div class="w-full max-w-md">
        <div class="bg-white shadow-lg rounded-lg p-8">
            <h1 class="text-2xl font-semibold text-slate-800 mb-1">S'inscrire</h1>
            <p class="text-sm text-slate-500 mb-6">Créez un compte pour commencer.</p>

            <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Nom</label>
                    <input
                        type="text"
                        bind:value={name}
                        required
                        class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-violet-500"
                        placeholder="Votre nom complet"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Email</label>
                    <input
                        type="email"
                        bind:value={email}
                        required
                        class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-violet-500"
                        placeholder="ex: vous@exemple.com"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Mot de passe</label>
                    <input
                        type="password"
                        bind:value={password}
                        required
                        class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-violet-500"
                        placeholder="Créer un mot de passe"
                    />
                </div>

                <div>
                    <label class="block text-sm font-medium text-slate-700 mb-1">Confirmer le mot de passe</label>
                    <input
                        type="password"
                        bind:value={confirmPassword}
                        required
                        class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-violet-500"
                        placeholder="Confirmez le mot de passe"
                    />
                </div>

                {#if error}
                    <p class="text-sm text-red-600">{error}</p>
                {/if}

                <button
                    type="submit"
                    class="w-full py-2 rounded-md bg-violet-600 text-white font-medium hover:bg-violet-700 transition disabled:opacity-60"
                    disabled={loading}
                >
                    {#if loading}
                        Inscription...
                    {:else}
                        S'inscrire
                    {/if}
                </button>

                <div class="text-center text-sm text-slate-600">
                    Déjà un compte ?
                    <a href="/login" class="text-violet-600 underline">Se connecter</a>
                </div>
            </form>
        </div>
    </div>
</div>
