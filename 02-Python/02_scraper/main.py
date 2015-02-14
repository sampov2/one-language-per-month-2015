#!/usr/bin/env python
# -*- coding: utf-8 -*-

import urllib3

http = urllib3.PoolManager()

# TODO: try including this _within_ the function. Maybe it can be scoped there?
class GetUrlTask:
    """A task for the getUrl function"""
    def __init__(self, url, depth):
        self.url = url
        self.depth = depth

class GetUrlResponse:
    """Response for a GetUrlTask"""
    def __init__(self, task, data):
        self.task = task
        self.data = data

## Recursive scraping using the generator pattern
def getUrl(url, maxDepth, http):
    "retrieves an url"
    # A list of tasks

    queue = [GetUrlTask(url, 0)]

    while len(queue) > 0:
        task = queue.pop(0)
        response = http.request("GET", task.url)

        if response.status == response.status:
            yield GetUrlResponse(task, response.data);
            if task.depth < maxDepth:
                queue.extend(findLinks(response.data, task.url, task.depth+1))

def findLinks(str, url, depth):
    ## TODO: implement
    return []

for response in getUrl("http://spatineo.com", 1, http):
    ## TODO: scrape for payload
    print response.data;
