import { error, json} from '@sveltejs/kit';

import type { Plant } from '$lib/types.js';

export function GET({ url }) {
	const min = Number(url.searchParams.get('min') ?? '0');
	const max = Number(url.searchParams.get('max') ?? '1');

	const d = max - min;

	if (isNaN(d) || d < 0) {
		error(400, 'min and max must be numbers, and min must be less than max');
	}

	const random = min + Math.random() * d;
	return new Response(String(random));
}

export async function POST({request}) {
    const {name, place} = await request.json();
    
    const plant: Plant = {
        id: 0,
        name: name, 
        place: place,
        lastWateredAt: new Date().toISOString()
    }

    const response  = await fetch(`http://localhost:8000/api/go/plants`, {
        method: 'POST',
        body: JSON.stringify({...plant}),
        headers: {
            'content-type': 'application/json'
        }
    })

    return json("Success")
}

export async function PUT({request}) {
    const {id, name, place, lastWateredAt} = await request.json();
    
    const response  = await fetch(`http://localhost:8000/api/go/plants`, {
        method: 'PUT',
        body: JSON.stringify({id, name, place, lastWateredAt}),
        headers: {
            'content-type': 'application/json'
        }
    })

    return json("Success")
}

export async function DELETE({request}) {
    const {id, name, place, lastWateredAt} = await request.json();
    
    const response  = await fetch(`http://localhost:8000/api/go/plants/${id}`, {
        method: 'DELETE',
    })

    return json("Success")
}