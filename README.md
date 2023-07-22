# Exporter

If you have different export credentials like me and don't want to copy and paste it everytime 
and don't want to waste your time by searching them you need them, you can configure them to exporter your terminal.
		
Exporter is a CLI library for saving different kind of export values with different names

## Usage
Adding a new environment:

`exporter set --name my-a-company-infra-creds --envs "key1=value1,key2=value2,key3=value3"`

Output an existing environment's export values:

```exporter environmentname```

Listing environment configurations

`exporter`

### Command

```sh
$ make install 
$ exporter --help
```





### Flags

```sh
  --help                help for exporter
  delete                deletes environment configuration
  set                   add new environment configuration
  update                updates existing environment configuration
```

### Developer (From Source) Install

If you would like to handle the build yourself, instead of fetching a binary,
this is how we recommend doing it.

- Make sure you have [Go](http://golang.org) installed.

- Clone this project

- In the project directory run
```sh
$ make install
```

## License
helm-ssm is available under the MIT license. See the LICENSE file for more info.


### Installed packages
- go get -u github.com/spf13/cobra