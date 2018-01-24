# unbound

## Name

*unbound* - perform recursive queries using libunbound

## Description

Via *unbound* you can perform recursive queries to the internet. Unbound uses DNSSEC by default when
resolving *and* it returns those records (DNSKEY, RRSIG, NSEC and NSEC3) back to the clients. This
also implies that OPT RR in the original request is not reflected upstream.

The internal answer cache of Unbound is disabled, so you may want to use the *cache* plugin.

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

## Metrics

If monitoring is enabled (via the *prometheus* directive) then the following metric is exported:

* TODO(miek)

## Examples

Resolve queries for all domains:
~~~ corefile
. {
    unbound
}
~~~

Resolve all queries within example.org.

~~~ corefile
. {
    unbound example.org
}
~~~

Resolve everything except queries for miek.nl or example.org

~~~ corefile
. {
    unbound . {
        except miek.nl example.org
    }
}
~~~

## Bugs

The *unbound* plugin depends on libunbound(3) which is C library, to compile this you have
a dependency on C and cgo. You can't compile CoreDNS completely static. For compilation you
also need the libunbound source code installed.

## See Also

See <https://unbound.net> for information on Unbound and unbound.conf(5).
