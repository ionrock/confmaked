# Confmaked

<blockquote class="twitter-tweet" data-lang="en"><p lang="en"
dir="ltr">Has no one built a config management system around Makefiles
yet? It always seemed like a natural fit.</p>&mdash; ionrock
(@ionrock) <a
href="https://twitter.com/ionrock/status/969082259499372544?ref_src=twsrc%5Etfw">March
1, 2018</a></blockquote>

Often times configuration management systems attempt to avoid
duplicating work. For example, [chef](https://chef.io) has a somewhat
complex two-phase compilation mode that will create a dependency graph
before actually running the cookbooks. This gets really copmlicated
and confusing, especially if you're not interested in learning the DLS
or YAML of your configuration management system of choice.

Developers have been dealing with avoiding extra work for a very long
time in Makefiles. So `confmaked` is a configuration management system
that allows specifiying Makefiles that will do the work!

I've mentioned this seemed like a good idea, so I figured I should put
my money where my mouth is.

For those who might be unfamilier with
[make](https://www.gnu.org/software/make/), it is a build tool that
works by creating targets. A target is often a file and targets won't
rebuild a target unless something that it is dependent on
changes. Here is a really simple example.


```
SOURCE_CODE=$(shell find . -name '*.go')

confmaked: $(SOURCE_CODE)
	go build ./cmd/confmaked
```

Make will look at all the files and only rebuild the binary if one of
those file changed. Hopefully, you can see where I'm going with this.

If you wanted to install a binary you might write something like:

```
/usr/local/bin/nginx:
	apt install nginx
```

If you wanted to write a config file you might do something like this:

```
config_data:
	wget -r https://config.internal.myorg.io/path/to/app/

/etc/myapp.conf:
	we -d config_data --template config_data/etc/myapp.conf.tmpl
```

Shameless plug. I wrote a tool called
[we](https://github.com/ionrock/we) that can read YAML/JSON and inject
it as environment variables or write a config before running a
command. So that is what is happening. You could also just download a
config file and build it ahead of time.

```
/etc/myapp.conf:
	curl -o /etc/myapp.conf https://config.internal.myorg.io/myapp.conf
```

You can start up services pretty easily since tools like systemd
should be OK with the noop if its already started.

```
start_myapp:
	service myapp start
```

There are a lot of good resources for how to be creative with
Makefiles so I'll leave it up to you to learn more there.

## How do I use it?

The idea is you run `confmaked` with a set of configured
Makefiles. The configuration will eventually provide support for
different protocols for grabbing the Makefiles such as HTTP and likely
gRPC, but for now it is just a list of local files.

Each file will be run in order and if anything fails it will
exit. Eventually it will allow running things on a cadence, but for
now, you can use your init system, cron or a process supervisor to try
things again.

Here is a simple example:

```
- name: myapp
  url: /etc/confmaked/myapp/Makefile
```

It will just call `make` in the directory at the moment. I'm sure I'll
eventually allow adding specific targets pretty soon here.

That is pretty much it.

You can also copy the Makefiles in to a container and get the same
behavior.

```
FROM ubuntu

RUN apt-get update && apt-get install -yq make
ADD confmaked /usr/local/bin/confmaked
ADD myapp.conf /etc/confmaked.yml
RUN confmaked --config /etc/confmaked.yml
```

Feel free to swap `docker` with `packer` or `vagrant` or `terraform`
and you can hopefully see how this might be a sort of OK idea.

## Seriously?

Sort of, yes! Often times developers already write code that is
effectively what a config management system does, so this might allow
some reuse there. Also, it makes it easy to reuse the same Makefiles
in things like containers, dev environments or whereever with minimal
setup. There is no need to get an ansible environment configured, set
up chef server or deal with chef zero.

From a scaling perspective, it also seems like it could be really nice
because the Makefiles can be hosted anywhere. The Makefiles can
download extra scripts, binaries, configuration data, etc. and the
file management aspects of make should just work.

What about security you ask? The only requirement is that make is
installed, not an entire compilier toolchain. The point is not to
allow building your code on your servers, but rather, define how to
configure your servers the same way you build code.

Some people really don't like make and that is OK too. I always chose
make because it is a lowest common denominator. Also, it should be
able to handle installing packages, writing config files and running
scripts without much trouble or extra dependencies. I think there are
other great build tools out there that could also work, but I'm trying
make first.
