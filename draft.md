htmx knowledge share

july 3rd week



talk about htmx from high level perspective


## what is html?
htmx is html6
htmx is high power tools for html
htmx is what html could be if we updated html instead of modern web development with all the javascript frameworks

htmx is a library that allows you to access modern browser features directly from HTML, rather than using javascript

## break down of htmx 

modern web dev apis send json responses and the framework puts builds that object into html before rendering that html
htmx lets you perform a request and swap in the response directly to the html element of your choosing

modern web development model utilizes json as the response body which is consumed by a js framework of choice and used for the final representation of html
while htmx does the request, accepts html and swaps it into the DOM

### declaring the AJAX request
htmx allows you to issue AJAX requests directly from HTML
get post put patch delete

### triggering requests
by default, natural events of an html tag is considered in order to trigger the AJAX request.
input, textarea & select are triggered on the change event
form is triggered on the submit event
everything else is triggered by the click event

mouseenter

#### modifying triggers
hx-trigger="mouseenter once"

#### trigger filters

#### special events

#### polling

#### load polling

#### request indicators

### targets

#### extended css selectors

### swapping

innerHTML by default

#### view transitions

#### swap options

### synchronization

### css transitions

#### details

### out of band swaps

#### troublesome tables
