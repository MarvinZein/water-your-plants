import type { PageLoad } from "./$types"

import { apiURL } from "$lib";

export const load: PageLoad = async ({ fetch }) => {
    const plantsResponse = await fetch(apiURL + '/plants');
    if (!plantsResponse.ok) {
        throw new Error('Failed to fetch plants');
    }

    const { plants } = await plantsResponse.json();

    return {
        plants: plants
    };
};
