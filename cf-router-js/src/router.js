import { Router } from 'itty-router';

const API_URL = 'https://europe-west3-ncpaste.cloudfunctions.net/ncpaste-func';

// now let's create a router (note the lack of "new")
const router = Router();

// GET item
router.get('/', () => {
	return new Response('INDEX');
});

// GET collection index
router.get('/*', async (req) => {
	const {pathname} = new URL(req.url);
	return await fetch(`${API_URL}${pathname}`);
});

// POST to the collection (we'll use async here)
router.post('/', async (request) => {
	const body = await request.text();
	return await fetch(`${API_URL}`, {
		method: 'POST',
		body,
	});
});

// 404 for everything else
router.all('*', () => new Response('Not Found.', { status: 404 }));

export default router;
