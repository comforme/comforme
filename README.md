# comfor.me
Everyone's personal community

#What is it?
Comfor.me (Community for Me) is a community-rated and identity-oriented 
social network/service listing. Users can find accepting communities and
services based on a wide array of keywords. Users can also start their own
communities categorized by aforementioned keywords. Comfor.me makes it easier
for an individual to find communities and services which accept them for who
they are.

#Technologies Used:
* The Go programming Language
* PostGreSQL
* Vagrant

## Instructions for setting up development server
### Deploying to Heroku
TODO

### Using Vagrant for Development
* First install VirtualBox and Vagrant
* Vagrant will be able to create a virtual machine hosting the
  comfor.me application locally. The application is installed
  via a shell script when the machine is provisioned, and starts
  the application on port 8080, which is forwarded to the local
  machine. Just run `vagrant up`.
* While developing, you can redeploy your changes to the code to
  the vagrant environment by changing to the /vagrant directory
  and executing the redeploy script `./vagrant_redeploy.sh` This
  will update the local copy, clear Go's cache, and reinstall
  the application.
* When you are finished developing you can dispose of the virtual
  machine by executing `vagrant destroy`.
