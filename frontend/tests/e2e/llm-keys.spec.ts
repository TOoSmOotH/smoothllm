import { test, expect } from '@playwright/test';

// Create a unique user for API key tests
const randomSuffix = Math.random().toString(36).substring(7);
const testUser = {
  username: `keyuser_${randomSuffix}`,
  email: `key_${randomSuffix}@example.com`,
  password: 'securePassword123!',
};

const testProvider = {
  name: `Key Test Provider ${randomSuffix}`,
  apiKey: 'sk-test-fake-key-for-testing',
};

const testApiKey = {
  name: `My Test Key ${randomSuffix}`,
};

// Helper to register and login
async function registerAndLogin(page: any) {
  await page.goto('/register');
  await page.fill('input#username', testUser.username);
  await page.fill('input#email', testUser.email);
  await page.fill('input#password', testUser.password);
  await page.fill('input#confirmPassword', testUser.password);
  await page.locator('input[type="checkbox"]').check();
  await page.click('button[type="submit"]');

  await page.waitForURL(/\/(dashboard|login)/, { timeout: 15000 });

  if (page.url().includes('/login')) {
    await page.fill('input#email', testUser.email);
    await page.fill('input#password', testUser.password);
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/dashboard/, { timeout: 10000 });
  }
}

async function loginUser(page: any) {
  await page.goto('/login');
  await page.fill('input#email', testUser.email);
  await page.fill('input#password', testUser.password);
  await page.click('button[type="submit"]');
  await page.waitForURL(/\/dashboard/, { timeout: 10000 });
}

async function createProvider(page: any) {
  await page.goto('/providers');
  await page.waitForLoadState('networkidle');

  // Wait for page to be ready
  await expect(page.locator('h1:has-text("LLM Providers")')).toBeVisible({ timeout: 10000 });

  const addButton = page.getByRole('button', { name: /Add Provider/i }).first();
  await expect(addButton).toBeVisible({ timeout: 5000 });
  await addButton.click();

  // Wait for modal
  const providerNameInput = page.locator('input[placeholder="My OpenAI Account"]');
  await expect(providerNameInput).toBeVisible({ timeout: 10000 });

  await providerNameInput.fill(testProvider.name);
  const apiKeyInput = page.locator('input[placeholder="sk-..."]');
  await apiKeyInput.fill(testProvider.apiKey);
  await page.getByRole('button', { name: 'Create Provider' }).click();

  await expect(page.getByText('Provider created successfully')).toBeVisible({ timeout: 10000 });
}

test.describe('API Keys Management', () => {
  test.describe.configure({ mode: 'serial' });

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    try {
      await registerAndLogin(page);
      await createProvider(page);
    } finally {
      await page.close();
    }
  });

  test('should display API keys page', async ({ page }) => {
    await loginUser(page);

    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    await expect(page).toHaveURL(/\/keys/);
    // Use text selector which is more reliable than role selector
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 15000 });
    await expect(page.getByText('Manage your proxy API keys')).toBeVisible();
  });

  test('should show empty state when no keys exist', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Should show empty state or keys list
    const emptyStateText = page.getByText('No API Keys');
    const hasEmptyState = await emptyStateText.isVisible({ timeout: 5000 }).catch(() => false);
    if (hasEmptyState) {
      await expect(page.getByRole('button', { name: 'Create Your First Key' })).toBeVisible();
    }
  });

  test('should open create key modal', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Click Create Key button
    const createButton = page.getByRole('button', { name: /Create Key/i }).first();
    await expect(createButton).toBeVisible({ timeout: 5000 });
    await createButton.click();

    // Modal should be visible - check for form fields using placeholder
    const keyNameInput = page.locator('input[placeholder="My API Key"]');
    await expect(keyNameInput).toBeVisible({ timeout: 10000 });
  });

  test('should require provider selection when creating key', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    const createButton = page.getByRole('button', { name: /Create Key/i }).first();
    await expect(createButton).toBeVisible({ timeout: 5000 });
    await createButton.click();

    // Wait for modal using placeholder
    const keyNameInput = page.locator('input[placeholder="My API Key"]');
    await expect(keyNameInput).toBeVisible({ timeout: 10000 });

    // Provider dropdown should exist
    const providerSelect = page.locator('select').first();
    await expect(providerSelect).toBeVisible();

    // Our test provider should be in the dropdown
    const options = await providerSelect.locator('option').allTextContents();
    const hasProvider = options.some(opt => opt.includes(testProvider.name));
    expect(hasProvider).toBe(true);
  });

  test('should create a new API key', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Wait a moment for the store to initialize with providers
    await page.waitForTimeout(1000);

    const createButton = page.getByRole('button', { name: /Create Key/i }).first();
    await expect(createButton).toBeVisible({ timeout: 5000 });
    await expect(createButton).toBeEnabled({ timeout: 5000 });
    await createButton.click();

    // Wait for modal using placeholder
    const keyNameInput = page.locator('input[placeholder="My API Key"]');
    await expect(keyNameInput).toBeVisible({ timeout: 10000 });

    // Fill in the form
    await keyNameInput.fill(testApiKey.name);

    // Wait for provider dropdown to populate
    const providerSelect = page.locator('select').first();
    await expect(providerSelect).toBeVisible({ timeout: 5000 });

    // Wait for options to be populated (more than just the placeholder)
    await page.waitForFunction(
      () => {
        const select = document.querySelector('select');
        if (!select) return false;
        const options = select.querySelectorAll('option:not([disabled])');
        return options.length > 0;
      },
      { timeout: 10000 }
    );

    // Select the first available (non-disabled) provider by index
    // Index 0 is the disabled "Select a provider..." placeholder
    await providerSelect.selectOption({ index: 1 });

    // Submit - click the button inside the modal (not the one in the page header)
    // The modal is teleported to body, so find the modal container first
    const modalSubmitButton = page.locator('.fixed.inset-0').getByRole('button', { name: 'Create Key' });
    await modalSubmitButton.click();

    // New key modal should appear
    await expect(page.getByRole('heading', { name: 'API Key Created' })).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('Copy your key now')).toBeVisible();

    // Key should be displayed
    await expect(page.locator('code')).toContainText('sk-smoothllm-');
  });

  test('should show warning about saving the key', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Create another key to test the warning
    const createButton = page.getByRole('button', { name: /Create Key/i }).first();
    await expect(createButton).toBeVisible({ timeout: 5000 });
    await createButton.click();

    // Wait for modal using placeholder
    const keyNameInput = page.locator('input[placeholder="My API Key"]');
    await expect(keyNameInput).toBeVisible({ timeout: 10000 });

    await keyNameInput.fill(`Another Key ${randomSuffix}`);
    const providerSelect = page.locator('select').first();

    // Wait for options to be populated
    await page.waitForFunction(
      () => {
        const select = document.querySelector('select');
        if (!select) return false;
        const options = select.querySelectorAll('option:not([disabled])');
        return options.length > 0;
      },
      { timeout: 10000 }
    );

    // Select the first available provider by index
    await providerSelect.selectOption({ index: 1 });

    // Submit - click the button inside the modal
    const modalSubmitButton = page.locator('.fixed.inset-0').getByRole('button', { name: 'Create Key' });
    await modalSubmitButton.click();

    await expect(page.getByRole('heading', { name: 'API Key Created' })).toBeVisible({ timeout: 10000 });

    // Warning should be visible
    await expect(page.getByText('Save this key securely')).toBeVisible();
    await expect(page.getByText('only time you\'ll see this key')).toBeVisible();

    // Close the modal
    await page.getByRole('button', { name: 'Done' }).click();
  });

  test('should display created keys in the list', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Our key should be in the list
    const keyHeading = page.locator(`h3:has-text("${testApiKey.name}")`);
    await expect(keyHeading).toBeVisible({ timeout: 10000 });

    // Should show key prefix
    await expect(page.locator('text=sk-smoothllm-').first()).toBeVisible();

    // Should show provider name
    await expect(page.getByText(testProvider.name).first()).toBeVisible();

    // Should show Active status
    await expect(page.getByText('Active').first()).toBeVisible();
  });

  test('should open edit modal for existing key', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Wait for our key heading to be visible
    const keyHeading = page.locator(`h3:has-text("${testApiKey.name}")`);
    await expect(keyHeading).toBeVisible({ timeout: 10000 });

    // Find the key card container and click the edit button (first button in the card)
    const keyCard = page.locator('div.rounded-lg').filter({ has: keyHeading }).first();
    await keyCard.locator('button').first().click();

    // Modal should open with edit title
    await expect(page.getByRole('heading', { name: 'Edit API Key' })).toBeVisible({ timeout: 5000 });

    // Name field should be pre-filled
    const keyNameInput = page.locator('input[placeholder="My API Key"]');
    await expect(keyNameInput).toHaveValue(testApiKey.name);
  });

  test('should update key name', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Wait for our key heading
    const keyHeading = page.locator(`h3:has-text("${testApiKey.name}")`);
    await expect(keyHeading).toBeVisible({ timeout: 10000 });

    // Open edit modal
    const keyCard = page.locator('div.rounded-lg').filter({ has: keyHeading }).first();
    await keyCard.locator('button').first().click();
    await expect(page.getByRole('heading', { name: 'Edit API Key' })).toBeVisible({ timeout: 5000 });

    // Update name
    const updatedName = `${testApiKey.name} Updated`;
    const keyNameInput = page.locator('input[placeholder="My API Key"]');
    await keyNameInput.clear();
    await keyNameInput.fill(updatedName);

    // Submit
    await page.getByRole('button', { name: 'Save Changes' }).click();

    // Wait for success
    await expect(page.getByText('API key updated successfully')).toBeVisible({ timeout: 10000 });

    // Updated name should appear
    await expect(page.getByText(updatedName)).toBeVisible({ timeout: 5000 });

    // Update the test key name for later tests
    testApiKey.name = updatedName;
  });

  test('should show revoke confirmation modal', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Wait for our key heading
    const keyHeading = page.locator(`h3:has-text("${testApiKey.name}")`);
    await expect(keyHeading).toBeVisible({ timeout: 10000 });

    // Find an active key and click revoke button (second button, after edit)
    const keyCard = page.locator('div.rounded-lg').filter({ has: keyHeading }).first();
    await keyCard.locator('button').nth(1).click();

    // Confirmation modal should appear
    await expect(page.getByRole('heading', { name: 'Revoke API Key' })).toBeVisible({ timeout: 5000 });
    await expect(page.getByText('immediately stop this key from working')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Revoke Key' })).toBeVisible();
  });

  test('should revoke an API key', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Find an active key - look for the Active status badge in an h3's sibling
    const activeStatus = page.getByText('Active').first();
    await expect(activeStatus).toBeVisible({ timeout: 10000 });

    // Click revoke on the first active key - find the card containing Active status
    const activeKeyCard = page.locator('div.rounded-lg').filter({ hasText: 'Active' }).first();
    await activeKeyCard.locator('button').nth(1).click();

    await expect(page.getByRole('heading', { name: 'Revoke API Key' })).toBeVisible({ timeout: 5000 });

    // Confirm revocation - target button in the modal
    const modalDialog = page.locator('.relative.bg-bg-primary.rounded-lg');
    await modalDialog.getByRole('button', { name: 'Revoke Key' }).click();

    // Success message
    await expect(page.getByText('API key revoked successfully')).toBeVisible({ timeout: 10000 });

    // Key should now show as Revoked
    await expect(page.getByText('Revoked').first()).toBeVisible({ timeout: 5000 });
  });

  test('should show delete confirmation modal', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Wait for our key heading
    const keyHeading = page.locator(`h3:has-text("${testApiKey.name}")`);
    await expect(keyHeading).toBeVisible({ timeout: 10000 });

    // Click delete button (last action button)
    const keyCard = page.locator('div.rounded-lg').filter({ has: keyHeading }).first();
    await keyCard.locator('button').last().click();

    // Confirmation modal should appear
    await expect(page.getByRole('heading', { name: 'Delete API Key' })).toBeVisible({ timeout: 5000 });
    await expect(page.getByText('This action cannot be undone')).toBeVisible();
  });

  test('should cancel delete when clicking cancel', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Wait for our key heading
    const keyHeading = page.locator(`h3:has-text("${testApiKey.name}")`);
    await expect(keyHeading).toBeVisible({ timeout: 10000 });

    // Click delete button
    const keyCard = page.locator('div.rounded-lg').filter({ has: keyHeading }).first();
    await keyCard.locator('button').last().click();

    await expect(page.getByRole('heading', { name: 'Delete API Key' })).toBeVisible({ timeout: 5000 });

    // Click cancel
    await page.getByRole('button', { name: 'Cancel' }).click();

    // Modal should close
    await expect(page.getByRole('heading', { name: 'Delete API Key' })).not.toBeVisible({ timeout: 3000 });

    // Key should still be visible
    await expect(page.getByText(testApiKey.name)).toBeVisible();
  });

  test('should delete an API key', async ({ page }) => {
    await loginUser(page);
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Wait for our key heading
    const keyHeading = page.locator(`h3:has-text("${testApiKey.name}")`);
    await expect(keyHeading).toBeVisible({ timeout: 10000 });

    // Get the name of the key we're about to delete for verification
    const keyCard = page.locator('div.rounded-lg').filter({ has: keyHeading }).first();
    await keyCard.locator('button').last().click();

    await expect(page.getByRole('heading', { name: 'Delete API Key' })).toBeVisible({ timeout: 5000 });

    // Confirm deletion
    await page.getByRole('button', { name: 'Delete Key' }).click();

    // Success message
    await expect(page.getByText('API key deleted successfully')).toBeVisible({ timeout: 10000 });

    // Key should no longer be visible (or empty state if last key)
    await page.waitForLoadState('networkidle');
  });
});

test.describe('API Keys - No Provider Warning', () => {
  test('should show warning when no providers exist', async ({ page }) => {
    // Create a fresh user with no providers
    const suffix = Math.random().toString(36).substring(7);
    const user = {
      username: `nokey_${suffix}`,
      email: `nokey_${suffix}@example.com`,
      password: 'securePassword123!',
    };

    await page.goto('/register');
    await page.fill('input#username', user.username);
    await page.fill('input#email', user.email);
    await page.fill('input#password', user.password);
    await page.fill('input#confirmPassword', user.password);
    await page.locator('input[type="checkbox"]').check();
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/(dashboard|login)/, { timeout: 15000 });

    if (page.url().includes('/login')) {
      await page.fill('input#email', user.email);
      await page.fill('input#password', user.password);
      await page.click('button[type="submit"]');
      await page.waitForURL(/\/dashboard/, { timeout: 10000 });
    }

    // Go to keys page without creating providers
    await page.goto('/keys');
    await page.waitForLoadState('networkidle');

    // Wait for page to load
    await expect(page.locator('h1:has-text("API Keys")')).toBeVisible({ timeout: 10000 });

    // Should show warning about no providers
    await expect(page.getByText('No providers configured')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('configure at least one LLM provider')).toBeVisible();

    // Create Key button should be disabled
    const createButton = page.getByRole('button', { name: /Create Key/i });
    await expect(createButton).toBeDisabled();
  });
});
