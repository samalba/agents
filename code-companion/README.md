# Code Companion

Code Companion is an agent designed to facilitate seamless development workflows by enabling efficient file synchronization between your local environment and a remote host. The agent relies on [Mutagen](https://mutagen.io/) for fast, bidirectional file replication, ensuring that changes made locally are reflected remotely and vice versa.

## File Replication with Mutagen

**Code Companion requires Mutagen to replicate files between your local machine and the host.** Before using the agent, you must set up file replication using Mutagen. This ensures that your code and changes are always up to date in both environments.

### Prerequisites
- [Install Mutagen](https://mutagen.io/documentation/introduction/installation) on your local machine.
- Ensure you have an SSH key pair available (e.g., `~/.ssh/id_ed25519.pub`).

### Starting File Replication

Run the following commands in your terminal to start file replication:

```sh
# 1. Create a local directory that will be used for replication
mkdir -p ./my-code

# 2. Start the Mutagen agent (this will use SSH to replicate a Dagger shared Cache volume)
dagger -m github.com/samalba/agents/code-companion call mutagen-agent up --authorized_keys ~/.ssh/id_ed25519.pub --ports 1222:22
```

In another terminal, run:

```sh
# 3. Create a sync session between your local code and the remote host
mutagen sync create --name=MyCode ./my-code root@localhost:1222:~/dagger

# 4. Monitor the sync session
mutagen sync monitor MyCode
```

Once replication is active, you can use Code Companion to interact with your code in the remote environment seamlessly.

## Use the Code Companion

```sh
dagger -m github.com/samalba/agents/code-companion
⋈ ask "Write me a hello world program in C."
```

You'll see the files appear in the `my-code` directory as the agent completes the assignment.
