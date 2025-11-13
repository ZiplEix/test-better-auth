<script lang="ts">
    import { authClient } from '$lib/auth-client';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    type SessionUser = {
        name?: string | null;
        email?: string | null;
        image?: string | null;
    };

    let loading = $state(true);
    let error = $state('');
    let user = $state<SessionUser | null>(null);

    onMount(async () => {
        // Fetch session on mount (client-side)
        try {
            const { data, error: sessionError } = await authClient.getSession();

            // handle possible error returned by better-auth client
            if (sessionError) {
                console.error(sessionError);
                error = "Impossible de récupérer votre session.";
                return;
            }

            // data is { user, session } | null
            user = data?.user ?? null;
        } catch (e) {
            console.error(e);
            error = "Impossible de récupérer votre session.";
        } finally {
            loading = false;
        }
    });

    const handleLogout = async () => {
        error = '';

        try {
            await authClient.signOut();
            goto('/login');
        } catch (e) {
            console.error(e);
            error = "Une erreur est survenue lors de la déconnexion.";
        }
    };

    // small helper to build initials when there is no image
    const getInitials = (user: SessionUser | null): string => {
        if (!user) return '?';

        if (user.name && user.name.trim().length > 0) {
            const parts = user.name.trim().split(' ');
            if (parts.length === 1) {
                return parts[0].charAt(0).toUpperCase();
            }
            return (
                (parts[0]?.charAt(0) ?? '').toUpperCase() +
                (parts[1]?.charAt(0) ?? '').toUpperCase()
            );
        }

        if (user.email && user.email.length > 0) {
            return user.email.charAt(0).toUpperCase();
        }

        return '?';
    };
</script>

<div class="min-h-screen flex items-center justify-center bg-linear-to-b from-slate-50 to-white px-4">
    <div class="w-full max-w-md">
        <div class="bg-white shadow-lg rounded-lg p-8">
            <h1 class="text-2xl font-semibold text-slate-800 mb-1">Mon compte</h1>
            <p class="text-sm text-slate-500 mb-6">
                Consultez les informations de votre profil.
            </p>

            {#if loading}
                <p class="text-slate-600 text-center">Chargement...</p>
            {:else if user}
                <div class="space-y-4">
                    <!-- Avatar block -->
                    <div class="flex items-center gap-3">
                        {#if user.image}
                            <img
                                src={user.image}
                                alt="Avatar"
                                class="h-12 w-12 rounded-full border border-slate-200 object-cover"
                            />
                        {:else}
                            <div
                                class="h-12 w-12 rounded-full bg-violet-600 text-white flex items-center justify-center text-lg font-semibold"
                            >
                                {getInitials(user)}
                            </div>
                        {/if}

                        <div>
                            <p class="text-sm font-medium text-slate-800">
                                {user.name || 'Nom non renseigné'}
                            </p>
                            <p class="text-xs text-slate-500">
                                {user.email}
                            </p>
                        </div>
                    </div>

                    <hr class="border-slate-200" />

                    <div>
                        <p class="text-sm font-medium text-slate-700 mb-1">Nom</p>
                        <p class="px-3 py-2 border rounded-md bg-slate-50">
                            {user.name || 'Non renseigné'}
                        </p>
                    </div>

                    <div>
                        <p class="text-sm font-medium text-slate-700 mb-1">Email</p>
                        <p class="px-3 py-2 border rounded-md bg-slate-50">
                            {user.email}
                        </p>
                    </div>

                    {#if error}
                        <p class="text-sm text-red-600">{error}</p>
                    {/if}

                    <button
                        type="button"
                        class="w-full py-2 rounded-md bg-violet-600 text-white font-medium hover:bg-violet-700 transition"
                        on:click={handleLogout}
                    >
                        Se déconnecter
                    </button>

                    <div class="text-center text-sm text-slate-600">
                        <a href="/" class="text-violet-600 underline">Retour à l'accueil</a>
                    </div>
                </div>
            {:else}
                <p class="text-center text-slate-600 mb-4">
                    Vous n'êtes pas connecté.
                </p>

                <a
                    href="/login"
                    class="block w-full text-center py-2 rounded-md bg-violet-600 text-white font-medium hover:bg-violet-700 transition"
                >
                    Se connecter
                </a>
            {/if}
        </div>
    </div>
</div>
