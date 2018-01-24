# unbound

## Name

*unbound* - perform recursive queries using libunbound

## Description

Via *unbound* you can perform recursive queries to the internet.

## Syntax

~~~
unbound [FROM]
~~~

* **FROM** is the base domain to match for the request to be resolved. If not specified the zones
  from the server block are used.

More features utilized with an expanded syntax:

~~~
unbound [FROM] {
    except IGNORED_NAMES...
}
~~~

* **FROM** as above.
* **IGNORED_NAMES** in `except` is a space-separated list of domains to exclude from resolving.

## Examples

Resolve all requests within example.org.

~~~
unbound example.org
~~~

Resolve everything except requests to miek.nl or example.org

~~~ corefile
. {
    unbound .{
        except miek.nl example.org
    }
}
~~~

## Bugs

The *unbound* plugin depends on libunbound(3) which is C library, to compile this you have
a dependency on C and cgo. You can't compile CoreDNS completely static. For compilation you
also need the libunbound source code installed.

## See Also

See <https://unbound.net> for information on Unbound.
