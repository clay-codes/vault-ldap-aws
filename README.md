# COOL-AD
### Remote AWS AD Environment

#### Prerequisites
- UNIX/LINUX
- go installed if not running on Apple M1
- Doormat CLI to authenticate to AWS
- (optional) openldap to test--can install via:
```
Ubuntu
sudo apt-get update
sudo apt-get install ldap-utils

On CentOS/RHEL:
sudo yum install openldap-clients

On macOS:
brew install openldap

```
  

#### Quick Setup Via Binary
- binary compiled on Apple M1 Max 32gb
- owners of machine with same architecture just downoad `cool-ad` and run `./coolad`

#### Setup
- clone repo; cd into 
- run `go run .`

### Usage
- first time use, enter `no/n` at cleanup prompt
- wait about 5 min for resources to be created
- once finished, will output ssh command to openssh to windows server/AD domain controller
- will also output ldapsearch test command; openldap not installed on server
- outputs powershell commands to view AD details, which can be executed on server
- if error recieved when attempting powershell commands, need to wait another few minutes for AD bootstrap completion

### Server Info
- domain controller
- username is `Administrator`
- password is `admin`
- AD domain is `vaultest.com`
- password complexities have been disabled
- once ssh connection established, all commands are run in batch, just type `powershell` to switch
- `exit` to close connection