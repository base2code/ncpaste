import { Router } from 'itty-router';

const API_URL = 'https://europe-west3-ncpaste.cloudfunctions.net/ncpaste-func';

// now let's create a router (note the lack of "new")
const router = Router();

// GET item
router.get('/', async (req) => {
	const url = req.url;
	const {pathname} = new URL(url);
	return await fetch(`${API_URL}`);
});

// GET collection index
router.get('/*', async (req) => {
	const url = req.url.replace("/raw", "");
	const {pathname} = new URL(url);
	return await fetch(`${API_URL}${pathname}`);
});

// POST with hastebin endpoint
router.post('/documents', async (request) => {
	const body = await request.text();
	const code = await fetch(`${API_URL}`, {
		method: 'POST',
		body,
	});
	// the body of the response will be a JSON string. Turn it into an json return. Set the value of "key"
	return new Response(JSON.stringify({key: await code.text()}), {
		headers: {
			'Content-Type': 'application/json',
		},
	});

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
