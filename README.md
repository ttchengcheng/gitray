# Gitray

add an icon to systray, to indicate uncommited / unpushed changes of your git repositories

![img](img/tray.png)
![img](img/menu.png)

## Usage

### download

```sh
git clone https://github.com/ttchengcheng/gitray.git
```

### build

```sh
cd gitray
go build
```

### add git repositories

```sh
# a project at /Users/yourname/project/project1
cd /Users/yourname/project/project1
# and the cloned gitray is at /Users/yourname/project/gitray
pwd >> /Users/yourname/project/gitray/config.txt

# There is another project at /Users/yourname/project/project2
cd /Users/yourname/project/project1
pwd >> /Users/yourname/project/gitray/config.txt
```

### Run

```sh
# macOS
./gitray &

```

PS: it is not tested on win, maybe it works 😛