<script lang="ts">
	import { writable } from 'svelte/store';
	import type { Plant } from './types';
	import { invalidate, invalidateAll } from '$app/navigation';

	let {plant} = $props()

	let newName = $state(plant.name)
	let newPlace = $state(plant.place)

	let showEdit = $state(false)

	let lastWateredDate = $derived(formatISODate(plant.lastWateredAt))

  	async function waterNow() {
		const currentDate = new Date().toISOString()
		plant.lastWateredAt = currentDate

		const response = await fetch(`/api/plants`, {
			method: 'PUT',
			body: JSON.stringify(plant),
			headers: {
				'content-type': 'application/json'
			}
		});
		console.log(plant)
	}

	async function postPlant() {
		plant.name = newName
		plant.place = newPlace

		const response = await fetch(`/api/plants`, {
			method: 'PUT',
			body: JSON.stringify(plant),
			headers: {
				'content-type': 'application/json'
			}
		});
		console.log(plant)
		showEdit = false
	}

	async function deletePlant() {
		const response = await fetch(`/api/plants`, {
			method: 'DELETE',
			body: JSON.stringify(plant),
			headers: {
				'content-type': 'application/json'
			}
		});
		console.log(plant)
		showEdit = false
		invalidateAll()
	}

	function formatISODate(isoDateString:string) {
    // Create a Date object from the ISO string
    const date = new Date(isoDateString);

    // Extract day, month, and year from the Date object
    const day = date.getDate().toString().padStart(2, '0');
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const year = date.getFullYear().toString().slice(-2);

    // Return the formatted date string
    return `${day}.${month}.${year}`;
}
</script>

<div 
	class="grid grid-cols-[1fr_1fr] items-center gap-y-4 p-2 py-4 shadow sm:grid-cols-[1fr_1fr_auto] sm:p-4"
>
	<div hidden={showEdit}>
		<h2 class="pb-2 text-xl font-bold">{plant.name}</h2>
		<p class="text-zinc-800">{plant.place}</p>
	</div>
	<div hidden={!showEdit}>
		<input bind:value={newName} type="text" class="pb-2 text-xl font-bold"/>
		<input bind:value={newPlace} type="text" class="text-zinc-800"/>
	</div>
	<div class="flex items-baseline justify-end gap-2 text-lg sm:justify-start">
		<svg
			class="h-4"
			viewBox="0 0 5.2922349 4.4979205"
			version="1.1"
			id="svg1"
			xmlns="http://www.w3.org/2000/svg"
		>
			<defs id="defs1" />
			<g id="layer1" transform="translate(-164.33876,-150.79365)">
				<path
					fill="currentColor"
					d="m 168.70495,151.71176 c -0.19579,0.19314 -0.24606,0.47096 -0.15346,0.70908 l -0.77258,0.77258 v -0.54768 c 0,-0.14552 -0.11907,-0.26459 -0.26459,-0.26459 h -0.27252 c 0.008,-0.045 0.008,-0.0873 0.008,-0.13229 0,-0.80433 -0.65087,-1.45521 -1.45521,-1.45521 a 1.4544146,1.4544146 0 0 0 -0.9261,2.57704 v 1.6563 c 0,0.14552 0.11906,0.26458 0.26458,0.26458 h 2.38125 c 0.14552,0 0.26459,-0.11906 0.26459,-0.26458 v -1.08744 l 1.14564,-1.14565 c 0.23813,0.0926 0.51594,0.045 0.70644,-0.14816 z m -3.82323,0.66939 c -0.005,-0.045 -0.0132,-0.0873 -0.0132,-0.13229 0,-0.51064 0.4154,-0.92604 0.92604,-0.92604 0.51065,0 0.92604,0.4154 0.92604,0.92604 0,0.045 -0.008,0.0873 -0.0132,0.13229 z"
					id="path1"
					style="stroke-width:0.264583"
				/>
			</g>
		</svg>

		<p>
			{lastWateredDate}
		</p>
	</div>
	<div
		class="w-full max-sm:col-span-2 max-sm:grid max-sm:grid-cols-[1fr_auto] max-sm:gap-2 max-sm:place-self-end"
	>
		<button hidden={showEdit} onclick={waterNow}  class="rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-500">Water now</button>
		<button hidden={!showEdit} onclick={postPlant}  class="rounded-lg bg-green-600 px-4 py-2 text-white hover:bg-green-500">Save</button>
		<button hidden={showEdit} onclick={() => showEdit = true} class="rounded-lg bg-zinc-300 px-4 py-2 hover:bg-zinc-200">Edit</button>
		<button hidden={!showEdit} onclick={() => {showEdit = false; newName = plant.name; newPlace=plant.place}} class="rounded-lg bg-zinc-300 px-4 py-2 hover:bg-zinc-200">Cancel</button>
	</div>

	<button hidden={!showEdit} onclick={deletePlant} class="bg-red-700 hover:bg-red-600 px-4 py-2 rounded-lg col-start-3 text-white">Delete Plant</button>
</div>
