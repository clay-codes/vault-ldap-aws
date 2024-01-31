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
- once finished, will output ssh command to openssh to windows server
- will also output ldapsearch command to be ran from any machine with openldap--can easily be installed via:

### Server Info
- username is `Administrator`
- password is `admin`
- test AD domain is `vaultest.com`
- password complexities have been disabled
- once ssh connection established, all commands are run in batch, just type `powershell` to switch
- `exit` to close connection