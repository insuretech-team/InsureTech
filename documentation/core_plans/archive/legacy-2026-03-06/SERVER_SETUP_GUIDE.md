# Initial Server Setup Guide (Ubuntu/Debian)

This guide covers the fundamental steps you should take immediately after logging into a new server as `root` for the first time. It will increase the security and usability of your server.

## Step 1: Create a New Regular User

Logging in as `root` is discouraged because it has absolute power over the system. It's best practice to create a regular user for daily tasks.

Run the following command to create a new user (replace `insureadmin` with your preferred username):

```bash
adduser insureadmin
```

You will be asked to create and confirm a new password for this user. Remember this password! You can press `Enter` to skip the other informational fields like Name, Room Number, etc.

## Step 2: Grant Administrative Privileges (sudo)

Your new user needs the ability to run administrative commands. We do this by adding the user to the `sudo` group.

```bash
usermod -aG sudo insureadmin
```

Now, when logged in as `insureadmin`, you can run superuser commands by typing `sudo` before them.

## Step 3: Copy SSH Keys to Your New User

If you logged into `root` using your SSH key, you should copy that key to your new user so you can log in as the new user securely.

As `root`, run this command to reliably copy the SSH configuration and keys to your new user (make sure you replace `insureadmin` with the username you created in Step 1):

```bash
rsync --archive --chown=insureadmin:insureadmin ~/.ssh /home/insureadmin
```

## Step 4: Set Up a Basic Firewall (UFW)

Ubuntu servers can use the UFW (Uncomplicated Firewall) to ensure only connections to approved services are allowed. Since we only want SSH access for now (and maybe HTTP/HTTPS later):

1. **Allow SSH connections** (CRITICAL, so you don't get locked out!):
   ```bash
   ufw allow OpenSSH
   ```

2. **Enable the firewall**:
   ```bash
   ufw enable
   ```
   *Type `y` and press `Enter` if it warns you about disrupting existing SSH connections.*

3. **Check the firewall status**:
   ```bash
   ufw status
   ```
   You should see `OpenSSH` allowed.

## Step 5: Test Your New User Login

**Before logging out of your `root` session**, open a **new** PowerShell window on your local machine and try connecting as your new user:

```powershell
ssh insureadmin@<YOUR_SERVER_IP_ADDRESS>
```

If you are logged in successfully (and did not get prompted for a password if you set up the SSH key), you're all set! 

From now on, you should log into the server using `ssh insureadmin@<IP>` instead of `root`. When you need to run administrative commands, just type `sudo` before the command (e.g., `sudo apt update`).
