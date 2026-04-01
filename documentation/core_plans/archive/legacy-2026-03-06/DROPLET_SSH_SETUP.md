# DigitalOcean Droplet SSH Setup Guide

This guide covers how to set up SSH key-based authentication for your DigitalOcean Droplet while retaining the ability to log in with your root password. This is useful for managing the InsureTech backend server.

## Step 1: Generate an SSH Key Pair (Local Machine)
If you don't already have an SSH key, generate one on your local Windows machine.

> [!NOTE]
> **Getting an error that `ssh-keygen` is not recognized?**
> You need to enable the OpenSSH Client in Windows:
> 1. Right-click the Start button and open **Terminal (Admin)** or **Windows PowerShell (Admin)**.
> 2. Run: `Add-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0`
> 3. Restart your normal terminal (e.g., VS Code) and try again.

Open **PowerShell** or Command Prompt and run:

```powershell
ssh-keygen -t ed25519 -C "admin@insuretech"
```

- Press **Enter** to accept the default file location (usually `C:\Users\<YourUser>\.ssh\id_ed25519`).
- Enter a passphrase (optional, but recommended for extra security), or press Enter twice to leave it empty.

## Step 2: Add the SSH Key to Your Droplet

### Option A: During Droplet Creation (Recommended if you haven't created it yet)
1. In the DigitalOcean dashboard, when creating a Droplet, go to the **Authentication** section.
2. Select **SSH Key** and click **New SSH Key**.
3. View the contents of your new public key by running this in PowerShell:
   ```powershell
   Get-Content ~/.ssh/id_ed25519.pub
   ```
4. Copy the output and paste it into the DigitalOcean dashboard.
5. Create a password for the Droplet as well if the UI allows, or set it later.

### Option B: Droplet is already created (Via DigitalOcean Console)
Since PowerShell `ssh` commands are not working, you can manually add your key using the DigitalOcean Console:

1. **Get your Public Key**: Run this command in your local PowerShell to view your public key:
   ```powershell
   type $env:USERPROFILE\.ssh\id_ed25519.pub
   ```
   *Copy all the output text (it should start with `ssh-ed25519...`).*

2. **Open the DigitalOcean Console**:
   - Go to your DigitalOcean dashboard.
   - Click on your Droplet.
   - Click the **Access** tab on the left menu.
   - Click **Launch Droplet Console** (or use the Recovery Console).
   - Log in with your `root` user and password if prompted.

3. **Add the Key to the Server**:
   Inside the DigitalOcean console terminal, run the following commands one by one to create the needed files and paste your key:

   ```bash
   # Create the .ssh directory if it doesn't exist
   mkdir -p ~/.ssh

   # Set correct permissions
   chmod 700 ~/.ssh

   # Open the authorized_keys file in the nano editor
   nano ~/.ssh/authorized_keys
   ```

4. **Paste Your Key**:
   - Paste the public key you copied from Step 1 into the `nano` editor.
   - Press `Ctrl+O` then `Enter` to save.
   - Press `Ctrl+X` to exit nano.

5. **Secure the File**:
   Run this final command in the console to secure the file:
   ```bash
   chmod 600 ~/.ssh/authorized_keys
   ```

You are now set up to securely manage the InsureTech backend using your SSH key, while retaining the fallback option to log in using the root password.
