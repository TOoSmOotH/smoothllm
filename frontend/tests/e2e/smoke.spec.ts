import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/', { waitUntil: 'networkidle' });

  // Wait for the app to be mounted and visible
  await expect(page.locator('#app').first()).toBeVisible({ timeout: 10000 });

  // Check title with a longer timeout and flexible match
  await expect(page).toHaveTitle(/SmoothLLM|SmoothWeb|Vite/i, { timeout: 10000 });

  // Verify some content is present on the home page
  await expect(page.getByRole('heading', { level: 1 })).toContainText(/Smooth/i);
});

test('backend health check', async ({ request }) => {
  const response = await request.get('http://backend:8080/health');
  expect(response.status()).toBe(200);
});
