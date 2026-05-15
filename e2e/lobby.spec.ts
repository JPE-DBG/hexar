import { test, expect } from '@playwright/test';
import { createAndEnterGame, joinGame, waitForGameActive } from './helpers';

test('create game shows 4-character room code', async ({ page }) => {
  await page.goto('/');
  await page.click('#createBtn');
  await page.waitForSelector('#displayCode');
  const code = await page.textContent('#displayCode');
  expect(code?.trim()).toMatch(/^[A-Z0-9]{4}$/);
});

test('join with invalid code shows error message', async ({ page }) => {
  await page.goto('/');
  await page.fill('#codeInput', 'XXXX');
  await page.click('#joinBtn');
  await page.waitForFunction(() => {
    const el = document.getElementById('lobbyStatus');
    return el !== null && el.textContent !== null && el.textContent.length > 0;
  });
  const status = await page.textContent('#lobbyStatus');
  expect(status?.trim()).toBeTruthy();
});

test('waiting overlay appears after creator enters game', async ({ page }) => {
  await createAndEnterGame(page);
  await expect(page.getByText('Waiting for opponent')).toBeVisible();
});

test('HUD activates when second player joins', async ({ browser }) => {
  const ctx1 = await browser.newContext();
  const ctx2 = await browser.newContext();
  try {
    const p1 = await ctx1.newPage();
    const p2 = await ctx2.newPage();

    const code = await createAndEnterGame(p1);
    await joinGame(p2, code);

    await waitForGameActive(p1);
    await waitForGameActive(p2);

    await expect(p1.locator('#hud')).toContainText('Gold');
    await expect(p2.locator('#hud')).toContainText('Gold');
  } finally {
    await ctx1.close();
    await ctx2.close();
  }
});
