# Simple URL Shortener

This is a simple URL shortener service written in Go.

## What does it do?
- Lets you create short links for long URLs.
- You can set how long the short link should work (default is 30 minutes).
- You can use your own shortcode, or the service will make one for you.
- You can see stats for each short link (like how many times it was clicked).
- When someone visits the short link, they get sent to the original URL, and the click is counted.

## How to run it

1. **Clone the repo and go to the folder:**

2. **Set up your .env file:**
   - Add your PORT (for example, `PORT=8080`)
   - Add log API settings.

3. **Build and run command:**
   make run
 
   The server will start. You can now use the API.

## API Endpoints

- POST /shorturls — Create a new short URL
- GET /shorturls/{shortcode} — Get stats for a short URL
- GET /{shortcode} — Redirect to the original URL and count the click

## How it works
- All data is stored in memory.
- There is a simple logging system for requests and errors.
- The code is split into files for main logic, handlers, storage, and utilities.

---
 