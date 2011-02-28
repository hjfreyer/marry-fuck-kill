#!/usr/bin/env python
#
# Copyright 2011 Hunter Freyer and Michael Kelly
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import logging
import random
import urllib2
from django.utils import simplejson

from google.appengine.api import users
from google.appengine.ext import db

import models

GOOGLE_API_KEY = 'ABQIAAAA4AIACTDq7g0UgEEe0e4XcBScM50iuTtmL4hn6SVBcuHEk5GnyBRYi46EgwfJeghlh-_jWgC9BbPapQ'
SEARCH_REFERER = 'http://marry-fuck-kill.appspot.com'


class EntityValidationError(Exception): pass


def GetStatsUrlsForTriple(triple, w=160, h=85):
  """Returns a list of stats URLs for the given triple.

  Args:
    triple: (Triple) triple to examine
    w: (int) Optional. Width of each chart image.
    h: (int) Optional. Height of each chart image.

  Returns:
    [str, str, str]: URLs for the Triple's three Entities.
  """
  counts = [GetEntityVoteCounts(triple.one),
            GetEntityVoteCounts(triple.two),
            GetEntityVoteCounts(triple.three)]

  urls = []
  overall_max = max([max(c) for c in counts])
  for count in counts:
    urls.append('http://chart.apis.google.com/chart'
                '?chxr=0,0,%(max)d'
                '&chxt=y'
                '&chbh=a'
                '&chs=%(w)dx%(h)d'
                '&cht=bvg'
                '&chco=9911BB,C76FDD,63067A'
                '&chds=0,%(max)d,0,%(max)d,0,%(max)d'
                '&chd=t:%(m)d|%(f)d|%(k)d'
                '&chdl=Marry|Fuck|Kill'
                '&chdlp=r' % (dict(m=count[0], f=count[1], k=count[2],
                              max=overall_max,
                              w=w,
                              h=h)))
  return urls


def GetEntityVoteCounts(entity):
  """Returns the marry, fuck, and kill vote counts for an entity.

  Args:
    entity: (Entity) the entity to examine
    w: (int) image width
    h: (int) image height

  Returns:
    [int, int, int]: marry, fuck, and kill vote counts
  """
  m = entity.assignment_reference_marry_set.count()
  f = entity.assignment_reference_fuck_set.count()
  k = entity.assignment_reference_kill_set.count()
  return [m, f, k]


def MakeEntity(name, query, user_ip, original_url):
  """Makes an Entity with the given attributes.

  Ensures that the URL actually comes from a google search for the
  query. If it doesn't, may throw an EntityValidationError.
  """
  # This must be synchronized with the number of pages the creation
  # interface shows. Raising it incurs a large performance penalty.
  check_pages = 2

  images = GetImagesForQuery(query, check_pages, user_ip)
  images_by_url = dict([(image['unescapedUrl'], image) for image in images])
  logging.info('validate_request: possible valid urls = %s',
               list(images_by_url))
  if original_url not in images_by_url:
    logging.error("URL '%s' is not in result set for query '%s'. "
                  "Result set over %d pages is: %s" % (
        orignal_url, query, check_pages, list(images_by_url)))
    raise EntityValidationError(
      "URL '%s' is not in result set for query '%s'." % (original_url,
                                                         query))

  tb_url = images_by_url[original_url]['tbUrl']

  # Get the thumbnail URL for the entity.  This could throw
  # URLError. We'll let it bubble up.
  fh = urllib2.urlopen(tb_url)
  logging.info('Downloading %s' % tb_url)
  data = fh.read()
  logging.info('Downloaded %s bytes' % len(data))

  entity = models.Entity(name=name,
                         data=data,
                         query=query,
                         original_url=original_url)
  entity.put()

  return entity


def MakeTriple(entities, creator, creator_ip):
  """Create the named triple.

  Args:
    entities: a data structure built in MakeSubmitHandler.
    creator: the user who created the Triple.
    creator_ip: IP address of the request to make this triple.

  The only non-obvious part of this is that we check that q[1-3] actually
  include u[1-3]. This is to prevent users from adding any URL they
  please.
  """
  for i in range(len(entities)):
    for k in ['n', 'u', 'q', 'ou']:
      if not entities[i][k]:
        raise ValueError("Entity %s missing attribute '%s'" % (i, k))

  # This may raise a URLError or EntityValidatationError.
  one = MakeEntity(name=entities[0]['n'],
                   query=entities[0]['q'],
                   original_url=entities[0]['ou'],
                   user_ip=creator_ip)
  two = MakeEntity(name=entities[1]['n'],
                   query=entities[1]['q'],
                   original_url=entities[1]['ou'],
                   user_ip=creator_ip)
  three = MakeEntity(name=entities[2]['n'],
                     query=entities[2]['q'],
                     original_url=entities[2]['ou'],
                     user_ip=creator_ip)

  triple = models.Triple(one=one, two=two, three=three,
                         creator=creator,
                         creatorip=creator_ip,
                         rand=random.random())
  triple.put()
  return triple


def MakeAssignment(triple_id, v1, v2, v3, user, user_ip):
  """Create a new assignment.

  Args:
    request: the POST request from the client
    user: the user who made the assignment request
  """
  values = [v1, v2, v3]
  if set(values) != set(['m', 'f', 'k']):
    return None

  try:
    triple_id = long(triple_id)
  except ValueError:
    logging.error("make_assignment: bad triple key '%s'", triple_id)

  triple = models.Triple.get_by_id(triple_id)
  logging.debug('triple = %s', triple)
  if triple is None:
    logging.error('make_assignment: No triple with key %s', triple_id)
    return None

  # We get an entity->action map from the client, but we need to reverse
  # it to action->entity to update the DB.
  triple_entities = [triple.one,
                     triple.two,
                     triple.three]
  entities = {}
  for i in range(len(values)):
    # Items in values are guaranteed to be 'm', 'f', 'k' (check above)
    entities[values[i]] = triple_entities[i]

  if (entities['m'] is None or entities['f'] is None
      or entities['k'] is None):
    logging.error('Not all non-None: marry = %s, fuck = %s, kill = %s',
            entities['m'], entities['f'], entities['k'])
    return None

  assign = models.Assignment(triple=triple,
                             marry=entities['m'],
                             fuck=entities['f'],
                             kill=entities['k'],
                             user=user,
                             userip=str(user_ip))
  assign.put()
  logging.info("Assigned m=%s, f=%s, k=%s to %s", entities['m'],
      entities['f'], entities['k'], triple)
  return assign


def GetImagesForQuery(query, check_pages, userip):
  start = 0
  images = []
  query = query.encode('utf-8')
  for page in range(check_pages):
    url = ('https://ajax.googleapis.com/ajax/services/search/images'
           '?v=1.0'
           '&safe=moderate'
           '&rsz=8'
           '&userip=%(ip)s'
           '&q=%(q)s'
           '&start=%(start)s'
           '&key=%(key)s' % {'q': urllib2.quote(query),
                             'start': start,
                             'ip': userip,
                             'key': GOOGLE_API_KEY})
    logging.info('GetImagesForQuery: query url=%s', url)
    req = urllib2.Request(url, None, {'Referer': SEARCH_REFERER})
    # This may raise a DownloadError
    results = simplejson.load(urllib2.urlopen(req))
    images += results['responseData']['results']

    # Make sure there are more pages before advancing 'start'
    pages = results['responseData']['cursor']['pages']
    logging.info('GetImagesForQuery: %d results so far,'
                 ' on page %s of %s (will try up to %s)',
                 len(images), page+1, len(pages), check_pages)
    if len(pages) > page + 1:
      start = pages[page+1]['start']
    else:
      break

  return images

