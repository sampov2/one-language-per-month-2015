#!/usr/bin/env python
# -*- coding: utf-8 -*-

import urllib3
from bs4 import BeautifulSoup
from urlparse import urljoin

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
        try:
            response = http.request("GET", task.url)
            if response.status == response.status:
                yield GetUrlResponse(task, response.data);
                if task.depth < maxDepth:
                    queue.extend(findLinks(response.data, task.url, task.depth+1))
        except urllib3.exceptions.RequestError:
            print "ERROR retrieving url "+task.url
            continue


def findLinks(html, baseUrl, depth):
    ret = []
    soup = BeautifulSoup(html)
    for a in soup.find_all('a'):
        href = a.get('href')
        if href != None:
            link = urljoin(baseUrl, href)
            ret.append(GetUrlTask(link, depth))

    return ret

def findImages(html, baseUrl):
    ret = []
    soup = BeautifulSoup(html)
    for img in soup.find_all('img'):
        src = img.get('src')
        if src != None:
            link = urljoin(baseUrl, src)
            ret.append(link)

    return ret

for response in getUrlGeneratorStyle("http://m.yle.fi", 2, http):
    for image in findImages(response.data, response.task.url):
        print image
