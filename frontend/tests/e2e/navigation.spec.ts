import { test, expect } from '@playwright/test';

test.describe('Navigation & UI', () => {
  test('should display 404 page for invalid routes', async ({ page }) => {
    // Navigate to a non-existent page
    await page.goto('/this-page-does-not-exist-12345');
    await page.waitForLoadState('networkidle');

    // Verify 404 page is shown - check for the h1 with "404"
    await expect(page.locator('h1').filter({ hasText: '404' })).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('Page Not Found')).toBeVisible();
  });

  test('should navigate home from 404 page', async ({ page }) => {
    // Navigate to a non-existent page
    await page.goto('/invalid-route-xyz');
    await page.waitForLoadState('networkidle');

    // Wait for 404 page
    await expect(page.locator('h1').filter({ hasText: '404' })).toBeVisible({ timeout: 10000 });

    // Click Return Home button
    const returnHomeButton = page.getByRole('button', { name: /Return Home/i });
    await returnHomeButton.click();

    // Should navigate to home
    await expect(page).toHaveURL('/');
  });

  test('home page CTAs should navigate correctly', async ({ page }) => {
    // Go to home page
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Check for Get Started button
    const getStartedButton = page.getByRole('link', { name: /get started/i });
    if (await getStartedButton.isVisible({ timeout: 3000 }).catch(() => false)) {
      await getStartedButton.click();
      await expect(page).toHaveURL(/\/(register|login)/);
    }
  });

  test('header navigation links work correctly', async ({ page }) => {
    // First register/login a user to see full navigation
    const suffix = Math.random().toString(36).substring(7);
    const user = {
      username: `navtest_${suffix}`,
      email: `navtest_${suffix}@example.com`,
      password: 'securePassword123!',
    };

    await page.goto('/register');
    await page.fill('input#username', user.username);
    await page.fill('input#email', user.email);
    await page.fill('input#password', user.password);
    await page.fill('input#confirmPassword', user.password);
    await page.locator('input[type="checkbox"]').check();
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 10000 });

    // Test navigation to home via logo/brand
    const logoLink = page.locator('header a').first();
    await logoLink.click();
    await expect(page).toHaveURL('/');

    // Navigate back to dashboard
    await page.goto('/dashboard');
    await page.waitForLoadState('networkidle');

    // Verify dashboard is accessible
    await expect(page.getByText(/Dashboard|Welcome/i)).toBeVisible({ timeout: 5000 });
  });
});

test.describe('Mobile Navigation', () => {
  test.use({ viewport: { width: 375, height: 667 } }); // iPhone SE size

  test('mobile menu should toggle on small screens', async ({ page }) => {
    // Go to home page on mobile viewport
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Look for mobile menu button (hamburger icon) - uses "Toggle menu" aria-label
    const menuButton = page.locator('button[aria-label="Toggle menu"]');

    if (await menuButton.isVisible({ timeout: 3000 }).catch(() => false)) {
      // Click to open menu
      await menuButton.click();

      // Wait for menu animation
      await page.waitForTimeout(300);

      // Verify mobile navigation is visible - look for the md:hidden nav
      const mobileNav = page.locator('nav.md\\:hidden');
      await expect(mobileNav).toBeVisible({ timeout: 3000 });

      // Click menu button again to close
      await menuButton.click();
      await page.waitForTimeout(300);
    }
  });

  test('mobile menu links should navigate and close menu', async ({ page }) => {
    // Go to home page on mobile viewport
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Look for mobile menu button
    const menuButton = page.locator('button[aria-label="Toggle menu"]');

    if (await menuButton.isVisible({ timeout: 3000 }).catch(() => false)) {
      // Open menu
      await menuButton.click();
      await page.waitForTimeout(300);

      // Look for Sign In link in mobile menu (inside the md:hidden nav)
      const mobileNav = page.locator('nav.md\\:hidden');
      const signInLink = mobileNav.getByRole('link', { name: /sign in/i });

      if (await signInLink.isVisible({ timeout: 3000 }).catch(() => false)) {
        await signInLink.click();

        // Should navigate to login
        await expect(page).toHaveURL(/\/login/, { timeout: 5000 });
      }
    }
  });

  test('authenticated mobile menu shows Sign Out', async ({ page }) => {
    // Register a user
    const suffix = Math.random().toString(36).substring(7);
    const user = {
      username: `mobile_${suffix}`,
      email: `mobile_${suffix}@example.com`,
      password: 'securePassword123!',
    };

    await page.goto('/register');
    await page.fill('input#username', user.username);
    await page.fill('input#email', user.email);
    await page.fill('input#password', user.password);
    await page.fill('input#confirmPassword', user.password);
    await page.locator('input[type="checkbox"]').check();
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 10000 });

    // On mobile, Sign Out button should be visible in header
    const signOutButton = page.getByRole('button', { name: /sign out/i });
    await expect(signOutButton).toBeVisible({ timeout: 5000 });
  });
});
