# unbound

## Name

*unbound* - perform recursive queries using libunbound.

## Description

Via *unbound* you can perform recursive queries. Unbound uses DNSSEC by default when resolving *and*
it returns those records (DNSKEY, RRSIG, NSEC and NSEC3) back to the clients. The *unbound* plugin
will remove those records when a client didn't ask for it. The internal (RR) answer cache of Unbound
is disabled, so you may want to use the *cache* plugin.

Libunbound can be configured via (a subset of) options, currently the following are set:

* `msg-cache-size`, set to 0
* `rrset-cache-size`, set to 0

The *unbound* plugin uses <https://github.com/miekg/unbound> to interface with libunbound.

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
    option NAME VALUE
}
~~~

* **FROM** as above.
* **IGNORED_NAMES** in `except` is a space-separated list of domains to exclude from resolving.
* `option` allows setting unbound options (see unbound.conf(5)), this can be specified multiple
  times.

## Metrics

If monitoring is enabled (via the *prometheus* directive) then the following metric is exported:

* `coredns_unbound_request_duration_seconds{}` - duration per query.
* `coredns_unbound_response_rcode_count_total{rcode}` - count of RCODEs.

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

or

~~~ corefile
example.org {
    unbound
}
~~~

Resolve everything except queries for example.org (or below):

~~~ corefile
. {
    unbound {
        except example.org
    }
}
~~~

Enable [DNS Query Name Minimisation](https://tools.ietf.org/html/rfc7816) by setting the option:

~~~ corefile
. {
    unbound {
        option qname-minimisation yes
    }
}
~~~

## Bugs

The *unbound* plugin depends on libunbound(3) which is C library, to compile this you have
a dependency on C and cgo. You can't compile CoreDNS completely static. For compilation you
also need the libunbound source code installed (`libunbound-dev` on Debian).

## See Also

See <https://unbound.net> for information on Unbound and unbound.conf(5). See
<https://github.com/miekg/unbound> for the (cgo) Go wrapper for libunbound.
