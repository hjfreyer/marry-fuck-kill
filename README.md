Intro
-----

This is the code for Marry Fuck Kill, a web game.

Up-to-date version always available at the
[Official Location](https://github.com/hjfreyer/marry-fuck-kill).

Hacking
-------

Prerequisites: node.js, ruby, sass. See "Installation Notes" below for brief
instructions for installing these.

In order to build the dev server, you need to have the appengine
python runtime installed, as well as npm. As a one-time set up (or if
you pull and something seems weird), run `npm install`.

After that, run `npm start` to bring up the dev server on
localhost:8080. Sass files and other preprocessed junk will get
automatically updated by grunt under the covers.

If nothing is listening on port 8080, make sure running `dev_appserver build/`
works. (It may, for instance, print a prompt to STDOUT when it's newly installed.)

Installation Notes
------------------
Here's how to get each of the prerequisites.

*ruby*: Use your system's ruby package. (`apt-get install ruby` on
Debian/Ubuntu.)

*node.js*: Install from http://nodejs.org/download/. Don't rely on OS packages,
they can get quite outdated.

*sass*: Since you already have ruby: `gem install sass`

This was tested on Ubuntu 14.04 LTS, but hopefully it's generic.

Authors
-------
* [Michael Kelly](http://michaelkelly.org)
* [Hunter Freyer](http://www.hjfreyer.com)
