<script lang="ts">
    import { authClient } from '$lib/auth-client';
    import { goto } from '$app/navigation';
    import { onMount } from 'svelte';

    type Todo = {
        id: number;
        title: string;
        completed: boolean;
    };

    // API base URL (you can move this to PUBLIC_API_URL later)
    const API_URL = 'http://localhost:8080/api';

    let loading = $state(true);
    let error = $state('');
    let todos = $state<Todo[]>([]);
    let newTitle = $state('');

    let jwtToken = $state<string | null>(null);

    // Fetch JWT token + todos on mount
    onMount(async () => {
        try {
            // get JWT token from Better Auth
            const { data, error: tokenError } = await authClient.token();

            if (tokenError || !data?.token) {
                console.error(tokenError);
                error = "Impossible de récupérer le token d'authentification.";
                // optional redirect to login
                goto('/login');
                return;
            }

            jwtToken = data.token;

            await fetchTodos();
        } catch (e) {
            console.error(e);
            error = "Une erreur est survenue lors du chargement des tâches.";
        } finally {
            loading = false;
        }
    });

    const fetchTodos = async () => {
        if (!jwtToken) {
            return;
        }

        error = '';

        try {
            const res = await fetch(`${API_URL}/todos`, {
                headers: {
                    Authorization: `Bearer ${jwtToken}`
                }
            });

            if (!res.ok) {
                throw new Error(`Erreur serveur: ${res.status}`);
            }

            const data: Todo[] = await res.json();
            todos = data;
        } catch (e) {
            console.error(e);
            error = "Impossible de récupérer la liste des tâches.";
        }
    };

    const handleAddTodo = async () => {
        if (!jwtToken) {
            error = "Vous n'êtes pas authentifié.";
            return;
        }

        if (!newTitle.trim()) {
            error = 'Le titre de la tâche ne peut pas être vide.';
            return;
        }

        error = '';

        try {
            const res = await fetch(`${API_URL}/todos`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${jwtToken}`
                },
                body: JSON.stringify({ title: newTitle.trim() })
            });

            if (!res.ok) {
                throw new Error(`Erreur serveur: ${res.status}`);
            }

            newTitle = '';
            await fetchTodos();
        } catch (e) {
            console.error(e);
            error = "Impossible d'ajouter la tâche.";
        }
    };

    const handleToggleTodo = async (todo: Todo) => {
        if (!jwtToken) {
            error = "Vous n'êtes pas authentifié.";
            return;
        }

        error = '';

        try {
            const res = await fetch(`${API_URL}/todos/${todo.id}/toggle`, {
                method: 'PATCH',
                headers: {
                    Authorization: `Bearer ${jwtToken}`
                }
            });

            if (!res.ok) {
                throw new Error(`Erreur serveur: ${res.status}`);
            }

            // re-fetch list after toggle
            await fetchTodos();
        } catch (e) {
            console.error(e);
            error = "Impossible de mettre à jour la tâche.";
        }
    };

    const handleDeleteTodo = async (todo: Todo) => {
        if (!jwtToken) {
            error = "Vous n'êtes pas authentifié.";
            return;
        }

        error = '';

        try {
            const res = await fetch(`${API_URL}/todos/${todo.id}`, {
                method: 'DELETE',
                headers: {
                    Authorization: `Bearer ${jwtToken}`
                }
            });

            if (!res.ok) {
                throw new Error(`Erreur serveur: ${res.status}`);
            }

            // re-fetch list after delete
            await fetchTodos();
        } catch (e) {
            console.error(e);
            error = "Impossible de supprimer la tâche.";
        }
    };
</script>

<div class="min-h-screen flex items-center justify-center bg-linear-to-b from-slate-50 to-white px-4">
    <div class="w-full max-w-md">
        <div class="bg-white shadow-lg rounded-lg p-8">
            <h1 class="text-2xl font-semibold text-slate-800 mb-1">Mes tâches</h1>
            <p class="text-sm text-slate-500 mb-6">
                Gérez une petite liste de tâches pour tester la stack Better Auth + Go.
            </p>

            {#if loading}
                <p class="text-slate-600 text-center">Chargement...</p>
            {:else}
                <div class="space-y-4">
                    <!-- Add todo form -->
                    <form
                        class="flex gap-2"
                        on:submit|preventDefault={handleAddTodo}
                    >
                        <input
                            type="text"
                            bind:value={newTitle}
                            class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-violet-500"
                            placeholder="Nouvelle tâche..."
                        />
                        <button
                            type="submit"
                            class="px-4 py-2 rounded-md bg-violet-600 text-white font-medium hover:bg-violet-700 transition"
                        >
                            Ajouter
                        </button>
                    </form>

                    {#if error}
                        <p class="text-sm text-red-600">{error}</p>
                    {/if}

                    <!-- Todos list -->
                    {#if todos.length === 0}
                        <p class="text-sm text-slate-500">
                            Aucune tâche pour le moment. Ajoutez-en une pour commencer.
                        </p>
                    {:else}
                        <ul class="space-y-2">
                            {#each todos as todo}
                                <li class="flex items-center justify-between gap-2 px-3 py-2 border rounded-md">
                                    <div class="flex items-center gap-2">
                                        <input
                                            type="checkbox"
                                            checked={todo.completed}
                                            on:change={() => handleToggleTodo(todo)}
                                            class="h-4 w-4 rounded border-slate-300 text-violet-600 focus:ring-violet-500"
                                        />
                                        <span
                                            class="text-sm text-slate-800"
                                            class:line-through={todo.completed}
                                            class:text-slate-400={todo.completed}
                                        >
                                            {todo.title}
                                        </span>
                                    </div>

                                    <button
                                        type="button"
                                        class="text-xs text-red-600 hover:text-red-700"
                                        on:click={() => handleDeleteTodo(todo)}
                                    >
                                        Supprimer
                                    </button>
                                </li>
                            {/each}
                        </ul>
                    {/if}

                    <div class="text-center text-sm text-slate-600 pt-2">
                        <a href="/me" class="text-violet-600 underline">Voir mon compte</a>
                    </div>
                </div>
            {/if}
        </div>
    </div>
</div>
