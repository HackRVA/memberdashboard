# Dev Setup
> Note: The Docker and Remote-Containers setup will soon be deprecated for this project. Please refer to [Getting Started](./getting_started.md) for more updated dev setup instructions.

## The Easy Way (using Remote-Containers)

To get started quickly, follow these steps:

1. Install [Docker](https://www.docker.com/products/docker-desktop), [VS Code](https://code.visualstudio.com/download), and the VS Code [Remote-Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension.
2. Clone the repo to a folder on WSL if you are using Windows for better performance.
3. Open the folder in VS Code and choose "Reopen in Container" or run the "Remote-Containers: Open Folder in Container..." command and select the local folder.

![Open from container](/img/openFromContainer.gif "Open from container")

You don't need to install any other dependencies.

To start the backend server, debug it in VS Code. Alternatively, you can use `sh buildandrun.sh` to start the server without debugging.

Please refer to the [UI README](/web/README.md) for instructions on starting the web app.

```
# navigate to ui folder
cd ui

# install node modules
$ npm ci

# run local env
$ npm run start
```

If you feel cramped in the VS Code terminal pane, you can still connect to dev container shell from your favorite terminal using

```
docker exec -it -u vscode memberdashboard_dev_1 bash
```

Now, go write code and implement features!

---
