switchconf
==========

switchconf is a simple utility to comment/uncomment regions of files.
Specifically it was created to change configuration files according to my
location.

Example
=======

For example, say you have an internal yum repository in your job and you
want to use that when you're in the office, but not when at home. You
have the .repo file like this at home:

    [fedora]
    name=Fedora $releasever - $basearch
    mirrorlist=https://mirrors.fedoraproject.org/metalink...
    enabled=1
    gpgcheck=1
    ...

And when you are at work you change it to

    [fedora]
    name=Fedora $releasever - $basearch
    baseurl=http://linux-ftp.mycompany.com/pub/mirrors/fedora/...
    enabled=1
    gpgcheck=1
    ...

Instead of changing it everytime, you can instead create this:

    [fedora]
    name=Fedora $releasever - $basearch
    #:switchconf: work
    baseurl=http://linux-ftp.mycompany.com/pub/mirrors/fedora/...
    #:switchconf: home
    mirrorlist=https://mirrors.fedoraproject.org/metalink...
    #:switchconf:
    enabled=1
    gpgcheck=1
    ...

Then you run add these lines to /etc/switchconf (or ~/.switchconf):

    /etc/yum.repos.d/fedora.repo:
       comment: #

Running `switchconf home` will result in the file changing to

    [fedora]
    name=Fedora $releasever - $basearch
    #:switchconf: work
    #:off:baseurl=http://linux-ftp.mycompany.com/pub/mirrors/fedora/...
    #:switchconf: home
    mirrorlist=https://mirrors.fedoraproject.org/metalink...
    #:switchconf:
    enabled=1
    gpgcheck=1
    ...

Similarly, running `switchconf work` will result in

    [fedora]
    name=Fedora $releasever - $basearch
    #:switchconf: work
    baseurl=http://linux-ftp.mycompany.com/pub/mirrors/fedora/...
    #:switchconf: home
    #:off:mirrorlist=https://mirrors.fedoraproject.org/metalink...
    #:switchconf:
    enabled=1
    gpgcheck=1
    ...

Installation
============

    go get github.com/robteix/switchconf

License
=======

Switchconf is licensed under the GPLv3.
