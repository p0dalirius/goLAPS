![](./.github/banner.png)

<p align="center">
  A simple way to read and write LAPS passwords from linux.
  <br>
  <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/p0dalirius/pyLAPS">
  <a href="https://twitter.com/intent/follow?screen_name=podalirius_" title="Follow"><img src="https://img.shields.io/twitter/follow/podalirius_?label=Podalirius&style=social"></a>
  <a href="https://www.youtube.com/c/Podalirius_?sub_confirmation=1" title="Subscribe"><img alt="YouTube Channel Subscribers" src="https://img.shields.io/youtube/channel/subscribers/UCF_x5O7CSfr82AfNVTKOv_A?style=social"></a>
  <br>
</p>

This script is a go setter/getter for property `ms-Mcs-AdmPwd` used by LAPS inspired by [@swisskyrepo](https://github.com/swisskyrepo/)'s [SharpLAPS](https://github.com/swisskyrepo/SharpLAPS) in C#.

Require (either):
  * Account with `ExtendedRight` or `GenericRead` to get LAPS passwords
  * Account with `ExtendedRight` or `GenericWrite` to set LAPS passwords
  * Domain Admin privileges

## Usage

```
                __    ___    ____  _____        
   ____ _____  / /   /   |  / __ \/ ___/       
  / __ `/ __ \/ /   / /| | / /_/ /\__ \      
 / /_/ / /_/ / /___/ ___ |/ ____/___/ /         
 \__, /\____/_____/_/  |_/_/    /____/    v1.2
/____/           @podalirius_                   

[!] Option -host <host> is required.
Usage of ./bin/goLAPS:
  -debug
    	Debug mode
  -domain string
    	(FQDN) domain to authenticate to.
  -hashes string
    	NT/LM hashes, format is LMhash:NThash.
  -host string
    	IP Address of the domain controller or KDC (Key Distribution Center) for Kerberos. If omitted it will use the domain part (FQDN) specified in the identity parameter.
  -password string
    	password to authenticate with.
  -port int
    	Port number to connect to LDAP server.
  -quiet
    	Show no information at all.
  -use-ldaps
    	Use LDAPS instead of LDAP.
  -username string
    	User to authenticate as.
```

## Contributing

Pull requests are welcome. Feel free to open an issue if you want to add other features.
