# go-ssr

This is a small experiment to see what’s possible with Go SSR (server-side rendering). The idea is to play to Go
strengths; simplicity, predictability, etc. and to use it for server-side rendering versus server-side rendering React.

Ultimately, the idea is to combine this approach with client-side rendered React so that pages are indexable by default
(because of static meta tags, which are dynamic but generated on the server). This mitigates many concerns; users get
client-side experiences while Google and bots can scrape site previews, and developers don’t need to think about SSR
concerns.

Finally, this approach should be applicable to any framework. Nothing about this approach in general should be concerned
with the implementation details of the client-side bundle. This does not implement client-side JavaScript but could be
adapted to as a starting point.
