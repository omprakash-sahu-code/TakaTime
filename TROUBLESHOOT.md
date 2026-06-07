# Troubleshooting Guide

This document covers common issues encountered while installing, configuring, and using TakaTime.

## 1. Extension Not Tracking Coding Activity

### Symptoms

* No data appears in MongoDB.
* Dashboard remains empty.
* Coding activity is not recorded.

### Possible Causes

* MongoDB URI is not configured.
* Helper binary failed to download.
* Extension is disabled.

### Solutions

#### Verify MongoDB Configuration

Open Command Palette:

```text
Ctrl+Shift+P
```

Run:

```text
TakaTime: Setup MongoDB URI
```

Ensure the connection string is valid:

```text
mongodb+srv://username:password@cluster.mongodb.net/
```

#### Check Extension Status

VS Code:

```text
Extensions → TakaTime → Enabled
```

Restart VS Code after enabling.

---

## 2. MongoDB Connection Failed

### Symptoms

```text
Failed to connect to MongoDB
Authentication failed
```

### Solutions

#### Verify Database User

MongoDB Atlas:

```text
Database Access → Users
```

Confirm:

* Username exists
* Password is correct
* User has read/write permissions

#### Verify Network Access

MongoDB Atlas:

```text
Network Access
```

Allow:

```text
0.0.0.0/0
```

or your specific IP address.

#### Test Connection

Use MongoDB Compass:

```text
Connect using the same URI
```

If Compass cannot connect, the issue is with MongoDB configuration rather than TakaTime.

---

## 3. Binary Download Failed

### Symptoms

```text
Failed to download helper binary
Binary not found
```

### Solutions

#### Check Internet Connection

TakaTime downloads platform-specific binaries from GitHub Releases.

Verify:

```bash
ping github.com
```

#### Antivirus Blocking Download

Some antivirus software may quarantine downloaded binaries.

Check:

```text
Windows Security → Protection History
```

Restore the binary if blocked.

#### Reinstall Binary

Run:

```text
TakaTime: Setup MongoDB URI
```

This triggers a binary verification/download process.

---

## 4. GitHub README Stats Not Updating

### Symptoms

* README graphs do not refresh.
* GitHub Action fails.

### Solutions

#### Verify Workflow Exists

```text
.github/workflows/takatime.yml
```

#### Check GitHub Actions

Navigate to:

```text
Repository → Actions
```

Inspect failed workflow logs.

#### Verify MongoDB Credentials

Ensure all required secrets are configured:

```text
MONGODB_URI
```

and any additional secrets required by the workflow.

---

## 5. Activity Missing After Offline Coding

### Symptoms

* Activity recorded while offline does not appear.

### Solutions

#### Reconnect to Internet

TakaTime supports offline caching and syncs automatically once connectivity is restored.

Wait a few minutes after reconnecting.

#### Restart Editor

Restart:

* VS Code
* Neovim

This can trigger synchronization.

---

## 6. Invalid MongoDB URI

### Symptoms

```text
Invalid connection string
```

### Common Mistakes

Incorrect:

```text
mongodb://cluster.mongodb.net
```

Correct:

```text
mongodb+srv://username:password@cluster.mongodb.net/
```

URL-encode special characters in passwords:

```text
@
#
$
%
```

Example:

```text
my%40password
```

---

## 7. No Data Appearing in Atlas

### Check Collection Creation

Open MongoDB Atlas:

```text
Database → Collections
```

Expected:

```text
takatime
```

or project-specific collections created by the tracker.

### Verify Writes

Monitor:

```text
Atlas → Metrics
```

for incoming operations.

---

## 8. GitHub Action Permission Errors

### Symptoms

```text
Permission denied
403 Forbidden
```

### Solution

Repository Settings:

```text
Actions → General
```

Enable:

```text
Read and write permissions
```

for GitHub Actions.

---

## 9. Neovim Plugin Not Working

### Verify Installation

Check:

```lua
require("taka-time").setup()
```

loads successfully.

### Run Initialization

Inside Neovim:

```vim
:TakaInit
```

### Check Plugin Manager

Ensure the plugin is installed and loaded by:

* lazy.nvim
* packer.nvim
* vim-plug

---

## 10. Debugging Checklist

Before opening an issue, verify:

* [ ] MongoDB URI is valid.
* [ ] Atlas network access is configured.
* [ ] Database user has permissions.
* [ ] Latest TakaTime version is installed.
* [ ] GitHub Actions workflow is configured correctly.
* [ ] Internet connection is available.
* [ ] Editor has been restarted.
* [ ] Logs and error messages are included in the issue report.

## Getting Help

If the issue persists:

1. Check existing GitHub Issues.
2. Collect logs/screenshots.
3. Open a new issue with:

   * Operating System
   * Editor version
   * TakaTime version
   * Error message
   * Steps to reproduce
