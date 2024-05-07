<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import PlantCard from '$lib/PlantCard.svelte';
	import type { Plant } from '$lib/types.js';

	const { data } = $props()

	let sortedPlants = $derived(data.plants.toSorted((x,y) => new Date(x.lastWateredAt) - new Date(y.lastWateredAt)))

	let showForm = $state(false);

	let name = $state("");
	let place = $state("");

	async function addPlant() {
		const response = await fetch(`/api/plants`, {
			method: 'POST',
			body: JSON.stringify({ name, place }),
			headers: {
				'content-type': 'application/json'
			}
		});
		showForm=false
		invalidateAll()
	}	
</script>

<div>
	<div class="mx-auto max-w-screen-md">
		<div class="flex items-center justify-center gap-x-4 py-12">
			<h1 class="text-4xl font-bold">My Plants</h1>
		</div>
		<div class="space-y-2 px-2">
			<div>
				<button
					onclick={() => (showForm = true)}
					hidden={showForm}
					class="mb-4 w-full rounded-lg bg-green-400 px-4 py-2 text-lg hover:bg-green-300"
					>+ Add Plant</button
				>
				<div hidden={!showForm} class="w-full shadow bg-green-300 p-4 rounded-lg grid sm:grid-cols-2 sm:grid-rows-[2fr_1fr] gap-4">
					<label class=" gap-2 flex justify-between items-center" for="name">
						<span class="font-bold">Name:</span>
						<input class="flex-grow" bind:value={name} type="text" name="name" />
					</label>
					<label class=" flex justify-between items-center gap-2" for="place"
						>
						<span class="font-bold">Place:</span>
						<input class="flex-grow" bind:value={place} type="text" name="place" />
					</label>
					<div class="ml-auto sm:col-start-2">
						<button onclick={addPlant} class="bg-blue-600 hover:bg-blue-500 text-white font-bold px-4 py-2 rounded-lg">OK</button>
						<button onclick={() => {name="";place="";showForm=false}} class="bg-zinc-300 hover:bg-zinc-200 px-4 py-2 rounded-lg">Cancel</button>
					</div>
				</div>
			</div>

			{#each sortedPlants as plant}
				<PlantCard {plant}/>
			{/each}
		</div>
	</div>
</div>
