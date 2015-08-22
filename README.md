re:Invent sessions API
---

## Description

Search sessions on AWS re:Invent  
https://api.supinf.co.jp/v1/reinvent-sessions?q=Expert%20400

## Basic Usage

### 1. Install go binary

```shell
$ go get github.com/supinf/reinvent-sessions-api/
```

### 2. Run this application

```shell
$ AWS_REGION=us-west-2 AWS_ACCESS_KEY_ID=? AWS_SECRET_ACCESS_KEY=? APP_PORT=9000 reinvent-sessions-api
```

### 3. Access the application

[http://localhost:9000/](http://localhost:9000/)

## Usage with DynamoDB Local, VirtualBox, CoreOS & Docker containers

### 1. Install VirtualBox & Vagrant

- [VirtualBox](https://www.virtualbox.org/)
- [Vagrant](http://www.vagrantup.com/)

### 2. Install vagrant-hostsupdater plugin

```shell
$ vagrant plugin install vagrant-hostsupdater
```

### 3. Change your working directory to vagrant folder

```shell
$ cd /path/to/this-repository-root/vagrant
```

### 4. Create a virtual machine with AWS credentials

```shell
$ AWS_REGION=us-west-2 AWS_ACCESS_KEY_ID=? AWS_SECRET_ACCESS_KEY=? vagrant up
```

### 5. Confirm whether a service is running

[http://reinvent-sessions-api.local/](http://reinvent-sessions-api.local/)

### 6. Test the application

```shell
$ vagrant ssh -c "docker run -it --rm -v /home/core/share:/go/src/github.com/supinf/reinvent-sessions-api supinf/reinvent-sessions:base go test github.com/supinf/reinvent-sessions-api/..."
```

### 7. Teardown the VM

```shell
$ vagrant halt
```

## Contribution

1. Fork ([https://github.com/supinf/reinvent-sessions-api/fork](https://github.com/supinf/reinvent-sessions-api/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request

## Copyright and license

Code and documentation copyright 2015 SUPINF Inc. Code released under the [MIT license](https://github.com/supinf/reinvent-sessions-api/blob/master/LICENSE).
