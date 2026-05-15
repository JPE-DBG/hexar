import { test, expect, type Browser } from '@playwright/test';
import { createAndEnterGame, joinGame, waitForGameActive } from './helpers';

async function setupGame(browser: Browser) {
  const ctx1 = await browser.newContext();
  const ctx2 = await browser.newContext();
  const p1 = await ctx1.newPage();
  const p2 = await ctx2.newPage();
  const code = await createAndEnterGame(p1);
  await joinGame(p2, code);
  await waitForGameActive(p1);
  await waitForGameActive(p2);
  return { ctx1, ctx2, p1, p2 };
}

test('both players see canvas and HUD after game starts', async ({ browser }) => {
  const { ctx1, ctx2, p1, p2 } = await setupGame(browser);
  try {
    await expect(p1.locator('#game')).toBeVisible();
    await expect(p1.locator('#hud')).toContainText('Player');
    await expect(p2.locator('#game')).toBeVisible();
    await expect(p2.locator('#hud')).toContainText('Player');
  } finally {
    await ctx1.close();
    await ctx2.close();
  }
});

test('gold counter increases as economy ticks run', async ({ browser }) => {
  const { ctx1, ctx2, p1 } = await setupGame(browser);
  try {
    const getGold = async () => {
      const text = await p1.locator('#hud .hud-value').first().textContent();
      return parseInt(text ?? '0', 10);
    };
    const before = await getGold();
    await p1.waitForTimeout(2_000);
    const after = await getGold();
    expect(after).toBeGreaterThan(before);
  } finally {
    await ctx1.close();
    await ctx2.close();
  }
});

test('sidebar is visible during active game', async ({ browser }) => {
  const { ctx1, ctx2, p1 } = await setupGame(browser);
  try {
    await expect(p1.locator('#sidebar')).toBeVisible();
  } finally {
    await ctx1.close();
    await ctx2.close();
  }
});
