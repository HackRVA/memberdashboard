# Dev Setup

## The Easy Way (using Remote-Containers)

Install the following

- [Docker](https://www.docker.com/products/docker-desktop)
- [VS Code](https://code.visualstudio.com/download)
- VS Code [Remote-Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension

That's it! Nothing else needs to be installed.

If you are on Windows, then it is recommended to clone to the repo to a folder on the WSL for best performance.

Either open the folder in VS Code and choose reopen in container or run the Remote-Containers: Open Folder in Container... command and select the local folder.

![Open from container](/img/openFromContainer.gif "Open from container")

The backend server can be started by debugging in VS Code. Otherwise, you can start the server without debugging using `sh buildandrun.sh`.

Start the web app as described in the [UI README](/web/README.md).

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
