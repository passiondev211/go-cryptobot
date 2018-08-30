# Deployment with Ansible

This document describes how to deploy this project and its environment with Ansible.

## Prepare machines (not applicable in case of aws ec2)

Ansible deployment implies you should have all access keys manually set and your user
is able to run passwordless `sudo`.

Also, SeLinux should be managed manually and can be turned off for debugging purposes
by running `sudo setenforce Permissive`.

Ensure the host machine can access the target machine by key (see `group_vars/all` for the location
of this key) and the user on the target machine is on the sudoers list with NOPASSWD option.

SSH as `<username>`, `su` and:

```
#!bash

# Install vim just for convenience.
sudo yum -y install vim

# Add sudo privileges to your user.
usermod -aG wheel <username>
```

Then exit `root` and add your public key:

```
#!bash

# Ensure the authorized_keys file exists.
mkdir ~/.ssh
chmod 700 ~/.ssh
touch ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys

# Drop your public key here.
vim ~/.ssh/authorized_keys
```

Now you must be able to login without password and add NOPASSWD option to the `wheel` group
in `/etc/sudoers`:

```
#!bash

# Edit /etc/sudoers enabling NOPASSWD option.
sudo chmod +w /etc/sudoers
sudo vim /etc/sudoers
sudo chmod -w /etc/sudoers
```

## Deploy

Before running deploy commands:
1. Ensure you can access the target machine specified in file `staging`
by running `ansible all -m ping -i staging`.
1. Ensure the target machine can access the target repository by creating
a keypair with `ssh-keygen` and adding the public key to the repo.
    
Then perform deployment:

    ansible-playbook -i staging site.yml
