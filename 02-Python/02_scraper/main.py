#!/usr/bin/env python
# -*- coding: utf-8 -*-

import urllib3
from bs4 import BeautifulSoup

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
## Good: nice example of the generator pattern
## Bad: HTTP requests are done in series
def getUrlGeneratorStyle(url, maxDepth, http):
    "retrieves an url"
    # A list of tasks

    queue = [GetUrlTask(url, 0)]

    while len(queue) > 0:
        task = queue.pop(0)
        response = http.request("GET", task.url)

        if response.status == response.status:
            yield GetUrlResponse(task, response.data);
            if task.depth < maxDepth:
                queue.extend(findLinks(response.data, task.url))

def findLinks(str, url):
    ## TODO
    return []

def findImages(str, url):
    ## TODO
    return []

for response in getUrlGeneratorStyle("http://spatineo.com", 1, http):
    images = findImages(response.data, response.task.url)
    print images
