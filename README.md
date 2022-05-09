# Matterfoss

Matterfoss is a community chat server based on code copied from the Mattermost™ project. It 
adds the following important features:

* **No trademarks**: Mattermost™ is a trademark of Mattermost Inc. In respect of their
[trademark policy](https://mattermost.org/trademark-standards-of-use/), we find it
inappropriate to use that name in a community free software project which can be changed and
remixed by anybody. As far as we know, the Matterfoss name is not trademarked.

* **No CLAs**: When you make a contribution to Matterfoss, you don't need to sign your code
over to anybody. There are no commercial licenses, everybody gets to copy this software under
the terms of the AGPL.

* **No DRM**: Every feature that is present in the open source software is fully unlocked, you
don't need any "license" file in order to use them.

* **No Tracking**: There is no telemetry, update-check or other "call home" feature which can
be used to track you. Note that if you use a 3rd party push notification provider, they will
have access to whatever data is sent to them.

## No warranties, no support

If you are considering using Matterfoss in a commercial setting, please note that this
software comes with no warranty. As far as we know, there is no available commercial
support.

This is a volunteer community project and we cannot even promise to fix important security
issues in a timely manner, so if you intend to use Matterfoss in a commercial setting, we
highly recommend that you instead choose Mattermost™. That way, you will be supporting the 
development of great open source software.

## Why?

Free communities need an alternative to closed proprietary systems like Slack. Unfortunately,
Mattermost™ has two features which made it unsuitable for such use:

1. Unless you buy an "enterprise license", even simple features like channel ownership are 
disabled, meaning anyone can delete any channel. Similarly, without an enterprise license,
you can't set a data retention policy, which can make GDPR compliance problematic.

2. Enterprise licenses are sold per head, so anyone hoping to set up a public server faces
potentially unlimited and significant license costs.

Attempts were made to engage with Mattermost Inc concerning these issues, but they were not
interested.

## Installing

We recommend that you compile Matterfoss on the server where you intend to use it, so that you
can quickly and easily make any changes that you need to. To compile it you will need
[golang](https://golang.org/dl/) and a database, either MariaDB or PostgreSQL.

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
the admin of your Matterfoss instance did not configure the links for their instance
properly. In order to avoid getting linked here, you need to edit your `config.json` file
or go into the System Console and edit the following entries:

* TermsOfServiceLink: url of your terms of service for your server
* PrivacyPolicyLink: url of your privacy policy
* AboutLink: url about the purpose of your server
* HelpLink: url with a link to help and support
* ReportAProblemLink: url where people can go to report problems with your service

## License

See the [LICENSE file](LICENSE.txt) for license rights and limitations.

