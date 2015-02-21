[![Build Status](https://travis-ci.org/comforme/comforme.svg)](https://travis-ci.org/comforme/comforme)
# comfor.me
Everyone's personal community

"All young people, regardless of sexual orientation or identity, deserve a safe and supportive environment in which to achieve their full potential"
-Harvey Milk

#What is it?
Comfor.me (Community for Me) is a social network / service listing based on communities of shared identities among users. Users can easily see contributions from other members in their communities, creating a safe and supportive space for shared needs and interests. Users can also create new communities to be found by each other and build new networks of support and helpful information. Comfor.me makes it easier for anyone to find the communities and services which embrace them for who they are.

##The User Experience
When you first visit comfor.me you are invited to create an account, followed by identifying some of the communities you identify with. Once you've joined, you can search and post pages on any topic you can imagine. Each page describes something, whether it be a business, local event, online article, piece of artwork, or anything else. Page descriptions are followed by individual contributions from other users, such as reviews for businesses, comments on articles, or any type of additional info they have that creates rich feedback from diverse communities of users like you. Whenever you have something you'd like to share, you can create your own contribution to a page or make a new page for anything you'd like to get thoughts on from other users.

###What Makes the Experience Better?
####Seamless Registration
Enter your email and pick your username, done! You can immediately discover your communities since email verification is automatically handled the next time you sign in.
####Modern, Responsive Site Design
With a website that offers a smooth, clean interface that gives you what you want with no unnecessary frills, and servers written in an optimized language that react quickly and make sure load times are a thing of the past. Comfor.me provides a great user experience for any community member.
####Freedom
Other social networks makes you use your real name, lock you into only one kind of content, limit finding communities you identify with, or do any number of other things to ensure you're easy to predict and generate as much advertising revenue as possible. Comfor.me does none of these things, instead we help you express your freedom to create and be a part of anything you want to, we build the tools for you to find people like you, then we step back when you step forward and become part of our community.

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
