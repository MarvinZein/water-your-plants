import type { PageLoad } from "./$types"

const apiUrl = 'http://localhost:8000/api/go'

export const load: PageLoad = async ({ fetch }) => {
    const plantsResponse = await fetch(apiUrl + '/plants');
    if (!plantsResponse.ok) {
        throw new Error('Failed to fetch plants');
    }

    const { plants } = await plantsResponse.json();

    return {
        plants: plants
    };
};
