# Spider

## Design

Firstly, we can see the raw code of `http://shakespeare.mit.edu/`. We can find that all hyperlinks under
Comedy, History, and Tragedy href to a index.html where we can see a hyper link to full version of this
set. For example, when I click `All's Well That Ends Well` in the homepage, I will go to a index page of
`All's Well That Ends Well` with a hyperlink `Entire play in one page` href to `http://shakespeare.mit.edu/allswell/full.html`.

So I can peek all hyperlinks under Comedy, History and Tragedy using regexp `<br><a href="(.*)/index.html">`.
Concat the submatch string with `/full.html` and the website path, we will get the target URL.

Secondly, the Poetry is a little different. Click `The Sonnets` will go to a index page that list many poetry
 pages whose link can be peeked by using regexp `<DT><A HREF="(.*)">` in the list page.

Finally, poetry except `The Sonnets` all can directly reach by using regexp `<em><a href="(.*)">` to peek their
 URL.

## Implement

By assessment, the target data is very small(<100MB) so I can easily put them all in memory. To avoid network
crash problems, I take proper reconnect operation in http client. To make best of multi-processor, I use multi-
threads to process network request and regexp running. It is worth a visit that I use **Atomic Integer** and
**Compare And Set** atomic operation to avoid Mutex to increase parallel performance. Additional buffered IO of
 disk file is necessary when frequent.

## Code Style

Standard Go code style but I like to use `goto`.