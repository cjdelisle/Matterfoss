# Matterfoss
Matterfoss is a community chat server based on code copied from the Mattermost(tm) project.

* **No trademarks**: Mattermost(tm) is a trademark of Mattermost Inc and in respect of their
[trademark policy](https://mattermost.org/trademark-standards-of-use/), we find it
inappropriate to use that name in a community free software project which can be changed and
remixed by anybody. The Matterfoss name is not (as far as we know) trademarked.
* **No CLAs**: When you make a contribution to Matterfoss, you don't need to sign your code
over to anybody. There are no commercial licenses, everybody gets to copy this software under
the terms of the AGPL.
* **No DRM**: Every feature that is present in the open source software is fully unlocked, you
don't need any "license" file in order to use them.
* **No Tracking**: There is no telemetry, update-check or other "call home" feature which can
be used to track you. Note that if you use a 3rd party push notification provider, they will
have access to whatever data is sent to them.

## No warrantees, no support
If you are considering using Matterfoss in a commercial setting, please note that this
software comes with no warrantee and (as far as we know) there is no available commercial
support.

This is a volunteer community project and we can not promise even to fix important security
issues in a timely manner, so if you intend to use Matterfoss in a commercial setting, we
highly recommend that you instead choose Mattermost(tm), you will be supporting the development
of great open source software.


## Installing
We recommend that you compile Matterfoss on the server where you intend to use it, so that you
can quickly and easily make any changes that you need to. To compile it you will need
[golang](https://golang.org/dl/) and a database, we recommend MariaDB.

```bash
# Download this server and the frontend
git clone https://github.com/cjdelisle/Matterfoss
git clone https://github.com/cjdelisle/MatterfossWeb

# Compile the frontend
cd MatterfossWeb
npm install
npm run build

# Compile matterfoss
cd ../Matterfoss
go build -o ./bin ./...

# Setup the database - note that mfuser and mfpass are defaults in the
# configuration. If your MariaDB instance is on another server, you will
# need to change 'mfuser'@'localhost' to 'mfuser'@'%' to allow outside
# access and clearly you will need to use a strong password.
echo "
create user 'mfuser'@'localhost' identified by 'mfpass';
create database matterfoss;
grant all privileges on matterfoss.* to 'mfuser'@'localhost';
" | mariadb

# Move the frontend into the server so that it will be served
mv ../MatterfossWeb/dist ./client

# Start matterfoss
./bin/matterfoss
```

## Customization
If you just clicked on a link and found yourself here, it's probably because you or
the admin if your Matterfoss instance did not configure the links for their instance
properly. In order to avoid getting linked here, you need to edit your config.json file
or go into the System Console and edit the following entries:

* TermsOfServiceLink: url of your terms of service for your server
* PrivacyPolicyLink: url of your privacy policy
* AboutLink: url about the purpose of your server
* HelpLink: url with a link to help and support
* ReportAProblemLink: url where people can go to report problems with your service
