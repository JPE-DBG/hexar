import { test, expect } from '@playwright/test';
import { createAndEnterGame, joinGame, waitForGameActive } from './helpers';

test('duplicate tab shows already-connected overlay', async ({ browser }) => {
  const ctx1 = await browser.newContext();
  const ctx2 = await browser.newContext();
  const ctx3 = await browser.newContext();
  try {
    const p1 = await ctx1.newPage();
    const p2 = await ctx2.newPage();

    const code = await createAndEnterGame(p1);
    await joinGame(p2, code);
    await waitForGameActive(p1);

    // P1's URL now has /#CODE:TOKEN written by history.replaceState
    const p1Url = p1.url();
    expect(p1Url).toContain('#');

    // Open same session URL in a third context — server rejects with 4001
    const p3 = await ctx3.newPage();
    await p3.goto(p1Url);
    await expect(p3.getByText('Already Connected')).toBeVisible({ timeout: 10_000 });
  } finally {
    await ctx1.close();
    await ctx2.close();
    await ctx3.close();
  }
});

test('URL hash reconnects after tab close', async ({ browser }) => {
  const ctx1 = await browser.newContext();
  const ctx2 = await browser.newContext();
  try {
    const p1 = await ctx1.newPage();
    const p2 = await ctx2.newPage();

    const code = await createAndEnterGame(p1);
    await joinGame(p2, code);
    await waitForGameActive(p1);

    const p1Url = p1.url();
    expect(p1Url).toContain('#');

    // Close P1 tab (disconnects), then reopen with hash URL
    await p1.close();

    const p1b = await ctx1.newPage();
    await p1b.goto(p1Url);
    // Hash lets main.ts call startGame directly — no lobby shown, game resumes
    await waitForGameActive(p1b);
    await expect(p1b.locator('#hud')).toContainText('Gold');
  } finally {
    await ctx1.close();
    await ctx2.close();
  }
});

test('pause banner appears when opponent disconnects', async ({ browser }) => {
  const ctx1 = await browser.newContext();
  const ctx2 = await browser.newContext();
  try {
    const p1 = await ctx1.newPage();
    const p2 = await ctx2.newPage();

    const code = await createAndEnterGame(p1);
    await joinGame(p2, code);
    await waitForGameActive(p1);
    await waitForGameActive(p2);

    await p2.close();

    await expect(p1.locator('.banner')).toContainText('Game paused', { timeout: 5_000 });
  } finally {
    await ctx1.close();
    await ctx2.close();
  }
});
